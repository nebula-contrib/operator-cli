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

package console

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/docker/cli/cli/streams"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const graphdServiceSelector = "app.kubernetes.io/cluster=%s,app.kubernetes.io/component=graphd,app.kubernetes.io/name=nebula-graph"

type Option struct {
	PodName           string
	Name              string // cluster name
	Namespace         string // cluster namespace
	Username          string
	Password          string
	Timeout           int32
	Eval              string
	File              string
	EnableSsl         bool
	SslRootCaPath     string
	SslCertPath       string
	SslPrivateKeyPath string
	GraphdServiceName string
}

const (
	sslRootCaKey     = "ssl.ca"
	sslCertKey       = "ssl.cert"
	sslPrivateKeyKey = "private.key"
	fileKey          = "file.nql"
	mountPathPrefix  = "/etc/nebula"
)

func RunShell(ctx context.Context, clientSet *kubernetes.Clientset, config *rest.Config, option *Option) error {

	svcName, err := getService(context.Background(), clientSet.CoreV1(), option)
	if err != nil {
		return err
	}
	option.GraphdServiceName = svcName

	req := PodExecReq(clientSet, option)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	in := streams.NewIn(os.Stdin)
	if err = in.SetRawTerminal(); err != nil {
		return err
	}
	defer in.RestoreTerminal()
	return exec.StreamWithContext(ctx,
		remotecommand.StreamOptions{
			Stdin:  in,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
}

func PodExecReq(clientSet *kubernetes.Clientset, option *Option) *rest.Request {
	podName, namespace := option.PodName, option.Namespace

	address := fmt.Sprintf("%s.%s.svc.cluster.local", option.GraphdServiceName, namespace)

	cmd := []string{"/usr/local/bin/nebula-console", "-addr", address, "-port", "9669", "-u", option.Username, "-p", option.Password}
	if option.Timeout > 0 {
		cmd = append(cmd, "-t", fmt.Sprintf("%d", option.Timeout))
	}
	if option.Eval != "" {
		cmd = append(cmd, "-e", fmt.Sprintf(`"%s"`, option.Eval))
	}
	if option.File != "" {
		cmd = append(cmd, "-f", path.Join(mountPathPrefix, fileKey))
	}
	if option.EnableSsl {
		cmd = append(cmd, "-enable_ssl")
	}
	if option.SslRootCaPath != "" {
		cmd = append(cmd, "-ssl_root_ca_path", path.Join(mountPathPrefix, sslRootCaKey))
	}
	if option.SslCertPath != "" {
		cmd = append(cmd, "-ssl_cert_path", path.Join(mountPathPrefix, sslCertKey))
	}
	if option.SslPrivateKeyPath != "" {
		cmd = append(cmd, "-ssl_private_key_path", path.Join(mountPathPrefix, sslPrivateKeyKey))
	}

	req := clientSet.CoreV1().RESTClient().
		Post().Resource("pods").
		Name(podName).Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: podName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)
	return req
}

func CratePod(name, namespace string, labels map[string]string, image string, option *Option) *corev1.Pod {
	const volumeName = "mount"
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:            name,
				Image:           image,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Command:         []string{"sh", "-c", "while true; do sleep 2; done"},
				Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("64Mi"),
				}},
			}},
		},
	}

	// add file which is mounted to configmap
	volume := corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: name,
			}}},
	}
	if option.SslRootCaPath != "" {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  sslRootCaKey,
			Path: sslRootCaKey,
		})
	}
	if option.SslCertPath != "" {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  sslCertKey,
			Path: sslCertKey,
		})
	}
	if option.SslPrivateKeyPath != "" {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  sslPrivateKeyKey,
			Path: sslPrivateKeyKey,
		})
	}
	if option.File != "" {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  fileKey,
			Path: fileKey,
		})
	}

	// mount configmap if file exists
	if len(volume.VolumeSource.ConfigMap.Items) > 0 {
		pod.Spec.Volumes = append(pod.Spec.Volumes, volume)
		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: mountPathPrefix,
		})
	}

	return pod
}

func readFileContent(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	all, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(all), nil
}

func CreateConfigMap(name, namespace string, labels map[string]string, option *Option) (*corev1.ConfigMap, error) {
	data := map[string]string{}
	var err error
	data[sslRootCaKey], err = readFileContent(option.SslRootCaPath)
	if err != nil {
		return nil, err
	}

	data[sslCertKey], err = readFileContent(option.SslCertPath)
	if err != nil {
		return nil, err
	}

	data[sslPrivateKeyKey], err = readFileContent(option.SslPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	data[fileKey], err = readFileContent(option.File)
	if err != nil {
		return nil, err
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: data,
	}

	return configMap, nil
}

func getService(ctx context.Context, coreV1 v1.CoreV1Interface, option *Option) (string, error) {
	selector := fmt.Sprintf(graphdServiceSelector, option.Name)
	svc, err := coreV1.Services(option.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return "", err
	}
	if len(svc.Items) == 0 {
		return "", errors.New("no service found")
	}
	return svc.Items[0].Name, nil
}
