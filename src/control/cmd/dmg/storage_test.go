//
// (C) Copyright 2019 Intel Corporation.
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

package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/daos-stack/daos/src/control/client"
)

func TestStorageCommands(t *testing.T) {
	runCmdTests(t, []cmdTest{
		{
			"Format without reformat",
			"storage format",
			"ConnectClients StorageFormat-false",
			nil,
		},
		{
			"Format with reformat",
			"storage format --reformat",
			"ConnectClients StorageFormat-true",
			nil,
		},
		{
			"Scan",
			"storage scan",
			strings.Join([]string{
				"ConnectClients",
				fmt.Sprintf("StorageScan-%+v", &client.StorageScanReq{}),
			}, " "),
			nil,
		},
		{
			"Prepare without force",
			"storage prepare",
			"ConnectClients",
			fmt.Errorf("consent not given"),
		},
		{
			"Prepare with nvme-only and scm-only",
			"storage prepare --force --nvme-only --scm-only",
			"ConnectClients",
			fmt.Errorf("nvme-only and scm-only options should not be set together"),
		},
		{
			"Prepare with scm-only",
			"storage prepare --force --scm-only",
			"ConnectClients StoragePrepare",
			nil,
		},
		{
			"Prepare with nvme-only",
			"storage prepare --force --nvme-only",
			"ConnectClients StoragePrepare",
			nil,
		},
		{
			"Prepare with non-existent option",
			"storage prepare --force --nvme",
			"",
			fmt.Errorf("unknown flag `nvme'"),
		},
		{
			"Prepare with force and reset",
			"storage prepare --force --reset",
			"ConnectClients StoragePrepare",
			nil,
		},
		{
			"Prepare with force",
			"storage prepare --force",
			"ConnectClients StoragePrepare",
			nil,
		},
		{
			"Set FAULTY device status",
			"storage setfaulty --devuuid abcd",
			"ConnectClients StorageSetFaulty-dev_uuid:\"abcd\" ",
			nil,
		},
		{
			"Set FAULTY device status without device specified",
			"storage setfaulty",
			"ConnectClients StorageSetFaulty",
			fmt.Errorf("the required flag `-u, --devuuid' was not specified"),
		},
		{
			"Nonexistent subcommand",
			"storage quack",
			"",
			fmt.Errorf("Unknown command"),
		},
	})
}
