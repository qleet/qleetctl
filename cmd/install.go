/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/qleet/qleetctl/internal/install"
	"github.com/qleet/qleetctl/internal/provider"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the QleetOS controle plane",
	Long:  `Install the QleetOS controle plane.`,
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
		apiDepsManifest.WriteString(install.APIDepsManifest)
		fmt.Println("QleetOS API dependencies written to /tmp directory")

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
			"threeport-api",
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
		apiServerManifest.WriteString(install.APIServerManifest)
		fmt.Println("QleetOS API server written to /tmp directory")

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
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	installCmd.Flags().StringP("provider", "p", "kind", "The infrasture tool or provider to install upon")
}
