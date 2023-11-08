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
	"log"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nebula-contrib/ngctl/pkg/util"
	"github.com/nebula-contrib/ngctl/pkg/version"
)

const OperatorSelector = "app.kubernetes.io/component=controller-manager,app.kubernetes.io/instance=nebula-operator"

func versionCmd() *cobra.Command {
	var (
		namespace string
	)
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show the version of ngctl and nebula operator",
		Long:  "show the version of ngctl and nebula operator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ngctlVersion(namespace)
		},
	}
	cmd.PersistentFlags().StringVar(&namespace, "operator-namespace", "nebula-operator-system", "namespace of nebula operator")
	return cmd
}

func ngctlVersion(namespace string) error {
	log.Printf("ngctl Version: %s", version.GetVersion())
	set, _, err := util.NewClientSet(kubeConfig, kubeContext)
	if err != nil {
		return err
	}
	ctx := context.Background()

	controllers, err := set.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: OperatorSelector,
	})
	if err != nil {
		return err
	}
	if len(controllers.Items) == 0 {
		log.Printf("nebula operator is not installed")
		return nil
	}
	log.Printf("Nebula Operator Version: %s", controllers.Items[0].Spec.Template.Spec.Containers[0].Image)
	return nil
}
