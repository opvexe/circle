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
	svc  circle.Service
	wc   chan WC
	stop chan struct{}
}

func NewAssignment() circle.Assignment {
	s := &assignment{
		svc:  NewService(),
		wc:   make(chan WC, 100),
		stop: make(chan struct{}),
	}
	go s.schedule()
	return s
}

func (s *assignment) schedule() {
	sched := schedule.NewFIFOScheduler()
	for {
		select {
		case list := <-s.wc:
			for _, l := range list.wc {
				sched.Schedule(func(ctx context.Context) {
					if err := s.svc.Do(ctx, l, list.token); err != nil {
						log.Printf("sched do task err: %s", err)
					}
				})
			}
			sched.WaitFinish(len(list.wc))
			return
		case <-s.stop:
			if sched != nil {
				sched.Stop()
			}
			return
		}
	}
}

type WC struct {
	wc    circle.WechatShares
	token string
}

func (s *assignment) Pub(source circle.Source) error {
	list, token, err := s.svc.List(context.Background(), source)
	if err != nil {
		return err
	}

	select {
	case s.wc <- WC{wc: list, token: token}:
	default:
	}
	return nil
}

func (s *assignment) Close() error {
	if s.stop != nil {
		close(s.stop)
	}
	return nil
}
