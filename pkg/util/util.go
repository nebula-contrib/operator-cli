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

package util

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newKubeConfig(kubeconfig string, context string) (*rest.Config, error) {
	if context == "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	// switch to the specified context
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

// NewClientSet creates a new kubernetes client set
func NewClientSet(kubeconfig string, context string) (*kubernetes.Clientset, *rest.Config, error) {
	config, err := newKubeConfig(kubeconfig, context)
	if err != nil {
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientSet, config, nil
}

// NewDynamicClient creates a new dynamic client
func NewDynamicClient(kubeConfig string, context string) (*dynamic.DynamicClient, *rest.Config, error) {
	config, err := newKubeConfig(kubeConfig, context)
	if err != nil {
		return nil, nil, err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return client, config, nil
}
