/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of qleetctl",
	Long:  `Print the version of qleetctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
