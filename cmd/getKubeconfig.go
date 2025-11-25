/*
Copyright © 2025 Christophe Jauffret <reg-github@geo6.net>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// getKubeconfigCmd represents the getKubeconfig command
var getKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Prints cluster kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {

		name, _ := cmd.Flags().GetString("name")

		// Get kubeconfig from container
		myCmd := exec.Command("container", "exec", name+"-control-plane", "cat", "/etc/kubernetes/admin.conf")
		if err := runCommand(myCmd, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieve kubeconfig: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getKubeconfigCmd)

	getKubeconfigCmd.Flags().StringP("name", "n", "kind", "the cluster context name")
}
