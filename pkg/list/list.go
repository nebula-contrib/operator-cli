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

package list

import (
	"context"

	"github.com/vesoft-inc/nebula-operator/apis/apps/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
)

func Clusters(ctx context.Context, client *dynamic.DynamicClient, allNamespaces bool, namespace string) ([]v1alpha1.NebulaCluster, error) {
	resource := v1alpha1.GroupVersion.WithResource("nebulaclusters")
	resourceInterface := client.Resource(resource)
	var clusterList *unstructured.UnstructuredList
	var err error
	if allNamespaces {
		clusterList, err = resourceInterface.List(ctx, v1.ListOptions{})
	} else {
		clusterList, err = resourceInterface.Namespace(namespace).List(ctx, v1.ListOptions{})
	}
	if err != nil {
		return nil, err
	}

	var clusters []v1alpha1.NebulaCluster
	for _, item := range clusterList.Items {
		cluster := v1alpha1.NebulaCluster{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &cluster)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, cluster)
	}
	return clusters, nil
}
