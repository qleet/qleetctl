/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qleet/qleetctl/internal/config"
	qout "github.com/qleet/qleetctl/internal/output"
	"github.com/qleet/qleetctl/internal/provider"
)

var deleteQleetOSInstanceName string

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
			provider.GetQleetKindClusterName(deleteQleetOSInstanceName),
		)
		if err := kindDelete.Run(); err != nil {
			qout.Error("failed to delete kind cluster", err)
			os.Exit(1)
		}
		qout.Info("kind cluster deleted")

		// get qleet config
		qleetConfig := &config.QleetConfig{}
		if err := viper.Unmarshal(qleetConfig); err != nil {
			qout.Error("failed to get Qleet config", err)
		}

		// update qleet config to remove the deleted Qleet OS instance and
		// current instance
		updatedQleetOSInstances := []config.QleetOSInstance{}
		for _, instance := range qleetConfig.QleetOSInstances {
			if instance.Name == deleteQleetOSInstanceName {
				continue
			} else {
				updatedQleetOSInstances = append(updatedQleetOSInstances, instance)
			}
		}

		viper.Set("QleetOSInstances", updatedQleetOSInstances)
		viper.Set("CurrentInstance", "")
		viper.WriteConfig()
		qout.Info("Qleet config updated")

		qout.Complete("QleetOS instance deleted")
	},
}

func init() {
	deleteCmd.AddCommand(deleteQleetosCmd)

	deleteQleetosCmd.Flags().StringVarP(&deleteQleetOSInstanceName, "name", "n", "", "name of Qleet OS instance")
	deleteQleetosCmd.MarkFlagRequired("name")
}
