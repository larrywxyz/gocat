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

package cmd

import (
	"github.com/palantir/stacktrace"
	"github.com/spf13/cobra"

	"gocat/relay"
)

func NewTCPToVsockCmd() *cobra.Command {
	var tcpToVsockPort string
	var tcpToVsockAddress string

	cmdInstance := &cobra.Command{
		Use:   "tcp-to-vsock",
		Short: "relay from a TCP source to vsock",
		Long:  `relay from a TCP source to vsock`,
		RunE: func(command *cobra.Command, args []string) error {
			// nolint: gocritic
			if len(tcpToVsockPort) < 0 {
				return stacktrace.NewError("blank/empty `src` specified")
			}

			// nolint: gocritic
			if len(tcpToVsockAddress) < 0 {
				return stacktrace.NewError("blank/empty `dst` specified")
			}

			err := relay.NewTCPToVsockRelay(
				tcpToVsockPort,
				tcpToVsockAddress,
			)
			if err != nil {
				return stacktrace.Propagate(err, "couldn't create relay from TCP to unix socket")
			}
			return nil
		},
	}

	cmdInstance.Flags().StringVar(&tcpToVsockPort, "src", "", "source of TCP address")
	_ = cmdInstance.MarkFlagRequired("src")
	cmdInstance.Flags().StringVar(
		&tcpToVsockAddress,
		"dst",
		"",
		"destination of vsock",
	)
	_ = cmdInstance.MarkFlagRequired("dst")
	return cmdInstance
}
