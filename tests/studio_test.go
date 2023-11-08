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
	"fmt"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"

	"github.com/nebula-contrib/ngctl/cmd"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

/*

Assuming the Kubernetes cluster is created by Minikube with default settings:
- The kubeconfig file is located at "~/.kube/config".
- The context defined in ~/.kube/config is set to "minikube".

Additionally, assuming the file example.yaml defines the Nebula cluster:
- The Nebula cluster is named "nebula"

*/

func TestStudio(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("studio install", func(t *testing.T) {
		// consider test in local environment with kubeconfig file in default path ~/.kube/config
		var namespace = "default"
		var name = "studio-test"
		set, _, err := util.NewClientSet(filepath.Join(homedir.HomeDir(), ".kube", "config"), "")
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

func TestVersion(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("version", func(t *testing.T) {
		// run version command
		command.SetArgs([]string{"version"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run version command error: %v", err)
		}
	})
}

func TestUse(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("use", func(t *testing.T) {
		// run use command
		command.SetArgs([]string{"use", "nebula"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run use command error: %v", err)
		}
	})
}

func TestInfo(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("info", func(t *testing.T) {
		// run info command
		command.SetArgs([]string{"info"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run info command error: %v", err)
		}
	})
}

func TestList(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("list", func(t *testing.T) {
		// run list command
		command.SetArgs([]string{"list"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run list command error: %v", err)
		}
	})

	t.Run("list -A", func(t *testing.T) {
		// run list command with -A flag
		command.SetArgs([]string{"list", "-A"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run list command error: %v", err)
		}
	})

	t.Run("list --context", func(t *testing.T) {
		// assume the context is "minikube"
		command.SetArgs([]string{"list", "--context", "minikube"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run list command error: %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	var command = cmd.RootCmd
	args := []string{"metad", "storaged", "graphd", "volume"}
	for _, v := range args {
		t.Run(fmt.Sprintf("get %s", v), func(t *testing.T) {
			command.SetArgs([]string{"get", v})
			err := command.Execute()
			if err != nil {
				t.Errorf("run get %s command error: %v", v, err)
			}
		})
	}

}

func TestConsole(t *testing.T) {
	var command = cmd.RootCmd
	t.Run("console", func(t *testing.T) {
		// run console command
		command.SetArgs([]string{"console", "-u", "root", "-p", "nebula"})
		err := command.Execute()
		if err != nil {
			t.Errorf("run console command error: %v", err)
		}
	})
}
