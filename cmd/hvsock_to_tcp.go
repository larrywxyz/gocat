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

func NewHvsockToTCPCmd() *cobra.Command {
	var hvsockToTcpPath string
	var hvsockToTcpAddress string

	cmdInstance := &cobra.Command{
		Use:   "hvsock-to-tcp",
		Short: "relay from a hvsock source to tcp clients",
		Long:  `relay from a hvsock source to tcp clients`,
		RunE: func(command *cobra.Command, args []string) error {
			// nolint: gocritic
			if len(hvsockToTcpPath) < 0 {
				return stacktrace.NewError("blank/empty `src` specified")
			}

			// nolint: gocritic
			if len(hvsockToTcpAddress) < 0 {
				return stacktrace.NewError("blank/empty `dst` specified")
			}

			err := relay.NewHvsockToTCPRelay(
				hvsockToTcpPath,
				hvsockToTcpAddress,
			)
			if err != nil {
				return stacktrace.Propagate(err, "couldn't create relay from hvsock to TCP")
			}
			return nil
		},
	}

	cmdInstance.Flags().StringVar(
		&hvsockToTcpPath,
		"src",
		"",
		"source of hvsock",
	)
	_ = cmdInstance.MarkFlagRequired("src")
	cmdInstance.Flags().StringVar(
		&hvsockToTcpAddress,
		"dst",
		"",
		"destination to TCP listen",
	)
	_ = cmdInstance.MarkFlagRequired("dst")

	return cmdInstance
}
