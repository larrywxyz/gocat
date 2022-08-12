// Copyright 2018 SumUp Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package relay

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/linuxkit/virtsock/pkg/hvsock"
	"github.com/palantir/stacktrace"
)

type HvsockTCP struct {
	AbstractDuplexRelay
}

func NewHvsockTcp(
	healthCheckInterval time.Duration,
	hvsockPath,
	tcpAddress string,
	bufferSize int,
) (*HvsockTCP, error) {
	tcpAddressParts := strings.Split(tcpAddress, ":")
	if len(tcpAddressParts) != 2 {
		return nil, stacktrace.NewError(
			"wrong format for tcp address %s. Expected <addr>:<port>",
			tcpAddress,
		)
	}

	_, err := strconv.ParseInt(tcpAddressParts[1], 10, 32)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"could not parse specified port number %s",
			tcpAddressParts[1],
		)
	}

	return &HvsockTCP{
		AbstractDuplexRelay{
			healthCheckInterval: healthCheckInterval,
			bufferSize:          bufferSize,
			sourceName:          "TCP connection",
			destinationName:     "hvsock",
			destinationAddr:     tcpAddress,
			dialSourceConn: func(ctx context.Context) (net.Conn, error) {
				dialer := &net.Dialer{}
				conn, err := dialer.DialContext(
					ctx,
					"tcp",
					tcpAddress,
				)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to dial TCP address: %s",
						tcpAddress,
					)
				}

				return conn, nil
			},
			listenTargetConn: func(ctx context.Context) (net.Listener, error) {
				split := strings.Split(hvsockPath, ":")

				VMID, err := hvsock.GUIDFromString(split[0])
				if err != nil {
					return nil, err
				}
				ServiceID, err := hvsock.GUIDFromString(split[1])
				if err != nil {
					return nil, err
				}
				addrz := hvsock.Addr{
					VMID:      VMID,
					ServiceID: ServiceID,
				}
				lis, err := hvsock.Listen(addrz)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to dial unix address: %s",
						hvsockPath,
					)
				}

				return lis, nil
			},
		},
	}, nil
}
