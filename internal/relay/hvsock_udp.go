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

type HvsockUDP struct {
	AbstractDuplexRelay
}

func NewHvsockUdp(
	healthCheckInterval time.Duration,
	hvsockPath,
	udpAddress string,
	bufferSize int,
) (*HvsockUDP, error) {
	udpAddressParts := strings.Split(udpAddress, ":")
	if len(udpAddressParts) != 2 {
		return nil, stacktrace.NewError(
			"wrong format for udp address %s. Expected <addr>:<port>",
			udpAddress,
		)
	}

	_, err := strconv.ParseInt(udpAddressParts[1], 10, 32)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"could not parse specified port number %s",
			udpAddressParts[1],
		)
	}

	return &HvsockUDP{
		AbstractDuplexRelay{
			healthCheckInterval: healthCheckInterval,
			bufferSize:          bufferSize,
			sourceName:          "hvsock",
			destinationName:     "UDP connection",
			destinationAddr:     udpAddress,
			dialSourceConn: func(ctx context.Context) (net.Conn, error) {
				dialer := &net.Dialer{}
				conn, err := dialer.DialContext(
					ctx,
					"tcp",
					udpAddress,
				)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to dial TCP address: %s",
						udpAddress,
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
				// var lc net.ListenConfig
				// listener, err := lc.Listen(ctx, "tcp", udpAddress)
				// if err != nil {
				// 	return nil, stacktrace.Propagate(
				// 		err,
				// 		"failed to listen at udp address: %s",
				// 		udpAddress,
				// 	)
				// }
				// log.Println("LISTEN SUCCESS")

				// return listener, nil
			},
		},
	}, nil
}
