//
// (C) Copyright 2018-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

// Package spdk provides Go bindings for SPDK
package spdk

// CGO_CFLAGS & CGO_LDFLAGS env vars can be used
// to specify additional dirs.

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lnvme_control -lspdk

#include "stdlib.h"
#include "daos_srv/control.h"
#include "spdk/stdinc.h"
#include "spdk/string.h"
#include "spdk/env.h"
#include "spdk/nvme.h"
#include "spdk/vmd.h"
#include "include/nvme_control.h"
#include "include/nvme_control_common.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/logging"
)

// Env is the interface that provides SPDK environment management.
type Env interface {
	InitSPDKEnv(logging.Logger, *EnvOptions) error
	FiniSPDKEnv(logging.Logger, *EnvOptions)
}

// EnvImpl is a an implementation of the Env interface.
type EnvImpl struct{}

// Rc2err returns error from label and rc.
func Rc2err(label string, rc C.int) error {
	return fmt.Errorf("%s: %d", label, rc)
}

// EnvOptions describe parameters to be used when initializing a processes
// SPDK environment.
type EnvOptions struct {
	MemSize        int      // size in MiB to be allocated to SPDK proc
	PciIncludeList []string // restrict SPDK device access
	DisableVMD     bool     // flag if VMD devices should not be included
}

func (o *EnvOptions) sanitizeIncludeList(log logging.Logger) error {
	if !o.DisableVMD {
		// DPDK will not accept VMD backing device addresses
		// so convert to VMD address
		newIncludeList, err := revertBackingToVmd(log, o.PciIncludeList)
		if err != nil {
			return err
		}
		o.PciIncludeList = newIncludeList
	}

	return nil
}

// pciListToC allocates memory and populate array of C SPDK PCI addresses.
func pciListToC(log logging.Logger, inAddrs []string) (*unsafe.Pointer, error) {
	var tmpAddr *C.struct_spdk_pci_addr

	structSize := unsafe.Sizeof(*tmpAddr)
	outAddrs := C.calloc(C.ulong(len(inAddrs)), C.ulong(structSize))

	for i, inAddr := range inAddrs {
		log.Debugf("adding %s to spdk_env_opts include list", inAddr)

		offset := uintptr(i) * structSize
		tmpAddr = (*C.struct_spdk_pci_addr)(unsafe.Pointer(uintptr(outAddrs) + offset))
		cs := C.CString(inAddr)

		if rc := C.spdk_pci_addr_parse(tmpAddr, C.CString(inAddr)); rc != 0 {
			C.free(unsafe.Pointer(outAddrs))

			return nil, Rc2err("spdk_pci_addr_parse()", rc)
		}

		C.free(unsafe.Pointer(cs))
	}

	return &outAddrs, nil
}

// revertBackingToVmd converts VMD backing device PCI addresses (with the VMD
// address encoded in the domain component of the PCI address) back to the PCI
// address of the VMD e.g. [5d0505:01:00.0, 5d0505:03:00.0] -> [0000:5d:05.5].
//
// Many assumptions are made as to the input and output PCI address structure in
// the conversion.
func revertBackingToVmd(log logging.Logger, pciAddrs []string) ([]string, error) {
	var outAddrs []string

	for _, inAddr := range pciAddrs {
		domain, _, _, _, err := common.ParsePCIAddress(inAddr)
		if err != nil {
			return nil, err
		}
		if domain == 0 {
			outAddrs = append(outAddrs, inAddr)
			continue
		}

		domainStr := fmt.Sprintf("%x", domain)
		if len(domainStr) != 6 {
			return nil, errors.New("unexpected length of domain")
		}

		outAddr := fmt.Sprintf("0000:%s:%s.%s",
			domainStr[:2], domainStr[2:4], domainStr[5:])
		if !common.Includes(outAddrs, outAddr) {
			log.Debugf("replacing backing device %s with vmd %s", inAddr, outAddr)
			outAddrs = append(outAddrs, outAddr)
		}
	}

	return outAddrs, nil
}

// InitSPDKEnv initializes the SPDK environment.
//
// SPDK relies on an abstraction around the local environment
// named env that handles memory allocation and PCI device operations.
// The library must be initialized first.
func (e *EnvImpl) InitSPDKEnv(log logging.Logger, opts *EnvOptions) error {
	log.Debugf("spdk init go opts: %+v", opts)

	if err := opts.sanitizeIncludeList(log); err != nil {
		return errors.Wrap(err, "sanitizing PCI include list")
	}

	listPtr, err := pciListToC(log, opts.PciIncludeList)
	if err != nil {
		return err
	}
	defer func() {
		if listPtr != nil {
			C.free(*listPtr)
		}
	}()

	envCtx := C.CString("--log-level=lib.eal:4")
	defer C.free(unsafe.Pointer(envCtx))

	retPtr := C.daos_spdk_init(C.int(opts.MemSize), envCtx,
		C.ulong(len(opts.PciIncludeList)),
		(*C.struct_spdk_pci_addr)(*listPtr))
	if err := checkRet(retPtr, "daos_spdk_init()"); err != nil {
		return err
	}
	clean(retPtr)

	if opts.DisableVMD {
		return nil
	}

	if rc := C.spdk_vmd_init(); rc != 0 {
		return Rc2err("spdk_vmd_init()", rc)
	}

	return nil
}

// FiniSPDKEnv initializes the SPDK environment.
func (e *EnvImpl) FiniSPDKEnv(log logging.Logger, opts *EnvOptions) {
	log.Debugf("spdk fini go opts: %+v", opts)

	C.spdk_env_fini()

	// TODO: enable when vmd_fini supported in daos spdk version
	//	if opts.DisableVMD {
	//		return nil
	//	}
	//
	//	if rc := C.spdk_vmd_fini(); rc != 0 {
	//		return Rc2err("spdk_vmd_fini()", rc)
	//	}
}
