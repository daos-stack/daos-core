//
// (C) Copyright 2019-2020 Intel Corporation.
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
package storage

import (
	"fmt"
	"strconv"
	"strings"
)

func concat(base string, idx int32) string {
	return fmt.Sprintf("%s-%d", base, idx)
}

func getIndex(varIdx ...int32) int32 {
	if len(varIdx) == 0 {
		varIdx = append(varIdx, 1)
	}

	return varIdx[0]
}

func MockNvmeDeviceHealth(varIdx ...int32) *NvmeDeviceHealth {
	idx := getIndex(varIdx...)
	tWarn := false
	if idx > 0 {
		tWarn = true
	}
	return &NvmeDeviceHealth{
		Temp:     uint32(idx),
		TempWarn: tWarn,
	}
}

func MockNvmeNamespace(varIdx ...int32) *NvmeNamespace {
	idx := getIndex(varIdx...)
	return &NvmeNamespace{
		ID:   idx,
		Size: uint64(idx),
	}
}

func MockNvmeController(varIdx ...int32) *NvmeController {
	idx := getIndex(varIdx...)

	return &NvmeController{
		Model:       concat("model", idx),
		Serial:      concat("serial", idx),
		PciAddr:     concat("pciAddr", idx),
		FwRev:       concat("fwRev", idx),
		SocketID:    idx,
		HealthStats: MockNvmeDeviceHealth(idx),
		Namespaces:  []*NvmeNamespace{MockNvmeNamespace(idx)},
	}
}

func MockScmModule(varIdx ...int32) *ScmModule {
	idx := uint32(getIndex(varIdx...))

	return &ScmModule{
		ChannelID:       idx,
		ChannelPosition: idx,
		ControllerID:    idx,
		SocketID:        idx,
		PhysicalID:      idx,
		Capacity:        uint64(idx),
	}
}

func MockScmNamespace(varIdx ...int32) *ScmNamespace {
	idx := getIndex(varIdx...)

	return &ScmNamespace{
		UUID:        MockUUID(varIdx...),
		BlockDevice: fmt.Sprintf("/dev/pmem%d", idx),
		Name:        fmt.Sprintf("pmem%d", idx),
		NumaNode:    uint32(idx),
		Size:        uint64(idx),
	}
}

func MockUUID(varIdx ...int32) string {
	idx := getIndex(varIdx...)
	idxStr := strconv.Itoa(int(idx))

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		strings.Repeat(idxStr, 8),
		strings.Repeat(idxStr, 4),
		strings.Repeat(idxStr, 4),
		strings.Repeat(idxStr, 4),
		strings.Repeat(idxStr, 12),
	)
}
