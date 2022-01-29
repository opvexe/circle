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
	"gitlab.com/taoshumin/go-schedule/schedule"
	"log"
	"x6t.io/circle"
)

type assignment struct {
	svc    circle.Service
	sched  schedule.Scheduler
	ctx    context.Context
	source circle.Source
}

func NewAssignment(source circle.Source, ctx context.Context) *assignment {
	s := &assignment{
		svc:    NewService(),
		sched:  schedule.NewFIFOScheduler(),
		source: source,
		ctx:    ctx,
	}
	go s.run()
	return s
}

func (s *assignment) run() error {
	defer s.sched.Stop()
	tasks, token, err := s.svc.Get(s.ctx, s.source)
	if err != nil {
		return err
	}

	shares, err := s.svc.UnfinishedWechatShares(tasks)
	if err != nil {
		return err
	}

	jb := func(ws circle.WechatShare) schedule.Job {
		return func(ctx context.Context) {
			if err := s.svc.Do(ctx, ws, token); err != nil {
				log.Printf("sched do task err: %s", err)
			}
		}
	}

	for _, share := range shares {
		s.sched.Schedule(jb(share))
	}

	s.sched.WaitFinish(len(shares))
	return nil
}

func (s *assignment) Close() error {
	if s.sched != nil {
		s.sched.Stop()
	}
	return nil
}
