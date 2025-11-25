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
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

type kincConfig struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Name       string `yaml:"name"`
	Image      string `yaml:"image,omitempty"`
	Networking struct {
		DisableDefaultCNI bool   `yaml:"disableDefaultCNI"`
		KubeProxyMode     string `yaml:"kubeProxyMode"`
		PodSubnet         string `yaml:"podSubnet"`
		ServiceSubnet     string `yaml:"serviceSubnet"`
		IPFamily          string `yaml:"ipFamily"`
	} `yaml:"networking"`
	Nodes []struct {
		Role string `yaml:"role"`
	} `yaml:"nodes"`
}

// createClusterCmd represents the createCluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Creates a local Kubernetes cluster",
	Long:  "Creates a local Kubernetes cluster using Apple container 'nodes'",
	Run: func(cmd *cobra.Command, args []string) {
		var config kincConfig

		image, _ := cmd.Flags().GetString("image")
		if image != "" {
			config.Image = image
		} else {
			config.Image = "kindest/node:v1.34.0"
		}

		configPath, _ := cmd.Flags().GetString("config")
		if configPath != "" {
			var err error
			config, err = loadConfig(configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}
		} else {
			config.Name = "kind"
			config.Kind = "Cluster"
			config.APIVersion = "kind.x-k8s.io/v1alpha4"
			config.Networking.PodSubnet = "10.244.0.0/16"
		}

		config.Name, _ = cmd.Flags().GetString("name")

		err := createCluster(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating cluster: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)

	createClusterCmd.Flags().String("config", "", "path to a kind config file")
	createClusterCmd.Flags().String("image", "", "node docker image to use for booting the cluster")
	createClusterCmd.Flags().String("kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	createClusterCmd.Flags().StringP("name", "n", "kind", "cluster name, overrides KIND_CLUSTER_NAME, config")
	createClusterCmd.Flags().Bool("retain", false, "retain nodes for debugging when cluster creation fails")
	createClusterCmd.Flags().Int("wait", 0, "wait for control plane node to be ready  (default 0s)")

}

func createCluster(config kincConfig) error {
	fmt.Printf("Creating cluster '%s' ...\n", config.Name)

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = " "
	// s.Suffix = fmt.Sprintf(" Ensuring node image (%s)", config.Image)
	// s.FinalMSG = fmt.Sprintf(" \033[32m✓\033[0m Ensuring node image (%s)\n", config.Image)

	// s.Start()

	// Pull the node image
	myCmd := exec.Command("container", "image", "pull", config.Image)
	if err := runCommand(myCmd, true); err != nil {
		cobra.CheckErr(fmt.Errorf("\nfailed to pull node image: %w", err))
	}
	fmt.Printf(" \033[32m✓\033[0m Ensuring node image (%s) 🖼\n", config.Image)

	// Prepare the node
	s.Suffix = " Preparing nodes 📦"
	s.FinalMSG = " \033[32m✓\033[0m Preparing nodes 📦\n"
	s.Start()
	myCmd = exec.Command("container", "run",
		"-d",
		"--name", config.Name+"-control-plane",
		"-m", "8G", "--disable-progress-updates",
		"-e", "KUBECONFIG=/etc/kubernetes/admin.conf",
		"-l", "io.x-k8s.kinc.cluster=kinc",
		"-p", "127.0.0.1:6443:6443",
		config.Image,
	)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("\nfailed to run node: %w", err)
	}
	s.Stop()

	// Writing config
	s.Suffix = " Writing configuration 📜"
	s.FinalMSG = " \033[32m✓\033[0m Writing configuration 📜\n"
	s.Start()
	myCmd = exec.Command("container", "exec", config.Name+"-control-plane", "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("\nfailed to write config: %w", err)
	}
	s.Stop()

	// Starting Control Plane
	s.Suffix = " Starting control-plane 🕹️"
	s.FinalMSG = " \033[32m✓\033[0m Starting control-plane 🕹️\n"
	s.Start()
	myCmd = exec.Command("container", "exec", config.Name+"-control-plane", "kubeadm", "init", "--pod-network-cidr="+config.Networking.PodSubnet)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("\nfailed to init cluster: %w", err)
	}

	myCmd = exec.Command("container", "exec", config.Name+"-control-plane", "kubectl", "taint", "nodes", "--all", "node-role.kubernetes.io/control-plane-")
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to remove taint: %w", err)
	}

	s.Stop()

	// Installing CNI
	s.Suffix = " Installing CNI 🔌"
	s.FinalMSG = " \033[32m✓\033[0m Installing CNI 🔌\n"
	s.Start()
	cniCmd := fmt.Sprintf("sed -e 's@{{ .PodSubnet }}@%s@' /kind/manifests/default-cni.yaml | kubectl apply -f -", config.Networking.PodSubnet)
	myCmd = exec.Command("container", "exec", config.Name+"-control-plane", "sh", "-euc", cniCmd)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("\nfailed to set up CNI: %w", err)
	}
	s.Stop()

	// Installing StorageClass 💾
	s.Suffix = " Installing StorageClass 💾"
	s.FinalMSG = " \033[32m✓\033[0m Installing StorageClass 💾\n"
	s.Start()
	storageCmd := "cat /kind/manifests/default-storage.yaml | kubectl apply -f -"
	myCmd = exec.Command("container", "exec", config.Name+"-control-plane", "sh", "-euc", storageCmd)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to set up StorageClass: %w", err)
	}
	s.Stop()

	// Set up kubeconfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	kubeDir := filepath.Join(homeDir, ".kube")
	if err := os.MkdirAll(kubeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .kube directory: %w", err)
	}

	// Get kubeconfig from container
	myCmd = exec.Command("container", "exec", config.Name+"-control-plane", "cat", "/etc/kubernetes/admin.conf")
	output, err := myCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	kubeconfigPath := filepath.Join(kubeDir, "config")
	if err := os.WriteFile(kubeconfigPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return nil
}

func loadConfig(configPath string) (kincConfig, error) {
	var config kincConfig

	// Read and parse the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}
