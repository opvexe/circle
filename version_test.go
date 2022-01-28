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

package circle

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want Info
	}{
		{
			name: "version",
			want: Info{
				GitVersion: gitVersion,
				BuildDate:  time.Now().Format("2006-01-02 15:04:05"),
				GoVersion:  runtime.Version(),
				Compiler:   runtime.Compiler,
				Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
