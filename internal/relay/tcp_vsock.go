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
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/mdlayher/vsock"
	"github.com/palantir/stacktrace"
)

type TCPtoVsock struct {
	AbstractDuplexRelay
}

func NewTCPtoVsock(
	healthCheckInterval time.Duration,
	tcpAddress,
	vsockPort string,
	bufferSize int,
) (*TCPtoVsock, error) {
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

	return &TCPtoVsock{
		AbstractDuplexRelay{
			healthCheckInterval: healthCheckInterval,
			sourceName:          "vsock",
			destinationName:     "TCP connection",
			destinationAddr:     vsockPort,
			bufferSize:          bufferSize,
			dialSourceConn: func(ctx context.Context) (net.Conn, error) {
				port, err := strconv.ParseInt(vsockPort, 10, 32)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"could not parse specified vsock port number %s",
						vsockPort,
					)
				}
				conn, err := vsock.Dial(vsock.Host, uint32(5001), nil)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to dial vsock port: %v",
						port,
					)
				}
				return conn, nil
			},
			listenTargetConn: func(ctx context.Context) (net.Listener, error) {
				var lc net.ListenConfig
				listener, err := lc.Listen(ctx, "tcp", tcpAddress)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to listen at tcp address: %s",
						tcpAddress,
					)
				}
				log.Println("LISTEN SUCCESS")

				return listener, nil
			},
		},
	}, nil
}
