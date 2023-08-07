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
	"context"
	"errors"
	"log"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resource is the generic type for kubernetes resource
type Resource interface {
	appsv1.Deployment | corev1.Service | corev1.Pod
}

// ResourceInterface is the generic interface for kubernetes resource management
type ResourceInterface[T Resource] interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*T, error)
	Create(ctx context.Context, resource *T, opts metav1.CreateOptions) (*T, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

type Status int

const (
	StatusAlready Status = iota
	StatusNotFound
	StatusConflicted
)

var (
	ErrConflict = errors.New("resource conflicted")
)

// Check checks target resource and returns status
func Check[T Resource](ctx context.Context, si ResourceInterface[T], name, label string) Status {
	resource, err := si.Get(ctx, name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		return StatusNotFound
	}
	/*
		resource.Labels / resource.GetLabels() are not available here
		the following code will cause error:
		if resource.GetLabels()[label] != name {
		.....
		}
	*/

	// get labels from resource by reflection
	labels, ok := reflect.ValueOf(resource).Elem().FieldByName("ObjectMeta").FieldByName("Labels").Interface().(map[string]string)
	if !ok || labels[label] != name {
		return StatusConflicted
	}
	return StatusAlready
}

// RemoveResource removes a resource
func RemoveResource[T Resource](ctx context.Context, si ResourceInterface[T], name, label string) error {
	status := Check(ctx, si, name, label)
	typ := reflect.TypeOf(new(T)).Elem().Name()
	if status == StatusNotFound {
		return nil
	} else if status == StatusConflicted {
		log.Printf("Resource %s %s is conflicted\n", typ, name)
		return nil
	}
	err := si.Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("Resource %s %s is removed\n", typ, name)
	return nil
}

// CreateResource creates a resource
func CreateResource[T Resource](ctx context.Context, si ResourceInterface[T], elem *T, name, label string) error {
	status := Check(ctx, si, name, label)
	typ := reflect.TypeOf(elem).Elem().Name()
	if status == StatusAlready {
		log.Printf("Resource %s %s is already created\n", typ, name)
		return nil
	}
	if status == StatusConflicted {
		log.Printf("Resource %s %s is conflicted\n", typ, name)
		return ErrConflict
	}
	_, err := si.Create(ctx, elem, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Printf("Resource %s %s is created\n", typ, name)
	return nil
}
