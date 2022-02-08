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
	"fmt"
	"x6t.io/circle"
)

type RunServerConfig struct {
	Account         string `envconfig:"optional"`
	Password        string `envconfig:"optional"`
	Tuisongclientid string `envconfig:"optional"`
	Express         string `envconfig:"optional"`
}

func NewRunServerConfig() *RunServerConfig {
	return &RunServerConfig{
		Account:         "chenxue",
		Password:        "ZHENdong123",
		Tuisongclientid: "e0d0171b89075356632758ca7df6a3ac",
		Express:         "@daily",
	}
}

func (c *RunServerConfig) Validate() error {
	if len(c.Account) == 0 {
		return fmt.Errorf("account must be not empty")
	}
	if len(c.Password) == 0 {
		return fmt.Errorf("password must be not empty")
	}
	if len(c.Tuisongclientid) == 0 {
		return fmt.Errorf("clientid must be not empty")
	}
	return nil
}

func (c *RunServerConfig) PreRun() circle.Source {
	return circle.Source{
		Account:         c.Account,
		Password:        c.Password,
		Tuisongclientid: c.Tuisongclientid,
	}
}
