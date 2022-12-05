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

var wsdCreateConfigPath string

// createWorkloadServiceDependencyCmd represents the workload-service-dependency command
var createWorkloadServiceDependencyCmd = &cobra.Command{
	Use:   "workload-service-dependency",
	Short: "Create a new workload service dependency",
	Long:  `Create a new workload service dependency.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		configContent, err := ioutil.ReadFile(wsdCreateConfigPath)
		if err != nil {
			panic(err)
		}
		var workloadServiceDependencyConfig config.WorkloadServiceDependencyConfig
		if err := yaml.Unmarshal(configContent, &workloadServiceDependencyConfig); err != nil {
			panic(err)
		}

		// create workload service dependency
		workloadServiceDependency, err := workloadServiceDependencyConfig.Create()
		if err != nil {
			panic(err)
		}

		//// get workload instance by name
		//workloadInstance, err := tpclient.GetWorkloadInstanceByName(
		//	workloadServiceDependencyConfig.WorkloadInstanceName,
		//	"http://localhost:1323", "",
		//)
		//if err != nil {
		//	panic(err)
		//}

		//// construct workload service dependency object
		//workloadServiceDependency := &tpapi.WorkloadServiceDependency{
		//	Name:               &workloadServiceDependencyConfig.Name,
		//	UpstreamHost:       &workloadServiceDependencyConfig.UpstreamHost,
		//	UpstreamPath:       &workloadServiceDependencyConfig.UpstreamPath,
		//	WorkloadInstanceID: &workloadInstance.ID,
		//}

		//// create workload instance in API
		//wsdJSON, err := json.Marshal(&workloadServiceDependency)
		//if err != nil {
		//	panic(err)
		//}
		//wsd, err := tpclient.CreateWorkloadServiceDependency(wsdJSON, "http://localhost:1323", "")
		//if err != nil {
		//	panic(err)
		//}

		fmt.Printf("workload service dependency %s created\n", *workloadServiceDependency.Name)
	},
}

func init() {
	createCmd.AddCommand(createWorkloadServiceDependencyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createWorkloadServiceDependencyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createWorkloadServiceDependencyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createWorkloadServiceDependencyCmd.Flags().StringVarP(&wsdCreateConfigPath, "config", "c", "", "path to file with workload service dependency config")
	createWorkloadServiceDependencyCmd.MarkFlagRequired("config")
}
