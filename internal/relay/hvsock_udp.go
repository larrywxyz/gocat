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

	"github.com/palantir/stacktrace"
	"github.com/sumup-oss/go-pkgs/logger"
	"github.com/linuxkit/virtsock/pkg/hvsock"
)

type HvsockUDP struct {
	AbstractDuplexRelay
}

func NewHvsockUdp(
	logger logger.Logger,
	healthCheckInterval time.Duration,
	hvsockPath,
	udpAddress string,
	bufferSize int,
) (*UnixSocketTCP, error) {
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

	return &UnixSocketTCP{
		AbstractDuplexRelay{
			healthCheckInterval: healthCheckInterval,
			logger:              logger,
			bufferSize:          bufferSize,
			sourceName:          "hvsock",
			destinationName:     "UDP connection",
			destinationAddr:     udpAddress,
			dialSourceConn: func(ctx context.Context) (net.Conn, error) {
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

				conn, err := hvsock.Dial(addrz)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to dial unix address: %s",
						hvsockPath,
					)
				}

				return conn, nil
			},
			listenTargetConn: func(ctx context.Context) (net.Listener, error) {
				var lc net.ListenConfig
				listener, err := lc.Listen(ctx, "udp", udpAddress)
				if err != nil {
					return nil, stacktrace.Propagate(
						err,
						"failed to listen at udp address: %s",
						udpAddress,
					)
				}
				return listener, nil
			},
		},
	}, nil
}
