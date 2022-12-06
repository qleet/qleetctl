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

var workloadDefinitionConfigPath string

// createWorkloadDefinitionCmd represents the workload-definition command
var createWorkloadDefinitionCmd = &cobra.Command{
	Use:   "workload-definition",
	Short: "Create a new workload definition",
	Long:  `Create a new workload definition.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		configContent, err := ioutil.ReadFile(workloadDefinitionConfigPath)
		if err != nil {
			panic(err)
		}
		var workloadDefinitionConfig config.WorkloadDefinitionConfig
		if err := yaml.Unmarshal(configContent, &workloadDefinitionConfig); err != nil {
			panic(err)
		}

		// create workload definition
		workloadDefinition, err := workloadDefinitionConfig.Create()
		if err != nil {
			panic(err)
		}

		//// get the content of the yaml document
		//definitionContent, err := ioutil.ReadFile(workloadDefinitionConfig.YAMLDocument)
		//if err != nil {
		//	panic(err)
		//}
		//stringContent := string(definitionContent)

		//// construct workload definition object
		//workloadDefinition := &tpapi.WorkloadDefinition{
		//	Name:         &workloadDefinitionConfig.Name,
		//	YAMLDocument: &stringContent,
		//	UserID:       &workloadDefinitionConfig.UserID,
		//}

		//// create workload definition in API
		//wdJSON, err := json.Marshal(&workloadDefinition)
		//if err != nil {
		//	panic(err)
		//}
		//wd, err := tpclient.CreateWorkloadDefinition(wdJSON, "http://localhost:1323", "")
		//if err != nil {
		//	panic(err)
		//}

		fmt.Printf("workload definition %s created\n", *workloadDefinition.Name)
	},
}

func init() {
	createCmd.AddCommand(createWorkloadDefinitionCmd)

	createWorkloadDefinitionCmd.Flags().StringVarP(&workloadDefinitionConfigPath, "config", "c", "", "path to file with workload definition config")
	createWorkloadDefinitionCmd.MarkFlagRequired("config")
}
