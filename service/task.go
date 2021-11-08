/*
Copyright 2021 The Gridsum Authors.

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
	"errors"
	"fmt"
	dottask "github.com/devfeel/dottask"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"sync"
	"x6t.io/circle"
)

const TaskName = "wechat"

type TaskService struct {
	taskSvc *dottask.TaskService
	source  circle.Source
	wg      errgroup.Group
	config   circle.Config
}

func NewTaskService(config circle.Config) *TaskService {
	return &TaskService{
		taskSvc: dottask.StartNewService(),
		source:  circle.Source{
			URL:             fmt.Sprintf("%s%s", URL, login),
			Account:         config.Account,
			Password:        config.Password,
			Tuisongclientid: config.Tuisongclientid,
		},
		config: config,
	}
}

func (s *TaskService) Open(ctx context.Context) error {
	_, err := s.taskSvc.CreateCronTask(TaskName, true, s.config.Express, s.Task, nil)
	if err != nil {
		return err
	}

	defer s.taskSvc.StartAllTask()
	return nil
}

func (s *TaskService) Task(ctx *dottask.TaskContext) error {
	userInfo, err := s.Login(context.Background(), s.source)
	if err != nil {
		return err
	}
	tasks, err := s.StatisticsTask(context.Background(), userInfo.Token)
	if err != nil {
		return err
	}

	return s.ProcessTask(context.Background(), tasks, userInfo.Token)
}

func (s *TaskService) Login(ctx context.Context, u circle.Source) (*circle.UserInfo, error) {
	user, err := NewUserService().New(s.source)
	if err != nil {
		return nil, err
	}
	userInfo, err := user.Login(context.Background(), s.source)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *TaskService) StatisticsTask(ctx context.Context, token string) (circle.Tasks, error) {
	fetcher, err := NewFetchService().New(circle.Source{
		URL:   fmt.Sprintf("%s%s", URL, tasks),
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	var (
		tasks circle.Tasks
		lock  sync.RWMutex
	)
	for j := 1; j < circle.Offset; j++ {
		i := j
		s.wg.Go(func() error {
			task, err := fetcher.Fetch(ctx, circle.Fitter{
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

	if len(tasks) == 0 {
		return nil, errors.New("tasks is empty")
	}

	// analyze unfinished tasks.
	return s.AnalyzeUnfinishedTasks(ctx, tasks)
}

func (s *TaskService) AnalyzeUnfinishedTasks(ctx context.Context, tasks circle.Tasks) (circle.Tasks, error) {
	var unfinished circle.Tasks
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
	return unfinished, nil
}

func (s *TaskService) ProcessTask(ctx context.Context, tasks circle.Tasks, token string) error {
	share, err := NewShareService().New(circle.Source{
		URL:   fmt.Sprintf("%s%s", URL, wechat),
		Token: token,
	})
	if err != nil {
		return err
	}
	for _, task := range tasks {
		if task.UserCircleCount < task.CircleCountRw { // 朋友圈
			for i := 0; i < task.CircleCountRw-task.UserCircleCount; i++ {
				_ = share.Wechat(ctx, circle.WechatShare{
					Microgrid: cast.ToString(task.ID),
					Type:      circle.Friends,
				})
			}
		}
		if task.UserGroupCount < task.GroupCountRw { // 微信群
			for i := 0; i < task.GroupCountRw-task.UserGroupCount; i++ {
				_ = share.Wechat(ctx, circle.WechatShare{
					Microgrid: cast.ToString(task.ID),
					Type:      circle.Group,
				})
			}
		}

		if task.UserReadCount < task.ReadCount { // 阅读

		}
	}
	return nil
}

func (s *TaskService) Close(ctx context.Context) error {
	s.taskSvc.StopAllTask()

	if err := s.wg.Wait(); err != nil {
		return err
	}
	return nil
}
