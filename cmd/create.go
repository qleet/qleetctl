/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create objects in the Qleet system",
	Long: `Create objects in the Qleet system.

The create command does nothing by itself.  Use one of the avilable subcommands
to create different objects in the system`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
