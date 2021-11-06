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
	"x6t.io/circle"
	"x6t.io/circle/queries"
)

const URL = "https://informationworld.zdjt.com/index.php/api/"

const (
	login   = "user/login"          // 登陆
	tasks   = "mytask/wwtask"       // 任务列表
	detail  = "mytask/microdetail"  // 任务详情
	friends = "mytask/forwardlogww" // 朋友圈
)

type UserService struct{}

func (s *UserService) New(src circle.Source) (circle.User, error) {
	client := &queries.Client{}

	if err := client.Connect(context.TODO(), &src); err != nil {
		return nil, err
	}
	return client, nil
}

type FetchService struct{}

func (s *FetchService) New(src circle.Source) (circle.Fetcher, error) {
	client := &queries.Client{}

	if err := client.Connect(context.TODO(), &src); err != nil {
		return nil, err
	}
	return client, nil
}

type ShareService struct {}

func (s *ShareService) New(src circle.Source) (circle.Share, error) {
	client := &queries.Client{}

	if err := client.Connect(context.TODO(), &src); err != nil {
		return nil, err
	}
	return client, nil
}