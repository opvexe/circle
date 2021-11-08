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

package main

import (
	"context"
	"github.com/vrischmann/envconfig"
	"time"
	"x6t.io/circle"
	"x6t.io/circle/service"
)

func main() {
	config := circle.NewConfig()
	if err:=envconfig.InitWithPrefix(config, "CCL");err!=nil{
		circle.CheckErr(err)
	}
	if err:=cmdRun(context.Background(),*config);err!=nil{
		circle.CheckErr(err)
	}
}

func cmdRun(ctx context.Context, o circle.Config) error {

	// Start the launcher and wait for it to exit on SIGINT or SIGTERM.
	runCtx := circle.WithStandardSignals(ctx)
	svc :=service.NewTaskService(o)
	if err:=svc.Open(runCtx);err!=nil{
		return err
	}

	<-runCtx.Done()

	// Tear down the launcher, allowing it a few seconds to finish any
	// in-progress requests.
	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return svc.Close(shutdownCtx)
}
