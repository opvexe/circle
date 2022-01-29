/*
Copyright 2021 The SHUMIN Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"context"
	"github.com/spf13/cobra"
	"x6t.io/circle"
	"x6t.io/circle/services"
)

func NewCircleCommand() *cobra.Command {
	s := NewRunServerConfig()

	cmd := &cobra.Command{
		Use:  "circle",
		Long: "",
		RunE: func(cmd *cobra.Command, args []string) error {

			return Run(s, circle.SetupSignalHandler())
		},
		SilenceUsage: false,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of circle",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(circle.Get())
		},
	}

	cmd.AddCommand(versionCmd)
	return cmd
}

func Run(opt *RunServerConfig, ctx context.Context) error {
	s := services.NewAssignment(opt.PreRun(), ctx)

	return nil
}
