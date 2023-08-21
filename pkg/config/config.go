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

package config

import (
	"os"
	"path"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/util/homedir"
)

const (
	NgctlConfigPath = ".ngctl/config"
)

type NgctlConfig struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func SaveConfig(namespace string, name string) error {
	configPath := path.Join(homedir.HomeDir(), NgctlConfigPath)
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// rwxr-x---
		if err = os.MkdirAll(dir, 0750); err != nil {
			return err
		}
	}
	ncConfig := NgctlConfig{
		Namespace: namespace,
		Name:      name,
	}
	content, err := json.Marshal(ncConfig)
	if err != nil {
		return err
	}
	// rw-r-----
	if err = os.WriteFile(configPath, content, 0640); err != nil {
		return err
	}
	return nil
}

func LoadConfig() (NgctlConfig, error) {
	var conf NgctlConfig
	configPath := path.Join(homedir.HomeDir(), NgctlConfigPath)
	content, err := os.ReadFile(configPath)
	if err != nil {
		return conf, err
	}
	if err = json.Unmarshal(content, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
