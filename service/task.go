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
}

func NewTaskService(source circle.Source) *TaskService {
	source = circle.Source{
		URL:             fmt.Sprintf("%s%s", URL, login),
		Account:         "chenxue",
		Password:        "ZHENdong123",
		Tuisongclientid: "e0d0171b89075356632758ca7df6a3ac",
	}
	return &TaskService{
		taskSvc: dottask.StartNewService(),
		source:  source,
	}
}

func (s *TaskService) Open() error {
	_, err := s.taskSvc.CreateCronTask(TaskName, true, "", s.Task, nil)
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

	return nil
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

	return nil, nil
}

func (s *TaskService) ProcessTask(ctx context.Context) error {

	return nil
}

func (s *TaskService) Close() error {
	defer s.taskSvc.StopAllTask()

	return nil
}
