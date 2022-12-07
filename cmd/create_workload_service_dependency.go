/*
Copyright © 2023 Qleet admin@qleet.io
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/qleet/qleetctl/internal/api"
	qout "github.com/qleet/qleetctl/internal/output"
)

var createWorkloadServiceDependencyConfigPath string

// createWorkloadServiceDependencyCmd represents the workload-service-dependency command
var createWorkloadServiceDependencyCmd = &cobra.Command{
	Use:          "workload-service-dependency",
	Example:      "qleetctl create workload-servicde-dependency -c /path/to/config.yaml",
	Short:        "Create a new workload service dependency",
	Long:         `Create a new workload service dependency.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		configContent, err := ioutil.ReadFile(createWorkloadServiceDependencyConfigPath)
		if err != nil {
			qout.Error("failed to read config file", err)
			os.Exit(1)
		}
		var workloadServiceDependency api.WorkloadServiceDependencyConfig
		if err := yaml.Unmarshal(configContent, &workloadServiceDependency); err != nil {
			qout.Error("failed to unmarshal config file yaml content", err)
			os.Exit(1)
		}

		// create workload service dependency
		wsd, err := workloadServiceDependency.Create()
		if err != nil {
			qout.Error("failed to create workload", err)
			os.Exit(1)
		}

		qout.Complete(fmt.Sprintf("workload service dependency %s created\n", *wsd.Name))
	},
}

func init() {
	createCmd.AddCommand(createWorkloadServiceDependencyCmd)

	createWorkloadServiceDependencyCmd.Flags().StringVarP(&createWorkloadServiceDependencyConfigPath, "config", "c", "", "path to file with workload service dependency config")
	createWorkloadServiceDependencyCmd.MarkFlagRequired("config")
}
