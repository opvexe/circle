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

package service

import (
	"context"
	"fmt"
	"testing"
	"x6t.io/circle"
)

func TestUserService_New_Login(t *testing.T) {
	type args struct {
		src circle.Source
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "service-login",
			args: args{
				src: circle.Source{
					URL:             fmt.Sprintf("%s%s", URL, login),
					Account:         "chenxue",
					Password:        "ZHENdong123",
					Tuisongclientid: "e0d0171b89075356632758ca7df6a3ac",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{}
			got, err := s.New(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			resp, err := got.Login(context.TODO(), tt.args.src)
			if err != nil {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Errorf("New() got = %v", got)
			}
		})
	}
}

func TestUserService_New_Task_List(t *testing.T) {
	type args struct {
		src circle.Source
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "service-fetch",
			args: args{
				src: circle.Source{
					URL: fmt.Sprintf("%s%s", URL, tasks),
					Token: "eb470e70-14dc-4c8b-8e89-60d3e7aba278",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FetchService{}
			got, err := s.New(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			resp, err := got.Fetch(context.TODO(),circle.Fitter{Page: "1",Version: "1.0.4"})
			if err != nil {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Errorf("New() got = %v", got)
			}
		})
	}
}

func TestUserService_New_Task_Detail(t *testing.T) {
	type args struct {
		src circle.Source
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "service-fetch",
			args: args{
				src: circle.Source{
					URL: fmt.Sprintf("%s%s", URL, detail),
					Token: "eb470e70-14dc-4c8b-8e89-60d3e7aba278",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FetchService{}
			got, err := s.New(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			resp, err := got.Detail(context.TODO(),"840")
			if err != nil {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Errorf("New() got = %v", got)
			}
		})
	}
}
