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
	dottask "github.com/devfeel/dottask"
	"testing"
	"x6t.io/circle"
)

func TestTaskService_Task(t *testing.T) {

	tests := []struct {
		name    string
		fields  *TaskService
		wantErr bool
	}{
		{
			name:    "task",
			fields:  NewTaskService(*circle.NewConfig()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.fields.Task(&dottask.TaskContext{}); (err != nil) != tt.wantErr {
				t.Errorf("Task() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
