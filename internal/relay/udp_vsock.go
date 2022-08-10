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

	"github.com/mdlayher/vsock"
	"github.com/palantir/stacktrace"
)

type UDPtoVsock struct {
	AbstractDuplexRelay
}

func NewUDPtoVsock(
	healthCheckInterval time.Duration,
	udpAddress,
	vsockPort string,
	bufferSize int,
) (*UDPtoVsock, error) {
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

	return &UDPtoVsock{
		AbstractDuplexRelay{
			healthCheckInterval: healthCheckInterval,
			sourceName:          "UDP connection",
			destinationName:     "vsock",
			destinationAddr:     vsockPort,
			bufferSize:          bufferSize,
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
				port, err := strconv.ParseInt(vsockPort, 10, 32)
				lis, err := vsock.Listen(uint32(port), nil)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to listen at vsock port: %s",
						vsockPort,
					)
				}
				return lis, nil
			},
		},
	}, nil
}
