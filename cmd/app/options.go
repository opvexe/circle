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
	"x6t.io/circle"
	"x6t.io/circle/services"
)

type RunServerConfig struct {
	*services.Options

	DebugMode bool
}

func NewRunServerConfig() *RunServerConfig {
	return &RunServerConfig{
		Options: services.NewOptions(),
	}
}

func (s *RunServerConfig) Flags() (fss services.NamedFlagSets) {
	fs := fss.FlagSet("circle")
	fs.BoolVar(&s.DebugMode, "debug", false, "Don't enable this if you don't know what it means.")
	s.Options.AddFlags(fss.FlagSet("assignment"), s.Options)
	return fss
}

func (c *RunServerConfig) PreRun() circle.Source {
	return circle.Source{
		Account:         c.Account,
		Password:        c.Password,
		Tuisongclientid: c.Tuisongclientid,
	}
}
