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
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

const gitVersion = "v1.0.0"

type Info struct {
	GitVersion string
	BuildDate  string
	GoVersion  string
	Compiler   string
	Platform   string
}

func (info Info) String() string {
	jsonString, _ := json.Marshal(info)
	return string(jsonString)
}

func Get() Info {
	return Info{
		GitVersion: gitVersion,
		BuildDate:  time.Now().Format("2006-01-02 15:04:05"),
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
