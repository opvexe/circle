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
	"testing"
	"x6t.io/circle"
)

func TestNewAssignment(t *testing.T) {
	assig := NewAssignment()
	defer assig.Close()

	if err := assig.Pub(circle.Source{
		Account:         "chenxue",
		Password:        "ZHENdong123",
		Tuisongclientid: "e0d0171b89075356632758ca7df6a3ac",
	}); err != nil {
		t.Errorf("pub assignment error :%s", err)
	}
}
