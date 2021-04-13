//
// (C) Copyright 2020-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package events

import (
	"github.com/daos-stack/daos/src/control/common"
	sharedpb "github.com/daos-stack/daos/src/control/common/proto/shared"
)

// EngineStateInfo describes details of a engine's state.
type EngineStateInfo struct {
	InstanceIdx uint32 `json:"instance_idx"`
	ExitErr     error  `json:"-"`
}

func (rsi *EngineStateInfo) isExtendedInfo() {}

// GetEngineStateInfo returns extended info if of type EngineStateInfo.
func (evt *RASEvent) GetEngineStateInfo() *EngineStateInfo {
	if ei, ok := evt.ExtendedInfo.(*EngineStateInfo); ok {
		return ei
	}

	return nil
}

// EngineStateInfoFromProto converts event info from proto to native format.
func EngineStateInfoFromProto(pbInfo *sharedpb.RASEvent_EngineStateInfo) (*EngineStateInfo, error) {
	rsi := &EngineStateInfo{
		InstanceIdx: pbInfo.EngineStateInfo.GetInstance(),
	}
	if pbInfo.EngineStateInfo.GetErrored() {
		rsi.ExitErr = common.ExitStatus(pbInfo.EngineStateInfo.GetError())
	}

	return rsi, nil
}

// EngineStateInfoToProto converts event info from native to proto format.
func EngineStateInfoToProto(rsi *EngineStateInfo) (*sharedpb.RASEvent_EngineStateInfo, error) {
	pbInfo := &sharedpb.RASEvent_EngineStateInfo{
		EngineStateInfo: &sharedpb.RASEvent_EngineStateEventInfo{
			Instance: rsi.InstanceIdx,
		},
	}
	if rsi.ExitErr != nil {
		pbInfo.EngineStateInfo.Errored = true
		pbInfo.EngineStateInfo.Error = rsi.ExitErr.Error()
	}

	return pbInfo, nil
}

// NewEngineFailedEvent creates a specific EngineFailed event from given inputs.
func NewEngineFailedEvent(hostname string, instanceIdx uint32, rank uint32, exitErr common.ExitStatus) *RASEvent {
	return New(&RASEvent{
		Msg:      "DAOS engine exited unexpectedly",
		ID:       RASEngineFailed,
		Hostname: hostname,
		Rank:     rank,
		Type:     RASTypeStateChange,
		Severity: RASSeverityError,
		ExtendedInfo: &EngineStateInfo{
			InstanceIdx: instanceIdx,
			ExitErr:     exitErr,
		},
	})
}
