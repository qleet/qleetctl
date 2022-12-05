/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/qleet/qleetctl/internal/config"
)

var workloadInstanceConfigPath string

// createWorkloadInstanceCmd represents the workload-instance command
var createWorkloadInstanceCmd = &cobra.Command{
	Use:   "workload-instance",
	Short: "Create a new workload instance",
	Long:  `Create a new workload instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		configContent, err := ioutil.ReadFile(workloadInstanceConfigPath)
		if err != nil {
			panic(err)
		}
		var workloadInstanceConfig config.WorkloadInstanceConfig
		if err := yaml.Unmarshal(configContent, &workloadInstanceConfig); err != nil {
			panic(err)
		}

		// create workload instance
		workloadInstance, err := workloadInstanceConfig.Create()
		if err != nil {
			panic(err)
		}

		//// get workload cluster by name
		//workloadCluster, err := tpclient.GetWorkloadClusterByName(
		//	workloadInstanceConfig.WorkloadClusterName,
		//	"http://localhost:1323", "",
		//)
		//if err != nil {
		//	panic(err)
		//}

		//// get workload definition by name
		//workloadDefinition, err := tpclient.GetWorkloadDefinitionByName(
		//	workloadInstanceConfig.WorkloadDefinitionName,
		//	"http://localhost:1323", "",
		//)
		//if err != nil {
		//	panic(err)
		//}

		//// construct workload instance object
		//workloadInstance := &tpapi.WorkloadInstance{
		//	Name:                 &workloadInstanceConfig.Name,
		//	WorkloadClusterID:    &workloadCluster.ID,
		//	WorkloadDefinitionID: &workloadDefinition.ID,
		//}

		//// create workload instance in API
		//wiJSON, err := json.Marshal(&workloadInstance)
		//if err != nil {
		//	panic(err)
		//}
		//wi, err := tpclient.CreateWorkloadInstance(wiJSON, "http://localhost:1323", "")
		//if err != nil {
		//	panic(err)
		//}

		fmt.Printf("workload instance %s created\n", *workloadInstance.Name)
	},
}

func init() {
	createCmd.AddCommand(createWorkloadInstanceCmd)

	createWorkloadInstanceCmd.Flags().StringVarP(&workloadInstanceConfigPath, "config", "c", "", "path to file with workload instance config")
	createWorkloadInstanceCmd.MarkFlagRequired("config")
}
