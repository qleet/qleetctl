/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	qout "github.com/qleet/qleetctl/internal/output"
	"github.com/qleet/qleetctl/internal/provider"
	"github.com/spf13/cobra"
)

// deleteQleetosCmd represents the delete qleetos command
var deleteQleetosCmd = &cobra.Command{
	Use:          "qleetos",
	Example:      "qleetctl delete qleetos",
	Short:        "A brief description of your command",
	Long:         `A long description of your command.`,
	SilenceUsage: true,
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
			qout.Error("failed to delete kind cluster", err)
			os.Exit(1)
		}
		qout.Info("kind cluster deleted")

		qout.Complete("QleetOS instance deleted")
	},
}

func init() {
	deleteCmd.AddCommand(deleteQleetosCmd)
}
