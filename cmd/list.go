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
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/nebula-contrib/ngctl/pkg/list"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

func listCmd() *cobra.Command {
	var (
		namespace     string
		allNamespaces bool
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all installed nebula graph clusters",
		Long:  "list all installed nebula graph clusters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listClusters(allNamespaces, namespace)
		},
	}
	cmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "namespace of the nebula graph cluster")
	cmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "if set, list the nebula graph clusters across all namespaces")
	return cmd
}

func listClusters(allNamespaces bool, namespace string) error {
	client, err := util.NewDynamicClient(kubeConfig)
	if err != nil {
		return err
	}
	ctx := context.Background()
	clusters, err := list.Clusters(ctx, client, allNamespaces, namespace)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		log.Printf("no nebula graph cluster found in namespace %s", namespace)
		return nil
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Namespace", "Name", "Graphd", "Metad", "Storaged"})
	for _, cluster := range clusters {
		status := cluster.Status
		spec := cluster.Spec
		graphdReady := fmt.Sprintf("%d/%d", status.Graphd.Workload.ReadyReplicas, *spec.Graphd.Replicas)
		metadReady := fmt.Sprintf("%d/%d", status.Metad.Workload.ReadyReplicas, *spec.Metad.Replicas)
		storagedReady := fmt.Sprintf("%d/%d", status.Storaged.Workload.ReadyReplicas, *spec.Storaged.Replicas)
		t.AppendRow(table.Row{cluster.Namespace, cluster.Name, graphdReady, metadReady, storagedReady})
	}
	t.Render()
	return nil
}
