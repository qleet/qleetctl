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
	Short: "Update QleetOS objects",
	Long: `Update QleetOS objects.

The update command does nothing by itself.  Use one of the avilable subcommands
to update different objects in the system`,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
