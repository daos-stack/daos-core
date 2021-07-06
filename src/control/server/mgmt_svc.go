//
// (C) Copyright 2018-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package server

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/daos-stack/daos/src/control/common"
	mgmtpb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	"github.com/daos-stack/daos/src/control/events"
	"github.com/daos-stack/daos/src/control/lib/control"
	"github.com/daos-stack/daos/src/control/logging"
	"github.com/daos-stack/daos/src/control/system"
)

// mgmtSvc implements (the Go portion of) Management Service, satisfying
// mgmtpb.MgmtSvcServer.
type mgmtSvc struct {
	mgmtpb.UnimplementedMgmtSvcServer
	log               logging.Logger
	harness           *EngineHarness
	membership        *system.Membership // if MS leader, system membership list
	sysdb             *system.Database
	rpcClient         control.UnaryInvoker
	events            *events.PubSub
	clientNetworkHint *mgmtpb.ClientNetHint
	joinReqs          joinReqChan
	groupUpdateReqs   chan struct{}
}

func newMgmtSvc(h *EngineHarness, m *system.Membership, s *system.Database, c control.UnaryInvoker, p *events.PubSub) *mgmtSvc {
	return &mgmtSvc{
		log:               h.log,
		harness:           h,
		membership:        m,
		sysdb:             s,
		rpcClient:         c,
		events:            p,
		clientNetworkHint: new(mgmtpb.ClientNetHint),
		joinReqs:          make(joinReqChan),
		groupUpdateReqs:   make(chan struct{}),
	}
}

// checkSystemRequest sanity checks that a request is not nil and
// has been sent to the correct system.
func (svc *mgmtSvc) checkSystemRequest(req proto.Message) error {
	if common.InterfaceIsNil(req) {
		return errors.New("nil request")
	}
	if sReq, ok := req.(interface{ GetSys() string }); ok {
		if sReq.GetSys() != svc.sysdb.SystemName() {
			return FaultWrongSystem(sReq.GetSys(), svc.sysdb.SystemName())
		}
	}
	return nil
}

// checkLeaderRequest performs sanity-checking on a request that must
// be run on the current MS leader.
func (svc *mgmtSvc) checkLeaderRequest(req proto.Message) error {
	if err := svc.sysdb.CheckLeader(); err != nil {
		return err
	}
	return svc.checkSystemRequest(req)
}

// checkReplicaRequest performs sanity-checking on a request that must
// be run on a MS replica.
func (svc *mgmtSvc) checkReplicaRequest(req proto.Message) error {
	if err := svc.sysdb.CheckReplica(); err != nil {
		return err
	}
	return svc.checkSystemRequest(req)
}
