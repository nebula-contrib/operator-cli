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
	"errors"
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-operator/apis/apps/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nebula-contrib/ngctl/pkg/config"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

func useCmd() *cobra.Command {
	var (
		namespace string
	)
	cmd := &cobra.Command{
		Use:   "use",
		Short: "Specify a Nebula Graph cluster to use",
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCluster(args, namespace)
		},
	}
	cmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "namespace of the nebula graph cluster")
	return cmd
}

func useCluster(args []string, namespace string) error {
	if len(args) == 0 {
		return errors.New("please specify the name of Nebula Graph cluster")
	}
	name := args[0]

	client, _, err := util.NewDynamicClient(kubeConfig, kubeContext)
	if err != nil {
		return err
	}

	resource := v1alpha1.GroupVersion.WithResource("nebulaclusters")

	ctx := context.Background()
	_, err = client.Resource(resource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	err = config.SaveConfig(namespace, name)
	if err != nil {
		return err
	}
	log.Printf("use nebula graph cluster %s in namespace %s", name, namespace)
	return nil
}
