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

// Query is processing parameters.
type Query struct {
	URL     string
	Page    int
	Version string
}

type Task struct {
	ID      int    `json:"id,omitempty"`       // 文章id
	Title   string `json:"title,omitempty"`    // 标题
	EndTime int    `json:"end_time,omitempty"` // 结束时间
	// The current number of tasks and the specified number of completed tasks.
	// containing user is the statistics of the number of tasks completed by the user prefix.
	ReadCount       int `json:"read_count,omitempty"`        // 任务阅读数
	UserReadCount   int `json:"user_read_count,omitempty"`   // 用户阅读数
	FinishScore     int `json:"finish_score,omitempty"`      // 任务积分
	UserScore       int `json:"user_score,omitempty"`        // 用户当前积分
	GroupCountRw    int `json:"group_count_rw,omitempty"`    // 任务微信群数
	UserGroupCount  int `json:"user_group_count,omitempty"`  // 用户微信群数
	CircleCountRw   int `json:"circle_count_rw,omitempty"`   // 任务朋友圈
	UserCircleCount int `json:"user_circle_count,omitempty"` // 用户朋友圈
	// Statistics on the number of user sharing tasks.
	// after adding the number of WeChat groups to the number of Moments.
	UserTaskCount int `json:"user_task_count,omitempty"` // 用户微信群+朋友圈分享次数
	// WeChat share pictures and connections.
	MicroURL string `json:"micro_url,omitempty"` // 微信分享链接
}

type Tasks []Task

// User is an interface for login.
type User interface {
	Login(ctx context.Context, username, password string) (string, error)
}

// Fetcher is an interface for fetch task list.
type Fetcher interface {
	Fetch(ctx context.Context, query Query) (Tasks, error)
}

// Share is an interface for share articles to wechat groups and friends.
type Share interface {
	WechatFriends(ctx context.Context, tasks Tasks) error
	WechatGroup(ctx context.Context, tasks Tasks) error
}

// Tasker is an interface for dispose task.
type Tasker interface {
	Share
	SharedByOtherRead(ctx context.Context) error
}

// Client is an interface for handle http request.
type Client interface {
	Get(ctx context.Context) error
	Post(ctx context.Context) error
}
