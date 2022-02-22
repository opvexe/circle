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

package services

import (
	"fmt"
	"github.com/spf13/pflag"
)

type Options struct {
	Account         string `json:",omitempty" yaml:"account" description:"Circle service host account"`
	Password        string `json:",omitempty" yaml:"password" description:"Circle service password"`
	Tuisongclientid string `json:",omitempty" yaml:"tuisongclientid" description:"Circle service Tuisong client id"`
	Express         string `json:",omitempty" yaml:"express" description:"Circle service express"`
}

func NewOptions() *Options {
	return &Options{
		Account:         "zhangming",
		Password:        "zhangming",
		Tuisongclientid: "e0d0171b8887997758ca7df6a3ac",
		Express:         "@daily",
	}
}

func (c *Options) Validate() []error {
	var errors []error

	if len(c.Account) == 0 {
		errors = append(errors, fmt.Errorf("account must be not empty"))
	}
	if len(c.Password) == 0 {
		errors = append(errors, fmt.Errorf("password must be not empty"))
	}
	if len(c.Tuisongclientid) == 0 {
		errors = append(errors, fmt.Errorf("clientid must be not empty"))
	}
	if len(c.Express) == 0 {
		errors = append(errors, fmt.Errorf("express must be not empty"))
	}
	return errors
}

func (s *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	fs.StringVar(&s.Account, "account", c.Account, "Circle service account name.")
	fs.StringVar(&s.Password, "password", c.Password, "Circle service account password")
	fs.StringVar(&s.Tuisongclientid, "tuisongclientid", c.Tuisongclientid, "Circle service Tuisong client id")
	fs.StringVar(&s.Express, "express", c.Express, "Circle service express")
}
