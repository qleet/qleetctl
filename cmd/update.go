/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update objects in the Qleet system",
	Long: `Update objects in the Qleet system.

The update command does nothing by itself.  Use one of the avilable subcommands
to update different objects in the system`,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
