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
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

// Flag Values for RootCmd
var (
	kubeConfig string
)

var RootCmd = &cobra.Command{
	Use:   "ngctl [command]",
	Short: "the command line tool for nebula operator",
	Long:  "ngctl is the command line tool for nebula operator.",
}

func init() {
	RootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "path of the kubernetes config file")
	RootCmd.AddCommand(studioCmd())
}
