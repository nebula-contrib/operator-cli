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
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/nebula-contrib/ngctl/pkg/config"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

const (
	graphd   = "graphd"
	metad    = "metad"
	storaged = "storaged"
	volume   = "volume"
)

func getCmd() *cobra.Command {
	var (
		allNamespaces bool
	)
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get component of nebula graph cluster",
		Long:  "get component of nebula graph cluster.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(args, allNamespaces)
		},
	}
	cmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "if set, list the nebula graph clusters across all namespaces")
	return cmd
}

func get(args []string, allNamespaces bool) error {
	if len(args) == 0 {
		return errors.New("please specify the kind of component")
	}
	kind := args[0]
	ctx := context.Background()
	client, err := util.NewClientSet(kubeConfig)
	if err != nil {
		return err
	}

	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}
	name, namespace := conf.Name, conf.Namespace

	switch kind {
	case graphd, metad, storaged:
		return getComponents(ctx, client, kind, name, namespace, allNamespaces)
	case volume:
		return getVolumes(ctx, client, name, namespace, allNamespaces)
	default:
		return errors.New("unsupported kind type of get command")
	}
}

func getComponents(ctx context.Context, client *kubernetes.Clientset, kind, name, namespace string, allNamespace bool) error {
	podList, err := getComponentPods(ctx, client, kind, name, namespace, allNamespace)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		return nil
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"NAME", "READY", "STATUS", "MEMORY", "CPU", "RESTARTS", "AGE", "NODE"})
	for _, pod := range podList.Items {
		t.AppendRow(table.Row{
			pod.Name,
			pod.Status.ContainerStatuses[0].Ready,
			pod.Status.Phase,
			pod.Spec.Containers[0].Resources.Requests.Memory().String(),
			pod.Spec.Containers[0].Resources.Requests.Cpu().String(),
			pod.Status.ContainerStatuses[0].RestartCount,
			// Age
			time.Since(pod.CreationTimestamp.Time).String(),
			//	HostIp
			pod.Status.HostIP,
		},
		)
	}
	t.Render()
	return nil
}

func getComponentPods(ctx context.Context, client *kubernetes.Clientset, kind, name string, namespace string, allNamespace bool) (*corev1.PodList, error) {
	if allNamespace {
		return client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
			// ignore namespace, kind, name
			LabelSelector: "app.kubernetes.io/name=nebula-graph,app.kubernetes.io/component",
		})
	}
	var selector string
	switch kind {
	case graphd, metad, storaged:
		selector = "app.kubernetes.io/cluster=" + name + ",app.kubernetes.io/component=" + kind + ",app.kubernetes.io/name=nebula-graph"
	default:
		selector = "app.kubernetes.io/cluster=" + name + ",app.kubernetes.io/name=nebula-graph"
	}

	return client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
}

func getPersistentVolume(ctx context.Context, client *kubernetes.Clientset, name string, namespace string, allNamespace bool) (*corev1.PersistentVolumeList, error) {
	var selector string
	if allNamespace {
		selector = "app.kubernetes.io/name=nebula-graph"
	} else {
		selector = "app.kubernetes.io/cluster=" + name + ",app.kubernetes.io/name=nebula-graph"
	}

	pvs, err := client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	var list corev1.PersistentVolumeList
	for _, pv := range pvs.Items {
		if pv.Spec.ClaimRef != nil && (allNamespace || pv.Spec.ClaimRef.Namespace == namespace) {
			list.Items = append(list.Items, pv)
		}
	}

	return &list, err
}

func getVolumes(ctx context.Context, client *kubernetes.Clientset, name string, namespace string, allNamespace bool) error {
	pvs, err := getPersistentVolume(ctx, client, name, namespace, allNamespace)
	if err != nil {
		return err
	}

	podList, err := getComponentPods(ctx, client, "", name, namespace, allNamespace)
	if err != nil {
		return err
	}

	podMap := make(map[string]string)
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				key := fmt.Sprintf("%s/%s", namespace, volume.PersistentVolumeClaim.ClaimName)
				// Node IP
				podMap[key] = pod.Status.HostIP
			}
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"VOLUME", "CLAIM", "STATUS", "CAPACITY", "HOST IP"})

	for _, pv := range pvs.Items {
		t.AppendRow(table.Row{
			pv.Name,
			// Claim
			pv.Spec.ClaimRef.Name,
			// Status
			pv.Status.Phase,
			// Capacity
			pv.Spec.Capacity.Storage().String(),
			// Host IP
			podMap[fmt.Sprintf("%s/%s", namespace, pv.Spec.ClaimRef.Name)],
		})
	}

	t.Render()
	return nil
}
