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
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-operator/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/nebula-contrib/ngctl/pkg/config"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

const serviceSelector = "app.kubernetes.io/cluster=%s,app.kubernetes.io/component=%s,app.kubernetes.io/name=nebula-graph"

func infoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "information of nebula graph clusters",
		Long:  "information of nebula graph clusters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return info()
		},
	}
	return cmd
}

func info() error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}
	namespace, name := conf.Namespace, conf.Name

	client, _, err := util.NewDynamicClient(kubeConfig, kubeContext)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cluster, err := getCluster(ctx, client, name, namespace)
	if err != nil {
		return err
	}
	clusterInfo(cluster)
	log.Println("Overview:")
	componentInfo(cluster)
	clientSet, _, err := util.NewClientSet(kubeConfig, kubeContext)
	if err != nil {
		return err
	}
	log.Println("Endpoints:")
	err = endpointsInfo(ctx, clientSet, name, namespace)
	if err != nil {
		return err
	}

	return nil
}

func endpointsInfo(ctx context.Context, clientSet *kubernetes.Clientset, name string, namespace string) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Component", "Name", "Type", "Endpoint"})
	coreV1 := clientSet.CoreV1()

	err := nodePort(ctx, name, coreV1, namespace, t)
	if err != nil {
		return err
	}

	for _, kind := range []string{"metad", "storaged", "graphd"} {
		err = clusterIp(ctx, coreV1, t, name, namespace, kind)
		if err != nil {
			return err
		}
	}
	t.Render()
	return nil
}

func clusterIp(ctx context.Context, coreV1 v1.CoreV1Interface, t table.Writer, name, namespace, kind string) error {
	label := fmt.Sprintf(serviceSelector, name, kind)
	svc, err := coreV1.Services(namespace).List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return err
	}
	if len(svc.Items) == 0 {
		return errors.New("no service found")
	}
	service := svc.Items[0]
	for _, port := range service.Spec.Ports {
		t.AppendRow(table.Row{
			kind,
			port.Name,
			"ClusterIP",
			fmt.Sprintf("%s.%s.svc.cluster.local:%d", service.Name, service.Namespace, port.Port),
		})
	}
	return nil
}

func nodePort(ctx context.Context, name string, coreV1 v1.CoreV1Interface, namespace string, t table.Writer) error {
	const kind = "graphd"
	// query graphd service
	label := fmt.Sprintf(serviceSelector, name, kind)
	svc, err := coreV1.Services(namespace).List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return err
	}
	if len(svc.Items) == 0 {
		return errors.New("no service found")
	}
	service := svc.Items[0]
	if service.Spec.Type != corev1.ServiceTypeNodePort {
		return nil // skip if not node port
	}

	// query ip of all nodes
	list, err := coreV1.Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	var IPs []string
	for _, node := range list.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				IPs = append(IPs, addr.Address)
			}
		}
	}

	// construct node port endpoints
	for _, port := range service.Spec.Ports {
		if port.NodePort < 0 {
			continue
		}
		for _, ip := range IPs {
			t.AppendRow(table.Row{
				kind,
				port.Name,
				"NodePort",
				fmt.Sprintf("%s:%d", ip, port.NodePort),
			})
		}
	}
	return nil
}

func componentInfo(cluster *v1alpha1.NebulaCluster) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"",
		"Phase", "Ready", "Desired",
		"CPU", "Memory", "DataVolume",
		"LogVolume", "Version"})

	status := cluster.Status
	spec := cluster.Spec
	//	Metad
	workload := status.Metad.Workload
	metad := spec.Metad
	t.AppendRow(table.Row{"Metad",
		status.Metad.Phase, workload.ReadyReplicas, *metad.Replicas,
		metad.Resources.Limits.Cpu(), metad.Resources.Limits.Memory(), metad.DataVolumeClaim.Resources.Requests.Storage(),
		metad.LogVolumeClaim.Resources.Requests.Storage(), metad.Version})

	workload = status.Storaged.Workload
	storaged := spec.Storaged

	// compute total storage of storaged
	storagedStorage := computeStoragedVolume(storaged.DataVolumeClaims)
	// Storaged
	t.AppendRow(table.Row{"Storaged",
		status.Storaged.Phase, workload.ReadyReplicas, *storaged.Replicas,
		storaged.Resources.Limits.Cpu(), storaged.Resources.Limits.Memory(), storagedStorage,
		storaged.LogVolumeClaim.Resources.Requests.Storage(), storaged.Version})

	// Graphd
	workload = status.Graphd.Workload
	graphd := spec.Graphd
	t.AppendRow(table.Row{"Graphd",
		status.Graphd.Phase, workload.ReadyReplicas, *graphd.Replicas,
		graphd.Resources.Limits.Cpu(), graphd.Resources.Limits.Memory(), "",
		graphd.LogVolumeClaim.Resources.Requests.Storage(), graphd.Version})
	t.Render()
}

func computeStoragedVolume(volumeClaims []v1alpha1.StorageClaim) string {
	var storagedStorage string
	if len(volumeClaims) > 0 {
		totalStorage := volumeClaims[0].Resources.Requests.Storage()
		for _, claim := range volumeClaims[1:] {
			totalStorage.Add(*claim.Resources.Requests.Storage())
		}
		storagedStorage = totalStorage.String()
	} else {
		storagedStorage = ""
	}
	return storagedStorage
}

func clusterInfo(cluster *v1alpha1.NebulaCluster) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendRow(table.Row{"Name", cluster.Name})
	t.AppendRow(table.Row{"Namespace", cluster.Namespace})
	t.AppendRow(table.Row{"CreationTimestamp", cluster.CreationTimestamp})
	t.Render()
}

func getCluster(ctx context.Context, client *dynamic.DynamicClient, name, namespace string) (*v1alpha1.NebulaCluster, error) {
	resource := v1alpha1.GroupVersion.WithResource("nebulaclusters")
	r, err := client.Resource(resource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	cluster := &v1alpha1.NebulaCluster{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(r.Object, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}
