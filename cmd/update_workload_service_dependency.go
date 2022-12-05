/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	tpclient "github.com/threeport/threeport-go-client"
	tpapi "github.com/threeport/threeport-rest-api/pkg/api/v0"
	"gopkg.in/yaml.v2"

	"github.com/qleet/qleetctl/internal/config"
)

var wsdUpdateConfigPath string

// updateWorkloadServiceDependencyCmd represents the workload-service-dependency command
var updateWorkloadServiceDependencyCmd = &cobra.Command{
	Use:   "workload-service-dependency",
	Short: "Update a new workload service dependency",
	Long:  `Update a new workload service dependency.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		configContent, err := ioutil.ReadFile(wsdUpdateConfigPath)
		if err != nil {
			panic(err)
		}
		var workloadServiceDependencyConfig config.WorkloadServiceDependencyConfig
		if err := yaml.Unmarshal(configContent, &workloadServiceDependencyConfig); err != nil {
			panic(err)
		}

		// get workload instance by name
		workloadInstance, err := tpclient.GetWorkloadInstanceByName(
			workloadServiceDependencyConfig.WorkloadInstanceName,
			"http://localhost:1323", "",
		)
		if err != nil {
			panic(err)
		}

		// construct workload service dependency object
		workloadServiceDependency := &tpapi.WorkloadServiceDependency{
			Name:               &workloadServiceDependencyConfig.Name,
			UpstreamHost:       &workloadServiceDependencyConfig.UpstreamHost,
			UpstreamPath:       &workloadServiceDependencyConfig.UpstreamPath,
			WorkloadInstanceID: &workloadInstance.ID,
		}

		// get existing workload service dependency by name to retrieve its ID
		existingWSD, err := tpclient.GetWorkloadServiceDependencyByName(
			workloadServiceDependencyConfig.Name,
			"http://localhost:1323", "",
		)
		if err != nil {
			panic(err)
		}

		// update workload instance in API
		wsdJSON, err := json.Marshal(&workloadServiceDependency)
		if err != nil {
			panic(err)
		}
		wsd, err := tpclient.UpdateWorkloadServiceDependency(existingWSD.ID, wsdJSON, "http://localhost:1323", "")
		if err != nil {
			panic(err)
		}

		fmt.Printf("workload service dependency %s updated\n", *wsd.Name)
	},
}

func init() {
	updateCmd.AddCommand(updateWorkloadServiceDependencyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateWorkloadServiceDependencyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateWorkloadServiceDependencyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateWorkloadServiceDependencyCmd.Flags().StringVarP(&wsdUpdateConfigPath, "config", "c", "", "path to file with workload service dependency config")
	updateWorkloadServiceDependencyCmd.MarkFlagRequired("config")
}
