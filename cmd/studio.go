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
	"github.com/nebula-contrib/operator-cli/internal/studio"
	"github.com/nebula-contrib/operator-cli/internal/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

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
		RunE: func(cmd *cobra.Command, args []string) error {
			clientSet, err := util.NewClientSet(kubeConfig)
			if err != nil {
				return err
			}
			labels := map[string]string{
				"nebula-studio/operator-cli": name,
			}

			deployment := studio.CreateDeployment(name, labels, 1, image)

			var ctx = context.Background()
			deployment, err = clientSet.AppsV1().
				Deployments(namespace).
				Create(ctx, deployment, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			log.Printf("Deployment created: %s\n", deployment.Name)
			service := studio.CreateService(name, labels, nodePort)
			service, err = clientSet.CoreV1().
				Services(namespace).
				Create(ctx, service, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			log.Printf("Service created: %s\n", service.Name)
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&name, "name", "studio", "name of the nebula graph studio")
	cmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "namespace of the nebula graph studio")
	cmd.PersistentFlags().Int32Var(&nodePort, "nodePort", 30080, "nodePort of the nebula graph studio")
	cmd.PersistentFlags().StringVar(&image, "image", "vesoft/nebula-graph-studio:v3.7.0", "image of the nebula graph studio")
	return cmd
}
