/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/qleet/qleetctl/internal/provider"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the QleetOS controle plane",
	Long:  `Uninstall the QleetOS controle plane.`,
	Run: func(cmd *cobra.Command, args []string) {
		// delete kind cluster
		fmt.Println("deleting kind cluster...")
		kindDelete := exec.Command(
			"kind",
			"delete",
			"cluster",
			"--name",
			provider.QleetKindClusterName,
		)
		if err := kindDelete.Run(); err != nil {
			panic(err)
		}
		fmt.Println("kind cluster deleted")
		fmt.Println("QleetOS instance uninstalled")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
