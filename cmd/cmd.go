/*
 * Copyright 2017 The GoStor Authors All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package cmd parses the parameters and runs GoFS
package cmd

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gofs",
		Short: "Gofs is a Fuse supported metadata server for object cloud storage server.",
		Long:  `GoFS is a fuse supported metadata server for object storage server. Use it to store photos, videos, VMs, containers, log files, or any blob of data as objects on your object storage server.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	cmd.AddCommand(
		newServerCommand(),
	)
	return cmd
}
