//
// (C) Copyright 2018 Intel Corporation.
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

package client

import (
	pb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

// Control interface provides connection handling capabilities.
type Control interface {
	connect(string, credentials.TransportCredentials) error
	disconnect() error
	connected() (connectivity.State, bool)
	getAddress() string
	getCtlClient() pb.MgmtCtlClient
	getSvcClient() pb.MgmtSvcClient
}

// control is an abstraction around the Mgmt{Control,Svc}Clients
// generated by gRPC. It provides a simplified mechanism so users can
// minimize their use of protobuf datatypes.
type control struct {
	ctlClient pb.MgmtCtlClient
	svcClient pb.MgmtSvcClient
	gconn     *grpc.ClientConn
}

// connect provides an easy interface to connect to Mgmt DAOS server.
//
// It takes address and port in a string.
//	addr: address and port number separated by a ":"
func (c *control) connect(addr string, creds credentials.TransportCredentials) (err error) {
	var opts []grpc.DialOption
	if creds != nil {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return
	}
	c.ctlClient = pb.NewMgmtCtlClient(conn)
	c.svcClient = pb.NewMgmtSvcClient(conn)
	c.gconn = conn

	return
}

// disconnect terminates the underlying channel used by the grpc
// client service.
func (c *control) disconnect() error { return c.gconn.Close() }

func checkState(state connectivity.State) bool {
	return (state == connectivity.Idle || state == connectivity.Ready)
}

// connected determines if the underlying socket connection is alive and well.
func (c *control) connected() (state connectivity.State, ok bool) {
	if c.gconn == nil {
		return
	}

	state = c.gconn.GetState()
	return state, checkState(state)
}

// getAddress returns the target address of the connection.
func (c *control) getAddress() string { return c.gconn.Target() }

func (c *control) getCtlClient() pb.MgmtCtlClient { return c.ctlClient }
func (c *control) getSvcClient() pb.MgmtSvcClient { return c.svcClient }
