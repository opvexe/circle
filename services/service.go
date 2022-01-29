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
	"context"
	"fmt"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"sync"
	"x6t.io/circle"
	"x6t.io/circle/queries"
)

const (
	URL    = "https://informationworld.zdjt.com/index.php/api/"
	login  = "user/login"
	tasks  = "mytask/wwtask"       // 任务列表
	detail = "mytask/microdetail"  // 任务详情
	wechat = "mytask/forwardlogww" // 1:朋友圈 2:微信群
)

type service struct {
	connect *Connect
	client  circle.Client
	wg      errgroup.Group
}

func NewService() circle.Service {
	return &service{
		connect: NewConnectService(),
		client:  &queries.Client{},
		wg:      errgroup.Group{},
	}
}

func (s *service) Get(ctx context.Context, source circle.Source) (circle.Tasks, string, error) {
	userInfo, err := s.login(ctx, source)
	if err != nil {
		return nil, "", err
	}

	tasks, err := s.fetch(ctx, userInfo.Token)
	if err != nil {
		return nil, userInfo.Token, err
	}
	return tasks, userInfo.Token, nil
}

func (s *service) login(ctx context.Context, source circle.Source) (*circle.UserInfo, error) {
	source.URL = fmt.Sprintf("%s%s", URL, login)
	client, err := s.connect.Connect(ctx, source)
	if err != nil {
		return nil, err
	}
	userInfo, err := client.Login(ctx, source)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *service) fetch(ctx context.Context, token string) (circle.Tasks, error) {
	url := fmt.Sprintf("%s%s", URL, tasks)
	client, err := s.connect.Connect(ctx, circle.Source{URL: url, Token: token})
	if err != nil {
		return nil, err
	}

	var (
		tasks circle.Tasks
		lock  sync.Mutex
	)

	for j := 1; j < circle.Offset; j++ {
		i := j
		s.wg.Go(func() error {
			task, err := client.Fetch(ctx, circle.Fitter{
				Page:    cast.ToString(i),
				Version: circle.Version,
			})
			if err != nil {
				return err
			}
			lock.Lock()
			tasks = append(tasks, task...)
			lock.Unlock()
			return nil
		})
	}

	if err := s.wg.Wait(); err != nil {
		return nil, err
	}
	return s.unfinished(tasks), nil
}

func (s *service) unfinished(tasks circle.Tasks) (unfinished circle.Tasks) {
	if len(tasks) == 0 {
		return unfinished
	}

	for _, task := range tasks {
		if task.UserCircleCount < task.CircleCountRw {
			unfinished = append(unfinished, task)
			continue
		}
		if task.UserGroupCount < task.GroupCountRw {
			unfinished = append(unfinished, task)
			continue
		}
		if task.UserReadCount < task.ReadCount {
			unfinished = append(unfinished, task)
			continue
		}
	}
	return unfinished
}

func (s *service) UnfinishedWechatShares(tasks circle.Tasks) (circle.WechatShares, error) {
	shares := make(circle.WechatShares, len(tasks))
	for _, task := range tasks {
		if task.UserCircleCount < task.CircleCountRw { // 朋友圈
			for i := 0; i < task.CircleCountRw-task.UserCircleCount; i++ {
				shares = append(shares, circle.WechatShare{
					Microgrid: cast.ToString(task.ID),
					Type:      circle.Friends,
				})
			}
		}

		if task.UserGroupCount < task.GroupCountRw { // 微信群
			for i := 0; i < task.GroupCountRw-task.UserGroupCount; i++ {
				shares = append(shares, circle.WechatShare{
					Microgrid: cast.ToString(task.ID),
					Type:      circle.Group,
				})
			}
		}

		// number of unread WeChat articles.
		if task.UserReadCount < task.ReadCount {
		}
	}

	if len(shares) == 0 {
		return nil, fmt.Errorf("no WeChat sharing task list")
	}
	return shares, nil
}

func (s *service) List(ctx context.Context, source circle.Source) (circle.WechatShares, string, error) {
	tasks, token, err := s.Get(context.Background(), source)
	if err != nil {
		return nil, "", err
	}

	shares, err := s.UnfinishedWechatShares(tasks)
	if err != nil {
		return nil, "", err
	}
	return shares, token, nil
}

func (s *service) Do(ctx context.Context, wc circle.WechatShare, token string) error {
	url := fmt.Sprintf("%s%s", URL, wechat)
	client, err := s.connect.Connect(ctx, circle.Source{URL: url, Token: token})
	if err != nil {
		return err
	}
	return client.Wechat(ctx, wc)
}
