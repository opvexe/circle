/*
Copyright 2021 The Xiadat Authors.

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

package circle

import "context"

const URL = "https://informationworld.zdjt.com/index.php/api/mytask/"

// Query processing parameters.
type Query struct {
	URL     string
	Page    int
	Version string
}

type Task struct {
	Title           string `json:"title"`
	EndTime         int    `json:"end_time"`
	ReadCount       int    `json:"read_count"`
	FinishScore     int    `json:"finish_score"`
	MicroURL        string `json:"micro_url"`
	SltImg          string `json:"slt_img"`

	CircleCountRw   int    `json:"circle_count_rw"`
	GroupCountRw    int    `json:"group_count_rw"`
	ID              int    `json:"id"`
	UserTaskCount   int    `json:"user_task_count"`
	UserReadCount   int    `json:"user_read_count"`
	UserScore       int    `json:"user_score"`
	UserCircleCount int    `json:"user_circle_count"`
	UserGroupCount  int    `json:"user_group_count"`
}

// Tasker is an interface for gathering, dispose task.
type Tasker interface {
	Fetch(ctx context.Context, query Query) ([]Task, error)
}
