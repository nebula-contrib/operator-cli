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

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/nebula-contrib/ngctl/pkg/studio"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

const studioLabelKey = "ngctl/nebula-studio"

func studioCmd() *cobra.Command {
	var (
		name      string
		namespace string
		nodePort  int32
		image     string
	)
	cmd := &cobra.Command{
		Use:   "studio",
		Short: "the command line tool for nebula graph studio",
		Long:  "studio is a command line tool for nebula graph studio.",
		Example: `  # installStudio nebula graph studio
  ngctl studio installStudio --name studio --nodePort 30180
  # uninstall nebula graph studio
  ngctl studio uninstall --name studio
`,
	}
	cmd.PersistentFlags().StringVar(&name, "name", "studio", "name of the nebula graph studio")
	cmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "namespace of the nebula graph studio")

	install := cobra.Command{
		Use:   "install",
		Short: "install nebula graph studio",
		Long:  "install nebula graph studio.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installStudio(name, namespace, image, nodePort)
		},
	}
	install.PersistentFlags().Int32Var(&nodePort, "nodePort", 30180, "nodePort of the nebula graph studio")
	install.PersistentFlags().StringVar(&image, "image", "vesoft/nebula-graph-studio:v3.7.0", "image of the nebula graph studio")

	cmd.AddCommand(&install)

	uninstall := cobra.Command{
		Use:   "uninstall",
		Short: "uninstall nebula graph studio",
		Long:  "uninstall nebula graph studio.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstallStudio(name, namespace)
		},
	}
	cmd.AddCommand(&uninstall)

	return cmd
}

func installStudio(name, namespace string, image string, nodePort int32) error {
	clientSet, err := util.NewClientSet(kubeConfig)
	if err != nil {
		return err
	}
	labels := map[string]string{
		studioLabelKey: name,
	}
	var ctx = context.Background()

	deploy := studio.CreateDeployment(name, namespace, labels, 1, image)
	deployInterface := clientSet.AppsV1().Deployments(namespace)
	err = util.CreateResource[appsv1.Deployment](ctx, deployInterface, deploy, name, studioLabelKey)
	if err != nil {
		return err
	}

	svc := studio.CreateService(name, namespace, labels, nodePort)
	svcInterface := clientSet.CoreV1().Services(namespace)
	err = util.CreateResource[corev1.Service](ctx, svcInterface, svc, name, studioLabelKey)
	if err != nil {
		return err
	}
	return nil
}

func uninstallStudio(name, namespace string) error {
	clientSet, err := util.NewClientSet(kubeConfig)
	if err != nil {
		return err
	}
	var ctx = context.Background()

	// remove deploy
	deployments := clientSet.AppsV1().Deployments(namespace)
	err = util.RemoveResource[appsv1.Deployment](ctx, deployments, name, studioLabelKey)
	if err != nil {
		return err
	}
	// remove service
	services := clientSet.CoreV1().Services(namespace)
	err = util.RemoveResource[corev1.Service](ctx, services, name, studioLabelKey)
	if err != nil {
		return err
	}
	return nil
}
