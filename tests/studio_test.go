/*
 * Copyright (c) 2023 The nebula-contrib Authors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package tests

import (
	"context"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"

	"github.com/nebula-contrib/ngctl/cmd"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

func TestStudio(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("studio install", func(t *testing.T) {
		// consider test in local environment with kubeconfig file in default path ~/.kube/config
		var namespace = "default"
		var name = "studio-test"
		set, err := util.NewClientSet(filepath.Join(homedir.HomeDir(), ".kube", "config"))
		if err != nil {
			t.Errorf("create client set error: %v", err)
		}
		ctx := context.Background()

		// run studio install command
		command.SetArgs([]string{"studio", "install", "--name", name})
		err = command.Execute()
		if err != nil {
			t.Errorf("run studio command error: %v", err)
		}
		// check whether the deployment is created
		deploy, err := set.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Errorf("get deployment error: %v", err)
		}
		if deploy.Name != name {
			t.Errorf("expect deployment name is %s, but got %s", name, deploy.Name)
		}
		//	 check whether the service is created
		svc, err := set.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Errorf("get service error: %v", err)
		}
		if svc.Name != name {
			t.Errorf("expect service name is %s, but got %s", name, svc.Name)
		}
	})
	t.Run("studio uninstall", func(t *testing.T) {
		var name = "studio-test"

		//   run studio uninstall command
		command.SetArgs([]string{"studio", "uninstall", "--name", name})
		err := command.Execute()
		if err != nil {
			t.Errorf("run studio command error: %v", err)
		}
	})
}
