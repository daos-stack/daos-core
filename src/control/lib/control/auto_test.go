//
// (C) Copyright 2020 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

package control

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/logging"
	"github.com/daos-stack/daos/src/control/security"
	"github.com/daos-stack/daos/src/control/server/config"
	"github.com/daos-stack/daos/src/control/server/ioserver"
	"github.com/daos-stack/daos/src/control/server/storage"
)

func TestControl_AutoConfig_checkStorage(t *testing.T) {
	hostRespOneScanFail := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "standard"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "bothFailed"),
		},
	}
	hostRespScanFail := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "bothFailed"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "bothFailed"),
		},
	}
	hostRespNoScmNs := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "standard"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "standard"),
		},
	}
	hostRespOneWithScmNs := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "withNamespace"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "standard"),
		},
	}
	hostRespWithScmNs := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "withNamespace"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "withNamespace"),
		},
	}
	hostRespWithScmNss := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "withNamespaces"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "withNamespaces"),
		},
	}
	hostRespWithSingleSSD := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "withSingleSSD"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "withSingleSSD"),
		},
	}
	hostRespWithSSDs := []*HostResponse{
		{
			Addr:    "host1",
			Message: MockServerScanResp(t, "withSpaceUsage"),
		},
		{
			Addr:    "host2",
			Message: MockServerScanResp(t, "withSpaceUsage"),
		},
	}
	ioCfgWithSSDs := func(numa int) *ioserver.Config {
		var pciAddrs []string
		for _, c := range MockServerScanResp(t, "withSpaceUsage").Nvme.Ctrlrs {
			if int(c.Socketid) == numa {
				pciAddrs = append(pciAddrs, c.Pciaddr)
			}
		}
		return ioserver.NewConfig().
			WithScmClass(storage.ScmClassDCPM.String()).
			WithScmMountPoint(fmt.Sprintf("/mnt/daos%d", numa)).
			WithScmDeviceList(fmt.Sprintf("/dev/pmem%d", numa)).
			WithBdevClass(storage.BdevClassNvme.String()).
			WithBdevDeviceList(pciAddrs...)
	}

	for name, tc := range map[string]struct {
		minPmem       int
		minNvme       int
		numNumaRet    int
		numNumaErr    error
		uErr          error
		hostResponses []*HostResponse
		expConfigOut  *config.Server
		expCheckErr   error
	}{
		"invoker error": {
			uErr:          errors.New("unary error"),
			hostResponses: hostRespOneWithScmNs,
			expCheckErr:   errors.New("unary error"),
		},
		"host storage scan failed": {
			hostResponses: hostRespOneScanFail,
			expCheckErr:   errors.New("1 host had errors"),
		},
		"host storage scan failed on multiple hosts": {
			hostResponses: hostRespScanFail,
			expCheckErr:   errors.New("2 hosts had errors"),
		},
		"host storage scan no hosts": {
			hostResponses: []*HostResponse{},
			expCheckErr:   errors.New("no host responses"),
		},
		"host storage mismatch": {
			hostResponses: hostRespOneWithScmNs,
			expCheckErr:   errors.New("storage hardware not consistent across hosts"),
		},
		"no min pmem and no numa": {
			hostResponses: hostRespNoScmNs,
			expCheckErr:   errors.New("no NUMA nodes reported on host"),
		},
		"no min pmem and numa error": {
			numNumaErr:    errors.New("get num numa nodes failed"),
			hostResponses: hostRespNoScmNs,
			expCheckErr:   errors.New("get num numa nodes failed"),
		},
		"no min pmem and 2 numa and 0 pmems present": {
			numNumaRet:    2,
			hostResponses: hostRespNoScmNs,
			expCheckErr:   errors.New("insufficient number of pmem devices, want 2 got 0"),
		},
		"2 min pmem and 1 pmems present": {
			minPmem:       2,
			hostResponses: hostRespWithScmNs,
			expCheckErr:   errors.New("insufficient number of pmem devices, want 2 got 1"),
		},
		"2 min pmem and 2 pmems present": {
			minPmem:       2,
			hostResponses: hostRespWithScmNss,
			expCheckErr:   errors.New("insufficient number of nvme devices for numa node 0, want 1 got 0"),
		},
		"no min nvme and 1 ctrlr present on 1 numa node": {
			minPmem:       2,
			hostResponses: hostRespWithSingleSSD,
			expCheckErr:   errors.New("insufficient number of nvme devices for numa node 1, want 1 got 0"),
		},
		"no min nvme and multiple ctrlrs present on 2 numa nodes": {
			minPmem:       2,
			hostResponses: hostRespWithSSDs,
			expConfigOut: config.DefaultServer().
				WithServers(ioCfgWithSSDs(0), ioCfgWithSSDs(1)),
		},
	} {
		t.Run(name, func(t *testing.T) {
			log, buf := logging.NewTestLogger(t.Name())
			defer common.ShowBufferOnFailure(t, buf)

			mi := NewMockInvoker(log, &MockInvokerConfig{
				UnaryError: tc.uErr,
				UnaryResponse: &UnaryResponse{
					Responses: tc.hostResponses,
				},
			})

			mockGetNumaCount := func() (int, error) {
				return tc.numNumaRet, tc.numNumaErr
			}

			req := &ConfigGenerateReq{
				NumPmem: tc.minPmem,
				NumNvme: tc.minNvme,
				Client:  mi,
				Log:     log,
			}

			gotCfg, gotCheckErr := req.checkStorage(context.Background(), mockGetNumaCount)
			common.CmpErr(t, tc.expCheckErr, gotCheckErr)
			if tc.expCheckErr != nil {
				return
			}

			cmpOpts := []cmp.Option{
				cmpopts.IgnoreUnexported(security.CertificateConfig{}, config.Server{}),
				cmpopts.IgnoreFields(config.Server{}, "GetDeviceClassFn"),
			}
			if diff := cmp.Diff(tc.expConfigOut, gotCfg, cmpOpts...); diff != "" {
				t.Fatalf("output cfg doesn't match (-want, +got):\n%s\n", diff)
			}
		})
	}
}
