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
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log"
	"x6t.io/circle"
	"x6t.io/circle/services"
)

func NewCircleCommand() *cobra.Command {
	s := NewRunServerConfig()

	cmd := &cobra.Command{
		Use: "circle",
		Long: `The Micro Circle server Service verification and configuration service objects.
provide a list of timed collection tasks, handle WeChat sharing and Moments sharing tasks.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.Validate(); len(err) != 0 {
				return circle.NewAggregate(err)
			}
			return Run(s, circle.SetupSignalHandler())
		},
		SilenceUsage: true,
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
	log.Printf("Circle server turning on")
	c := cron.New()
	defer c.Stop()

	_, err := c.AddFunc(opt.Express, func() {
		if err := services.NewAssignment().Pub(opt.PreRun()); err != nil {
			log.Printf("corn task assignment tasks")
		}
	})

	if err != nil {
		return err
	}
	c.Start()

	<-ctx.Done()
	log.Printf("Circle server exiting")
	return nil
}
