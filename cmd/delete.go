/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete QleetOS objects",
	Long: `Delete QleetOS objects.

The delete command does nothing by itself.  Use one of the avilable subcommands
to delete different objects in the system`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
