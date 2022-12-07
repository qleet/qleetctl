/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	tpclient "github.com/threeport/threeport-go-client"
	tpapi "github.com/threeport/threeport-rest-api/pkg/api/v0"
	kubeclient "k8s.io/client-go/tools/clientcmd"

	"github.com/qleet/qleetctl/internal/install"
	"github.com/qleet/qleetctl/internal/provider"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the QleetOS control plane",
	Long:  `Install the QleetOS control plane.`,
	Run: func(cmd *cobra.Command, args []string) {
		// write kind config file to /tmp directory
		configFile, err := os.Create(provider.QleetKindConfigPath)
		if err != nil {
			panic(err)
		}
		defer configFile.Close()
		configFile.WriteString(provider.KindConfig())
		fmt.Println("kind config written to /tmp directory")

		// start kind cluster
		fmt.Println("creating kind cluster... (this could take a few minutes)")
		kindCreate := exec.Command(
			"kind",
			"create",
			"cluster",
			"--config",
			provider.QleetKindConfigPath,
		)
		if err := kindCreate.Run(); err != nil {
			panic(err)
		}
		fmt.Println("kind cluster created")

		// write API dependencies manifest to /tmp directory
		apiDepsManifest, err := os.Create(install.APIDepsManifestPath)
		if err != nil {
			panic(err)
		}
		defer apiDepsManifest.Close()
		apiDepsManifest.WriteString(install.APIDepsManifest())
		fmt.Println("QleetOS API dependencies manifest written to /tmp directory")

		// install API dependencies on kind cluster
		fmt.Println("installing QleetOS API dependencies")
		apiDepsCreate := exec.Command(
			"kubectl",
			"apply",
			"-f",
			install.APIDepsManifestPath,
		)
		if err := apiDepsCreate.Run(); err != nil {
			panic(err)
		}
		psqlConfigCreate := exec.Command(
			"kubectl",
			"create",
			"configmap",
			"postgres-config-data",
			"-n",
			install.ThreeportControlPlaneNs,
		)
		if err := psqlConfigCreate.Run(); err != nil {
			panic(err)
		}

		fmt.Println("QleetOS API dependencies created")

		// write API server manifest to /tmp directory
		apiServerManifest, err := os.Create(install.APIServerManifestPath)
		if err != nil {
			panic(err)
		}
		defer apiServerManifest.Close()
		apiServerManifest.WriteString(install.APIServerManifest())
		fmt.Println("QleetOS API server manifest written to /tmp directory")

		// install QleetOS API
		fmt.Println("installing QleetOS API server")
		apiServerCreate := exec.Command(
			"kubectl",
			"apply",
			"-f",
			install.APIServerManifestPath,
		)
		if err := apiServerCreate.Run(); err != nil {
			panic(err)
		}

		fmt.Println("QleetOS API server created")

		// write workload controller manifest to /tmp directory
		workloadControllerManifest, err := os.Create(install.WorkloadControllerManifestPath)
		if err != nil {
			panic(err)
		}
		defer workloadControllerManifest.Close()
		workloadControllerManifest.WriteString(install.WorkloadControllerManifest())
		fmt.Println("QleetOS workload controller manifest written to /tmp directory")

		// install workload controller
		fmt.Println("installing QleetOS workload controller")
		workloadControllerCreate := exec.Command(
			"kubectl",
			"apply",
			"-f",
			install.WorkloadControllerManifestPath,
		)
		if err := workloadControllerCreate.Run(); err != nil {
			panic(err)
		}

		fmt.Println("QleetOS workload controller created")

		// wait a few seconds for everything to come up
		fmt.Println("waiting for components to spin up...")
		time.Sleep(time.Second * 200)

		// get kubeconfig
		defaultLoadRules := kubeclient.NewDefaultClientConfigLoadingRules()

		clientConfigLoadRules, err := defaultLoadRules.Load()
		if err != nil {
			panic(err)
		}

		clientConfig := kubeclient.NewDefaultClientConfig(*clientConfigLoadRules, &kubeclient.ConfigOverrides{})
		kubeConfig, err := clientConfig.RawConfig()
		if err != nil {
			panic(err)
		}

		// get cluster CA and server endpoint
		var caCert string
		clusterFound := false
		for clusterName, cluster := range kubeConfig.Clusters {
			if clusterName == kubeConfig.CurrentContext {
				caCert = string(cluster.CertificateAuthorityData)
				clusterFound = true
			}
		}
		if !clusterFound {
			fmt.Println("kubeconfig cluster for qleet-os not found")
			os.Exit(1)
		}

		// get client certificate and key
		var cert string
		var key string
		userFound := false
		for userName, user := range kubeConfig.AuthInfos {
			if userName == kubeConfig.CurrentContext {
				cert = string(user.ClientCertificateData)
				key = string(user.ClientKeyData)
				userFound = true
			}
		}
		if !userFound {
			fmt.Println("kubeconfig user for qleet-os not found")
			os.Exit(1)
		}

		// setup default compute space cluster
		clusterName := "default-qleet-compute-space"
		clusterRegion := "local"
		clusterProvider := "kind"
		server := "kubernetes.default"
		workloadCluster := tpapi.WorkloadCluster{
			Name:          &clusterName,
			Region:        &clusterRegion,
			Provider:      &clusterProvider,
			APIEndpoint:   &server,
			CACertificate: &caCert,
			Certificate:   &cert,
			Key:           &key,
		}
		wcJSON, err := json.Marshal(&workloadCluster)
		if err != nil {
			panic(err)
		}
		wc, err := tpclient.CreateWorkloadCluster(wcJSON, "http://localhost:1323", "")
		if err != nil {
			panic(err)
		}

		fmt.Printf("default workload cluster %s for compute space set up\n", *wc.Name)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("provider", "p", "kind", "The infrasture tool or provider to install upon")
}
