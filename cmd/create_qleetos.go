/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tpclient "github.com/threeport/threeport-go-client"
	tpapi "github.com/threeport/threeport-rest-api/pkg/api/v0"
	kubeclient "k8s.io/client-go/tools/clientcmd"

	"github.com/qleet/qleetctl/internal/config"
	"github.com/qleet/qleetctl/internal/install"
	qout "github.com/qleet/qleetctl/internal/output"
	"github.com/qleet/qleetctl/internal/provider"
)

var (
	createQleetOSInstanceName string
	forceOverwriteConfig      bool
)

// createQleetosCmd represents the create qleetos command
var createQleetosCmd = &cobra.Command{
	Use:          "qleetos",
	Example:      "qleetctl create qleetos",
	Short:        "Create a new instance of the QleetOS control plane",
	Long:         `Create a new instance of the QleetOS control plane.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// get qleet config
		qleetConfig := &config.QleetConfig{}
		if err := viper.Unmarshal(qleetConfig); err != nil {
			qout.Error("failed to get Qleet config", err)
		}

		// check qleet config for exisiting instance
		qleetOSInstanceConfigExists := false
		for _, instance := range qleetConfig.QleetOSInstances {
			if instance.Name == createQleetOSInstanceName {
				qleetOSInstanceConfigExists = true
				if !forceOverwriteConfig {
					qout.Error(
						"interupted creation of QleetOS instance",
						errors.New(fmt.Sprintf("instance of QleetOS with name %s already exists", instance.Name)),
					)
					qout.Info("if you wish to overwrite the existing config use --force-overwrite-config flag")
					qout.Warning("you will lose the ability to connect to the existing QleetOS instance if it still exists")
					os.Exit(1)
				}
			}
		}

		// write kind config file to /tmp directory
		configFile, err := os.Create(provider.QleetKindConfigPath)
		if err != nil {
			qout.Error("failed to write kind config file to disk", err)
			os.Exit(1)
		}
		defer configFile.Close()
		configFile.WriteString(provider.KindConfig(createQleetOSInstanceName))
		qout.Info("kind config written to /tmp directory")

		// start kind cluster
		qout.Info("creating kind cluster... (this could take a few minutes)")
		kindCreate := exec.Command(
			"kind",
			"create",
			"cluster",
			"--config",
			provider.QleetKindConfigPath,
		)
		if err := kindCreate.Run(); err != nil {
			qout.Error("failed to create new kind cluster", err)
			os.Exit(1)
		}
		qout.Info("kind cluster created")

		// write API dependencies manifest to /tmp directory
		apiDepsManifest, err := os.Create(install.APIDepsManifestPath)
		if err != nil {
			qout.Error("failed to write API dependency manifests to disk", err)
			os.Exit(1)
		}
		defer apiDepsManifest.Close()
		apiDepsManifest.WriteString(install.APIDepsManifest())
		qout.Info("QleetOS API dependencies manifest written to /tmp directory")

		// install API dependencies on kind cluster
		qout.Info("installing QleetOS API dependencies")
		apiDepsCreate := exec.Command(
			"kubectl",
			"apply",
			"-f",
			install.APIDepsManifestPath,
		)
		if err := apiDepsCreate.Run(); err != nil {
			qout.Error("failed to install API dependencies to kind cluster", err)
			os.Exit(1)
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
			qout.Error("failed to create API database config", err)
			os.Exit(1)
		}

		qout.Info("QleetOS API dependencies created")

		// write API server manifest to /tmp directory
		apiServerManifest, err := os.Create(install.APIServerManifestPath)
		if err != nil {
			qout.Error("failed to write API manifest to disk", err)
			os.Exit(1)
		}
		defer apiServerManifest.Close()
		apiServerManifest.WriteString(install.APIServerManifest())
		qout.Info("QleetOS API server manifest written to /tmp directory")

		// install QleetOS API
		qout.Info("installing QleetOS API server")
		apiServerCreate := exec.Command(
			"kubectl",
			"apply",
			"-f",
			install.APIServerManifestPath,
		)
		if err := apiServerCreate.Run(); err != nil {
			qout.Error("failed to create API server", err)
			os.Exit(1)
		}

		qout.Info("QleetOS API server created")

		// write workload controller manifest to /tmp directory
		workloadControllerManifest, err := os.Create(install.WorkloadControllerManifestPath)
		if err != nil {
			qout.Error("failed to write workload controller manifest to disk", err)
			os.Exit(1)
		}
		defer workloadControllerManifest.Close()
		workloadControllerManifest.WriteString(install.WorkloadControllerManifest())
		qout.Info("QleetOS workload controller manifest written to /tmp directory")

		// install workload controller
		qout.Info("installing QleetOS workload controller")
		workloadControllerCreate := exec.Command(
			"kubectl",
			"apply",
			"-f",
			install.WorkloadControllerManifestPath,
		)
		if err := workloadControllerCreate.Run(); err != nil {
			qout.Error("failed to create workload controller", err)
			os.Exit(1)
		}

		qout.Info("QleetOS workload controller created")

		// wait a few seconds for everything to come up
		qout.Info("waiting for components to spin up...")
		time.Sleep(time.Second * 200)

		// get kubeconfig
		defaultLoadRules := kubeclient.NewDefaultClientConfigLoadingRules()

		clientConfigLoadRules, err := defaultLoadRules.Load()
		if err != nil {
			qout.Error("failed to load default kubeconfig rules", err)
			os.Exit(1)
		}

		clientConfig := kubeclient.NewDefaultClientConfig(*clientConfigLoadRules, &kubeclient.ConfigOverrides{})
		kubeConfig, err := clientConfig.RawConfig()
		if err != nil {
			qout.Error("failed to load kubeconfig", err)
			os.Exit(1)
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
			qout.Error(
				"failed to get Kubernetes cluster CA and endpoint",
				errors.New("cluster config not found in kubeconfig"),
			)
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
			qout.Error(
				"failed to get user credentials to Kubernetes cluster",
				errors.New("kubeconfig user for qleet-os cluster not found"),
			)
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
			qout.Error("failed to marshal workload cluster to json", err)
			os.Exit(1)
		}
		wc, err := tpclient.CreateWorkloadCluster(wcJSON, install.GetQleetOSAPIEndpoint(), "")
		if err != nil {
			qout.Error("failed to create workload cluster in Qleet API", err)
			os.Exit(1)
		}
		qout.Info(fmt.Sprintf("default workload cluster %s for compute space set up", *wc.Name))

		// create qleet config for new instance
		newQleetOSInstance := &config.QleetOSInstance{
			Name:      createQleetOSInstanceName,
			APIServer: install.GetQleetOSAPIEndpoint(),
		}

		// update qleet config to add the new instance and current instance
		if qleetOSInstanceConfigExists {
			for n, instance := range qleetConfig.QleetOSInstances {
				if instance.Name == createQleetOSInstanceName {
					qleetConfig.QleetOSInstances[n] = *newQleetOSInstance
				}
			}
		} else {
			qleetConfig.QleetOSInstances = append(qleetConfig.QleetOSInstances, *newQleetOSInstance)
		}
		viper.Set("QleetOSInstances", qleetConfig.QleetOSInstances)
		viper.Set("CurrentInstance", createQleetOSInstanceName)
		viper.WriteConfig()
		qout.Info("Qleet config updated")

		qout.Complete("QleetOS instance created")
	},
}

func init() {
	createCmd.AddCommand(createQleetosCmd)

	createQleetosCmd.Flags().StringP("provider", "p", "kind", "the infrasture provider to install upon")
	createQleetosCmd.Flags().StringVarP(&createQleetOSInstanceName, "name", "n", "", "name of Qleet OS instance")
	createQleetosCmd.MarkFlagRequired("name")
	createQleetosCmd.Flags().BoolVar(&forceOverwriteConfig, "force-overwrite-config", false, "force the overwrite of an existing QleetOS instance with the same name in Qleet config")
}
