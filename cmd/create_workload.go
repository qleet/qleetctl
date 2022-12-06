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

var workloadCreateConfigPath string

// createWorkloadCmd represents the workload command
var createWorkloadCmd = &cobra.Command{
	Use:   "workload",
	Short: "Create a new workload",
	Long:  `Create a new workload.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load config
		configContent, err := ioutil.ReadFile(workloadCreateConfigPath)
		if err != nil {
			panic(err)
		}
		var workloadConfig config.WorkloadConfig
		if err := yaml.Unmarshal(configContent, &workloadConfig); err != nil {
			panic(err)
		}

		// create workload
		if err := workloadConfig.Create(); err != nil {
			panic(err)
		}

		fmt.Printf("workload %s created\n", workloadConfig.Name)
		//fmt.Println("############################")
		//fmt.Printf("%+v\n", workloadConfig)
	},
}

func init() {
	createCmd.AddCommand(createWorkloadCmd)

	createWorkloadCmd.Flags().StringVarP(&workloadCreateConfigPath, "config", "c", "", "path to file with workload config")
	createWorkloadCmd.MarkFlagRequired("config")
}
