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

	"github.com/spf13/cobra"
)

// deleteClustersCmd represents the deleteClusters command
var deleteClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Deletes one or more clusters",
	Long: `Deletes one or more Kind clusters from the system.

This is an idempotent operation, meaning it may be called multiple times without
failing (like "rm -f"). If the cluster resources exist they will be deleted, and
if the cluster is already gone it will just return success.

Errors will only occur if the cluster resources exist and are not able to be deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteClusters called")
	},
}

func init() {
	deleteCmd.AddCommand(deleteClustersCmd)

	// Here you will define your flags and configuration settings.

	deleteClustersCmd.Flags().BoolP("all", "A", false, "delete all clusters")
}
