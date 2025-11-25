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
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Deletes a cluster",
	Long: `Deletes a Kind cluster from the system.

This is an idempotent operation, meaning it may be called multiple times without
failing (like "rm -f"). If the cluster resources exist they will be deleted, and
if the cluster is already gone it will just return success.

Errors will only occur if the cluster resources exist and are not able to be deleted.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := cmd.Flags().Lookup("name").Value.String()

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = " "
		s.Suffix = fmt.Sprintf(" Deleting cluster '%s' ...", name)
		s.FinalMSG = fmt.Sprintf(" \033[32m✓\033[0m Cluster '%s' deleted.\n", name)
		s.Start()

		myCmd := exec.Command("container", "rm", "-f", name+"-control-plane")
		_ = myCmd.Run()
		s.Stop()

	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	deleteClusterCmd.Flags().StringP("name", "n", "kind", "the cluster name")
}
