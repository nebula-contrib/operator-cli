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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/nebula-contrib/ngctl/pkg/config"
	"github.com/nebula-contrib/ngctl/pkg/console"
	"github.com/nebula-contrib/ngctl/pkg/util"
)

const consoleLabel = "ngctl/nebula-console"

func consoleCmd() *cobra.Command {
	var (
		image string
	)
	option := console.Option{}
	cmd := &cobra.Command{
		Use:   "console",
		Short: "nebula console client for nebula graph ",
		Long:  "nebula console client for nebula graph.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(option, image)
		},
	}

	cmd.PersistentFlags().StringVar(&image, "image", "vesoft/nebula-console:v3.5", "image of the nebula graph console")
	cmd.PersistentFlags().StringVarP(&option.PodName, "pod_name", "n", "nebula-console", "set the name of the console pod. ")

	cmd.PersistentFlags().StringVarP(&option.Username, "user", "u", "root", "set the username of the NebulaGraph account. ")
	cmd.PersistentFlags().StringVarP(&option.Password, "password", "p", "", "set the password of the NebulaGraph account. ")
	cmd.PersistentFlags().Int32VarP(&option.Timeout, "timeout", "t", 120, "set the connection timeout in milliseconds. ")
	cmd.PersistentFlags().StringVarP(&option.Eval, "eval", "e", "", "set the nGQL statement in string type. ")
	cmd.PersistentFlags().StringVarP(&option.File, "file", "f", "", "set the path of the file that stores nGQL statements. ")
	cmd.PersistentFlags().BoolVarP(&option.EnableSsl, "enable_ssl", "", false, "connect to NebulaGraph using SSL encryption and two-way authentication. ")
	cmd.PersistentFlags().StringVarP(&option.SslRootCaPath, "ssl_root_ca_path", "", "", "specify the path of the CA root certificate. ")
	cmd.PersistentFlags().StringVarP(&option.SslCertPath, "ssl_cert_path", "", "", "specify the path of the SSL public key certificate. ")
	cmd.PersistentFlags().StringVarP(&option.SslPrivateKeyPath, "ssl_private_key_path", "", "", "specify the path of the SSL key. ")

	return cmd
}

func run(option console.Option, image string) error {
	ngctlConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}
	option.Name, option.Namespace = ngctlConfig.Name, ngctlConfig.Namespace

	ctx := context.Background()

	clientSet, conf, err := util.NewClientSet(kubeConfig, kubeContext)
	if err != nil {
		return nil
	}

	err = initConsole(ctx, clientSet, &option, image)
	if err != nil {
		return err
	}

	err = console.RunShell(ctx, clientSet, conf, &option)
	if err != nil {
		return err
	}

	return nil
}

func initConsole(ctx context.Context, clientSet *kubernetes.Clientset, option *console.Option, image string) error {
	podName, namespace := option.PodName, option.Namespace
	if err := loadConfig(ctx, clientSet, option, podName, namespace); err != nil {
		return err
	}
	pods := clientSet.CoreV1().Pods(namespace)
	status := util.Check[corev1.Pod](ctx, pods, podName, consoleLabel)
	if status == util.StatusConflicted {
		return errors.New("console pod is conflicted with already exist pod, please check")
	} else if status == util.StatusAlready {
		log.Printf("console pod is already ready, skip init pod")
	} else {
		// create pod
		labels := map[string]string{
			consoleLabel: podName,
		}
		pod := console.CratePod(podName, namespace, labels, image, option)
		_, err := clientSet.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// watch pod status until it is running
	watch, err := clientSet.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + podName,
	})
	if err != nil {
		return err
	}
	for event := range watch.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			log.Printf("convert event object to pod failed")
			continue
		}
		if pod.Status.Phase == corev1.PodRunning {
			log.Printf("console pod is ready")
			break
		}
	}
	return nil
}

func loadConfig(ctx context.Context, clientSet *kubernetes.Clientset, option *console.Option, name, namespace string) error {
	configMaps := clientSet.CoreV1().ConfigMaps(namespace)
	labels := map[string]string{
		consoleLabel: name,
	}
	configMap, err := console.CreateConfigMap(name, namespace, labels, option)
	if err != nil {
		return err
	}
	err = util.CreateOrUpdate[corev1.ConfigMap](ctx, configMaps, configMap, name, consoleLabel)
	if err != nil {
		return err
	}
	return nil
}
