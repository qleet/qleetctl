package config

import (
	"encoding/json"
	"io/ioutil"

	tpclient "github.com/threeport/threeport-go-client"
	tpapi "github.com/threeport/threeport-rest-api/pkg/api/v0"
)

type WorkloadConfig struct {
	Name                      string                          `yaml:"Name"`
	WorkloadDefinition        WorkloadDefinitionConfig        `yaml:"WorkloadDefinition"`
	WorkloadInstance          WorkloadInstanceConfig          `yaml:"WorkloadInstance"`
	WorkloadServiceDependency WorkloadServiceDependencyConfig `yaml:"WorkloadServiceDependency"`
}

type WorkloadDefinitionConfig struct {
	Name         string `yaml:"Name"`
	YAMLDocument string `yaml:"YAMLDocument"`
	UserID       uint   `yaml:"UserID"`
}

type WorkloadInstanceConfig struct {
	Name                   string `yaml:"Name"`
	WorkloadClusterName    string `yaml:"WorkloadClusterName"`
	WorkloadDefinitionName string `yaml:"WorkloadDefinitionName"`
}

type WorkloadServiceDependencyConfig struct {
	Name                 string `yaml:"Name"`
	UpstreamHost         string `yaml:"UpstreamHost"`
	UpstreamPath         string `yaml:"UpstreamPath"`
	WorkloadInstanceName string `yaml:"WorkloadInstanceName"`
}

func (wc *WorkloadConfig) Create() error {
	// create the definition
	_, aerr := wc.WorkloadDefinition.Create()
	if aerr != nil {
		return aerr
	}

	// create the instance
	_, berr := wc.WorkloadInstance.Create()
	if berr != nil {
		return berr
	}

	// create the service dependency
	_, cerr := wc.WorkloadServiceDependency.Create()
	if cerr != nil {
		return cerr
	}

	return nil
}

func (wdc *WorkloadDefinitionConfig) Create() (*tpapi.WorkloadDefinition, error) {
	// get the content of the yaml document
	definitionContent, err := ioutil.ReadFile(wdc.YAMLDocument)
	if err != nil {
		return nil, err
	}
	stringContent := string(definitionContent)

	// construct workload definition object
	workloadDefinition := &tpapi.WorkloadDefinition{
		Name:         &wdc.Name,
		YAMLDocument: &stringContent,
		UserID:       &wdc.UserID,
	}

	// create workload definition in API
	wdJSON, err := json.Marshal(&workloadDefinition)
	if err != nil {
		return nil, err
	}
	wd, err := tpclient.CreateWorkloadDefinition(wdJSON, "http://localhost:1323", "")
	if err != nil {
		return nil, err
	}

	return wd, nil
}

func (wic *WorkloadInstanceConfig) Create() (*tpapi.WorkloadInstance, error) {
	// get workload cluster by name
	workloadCluster, err := tpclient.GetWorkloadClusterByName(
		wic.WorkloadClusterName,
		"http://localhost:1323", "",
	)
	if err != nil {
		return nil, err
	}

	// get workload definition by name
	workloadDefinition, err := tpclient.GetWorkloadDefinitionByName(
		wic.WorkloadDefinitionName,
		"http://localhost:1323", "",
	)
	if err != nil {
		return nil, err
	}

	// construct workload instance object
	workloadInstance := &tpapi.WorkloadInstance{
		Name:                 &wic.Name,
		WorkloadClusterID:    &workloadCluster.ID,
		WorkloadDefinitionID: &workloadDefinition.ID,
	}

	// create workload instance in API
	wiJSON, err := json.Marshal(&workloadInstance)
	if err != nil {
		return nil, err
	}
	wi, err := tpclient.CreateWorkloadInstance(wiJSON, "http://localhost:1323", "")
	if err != nil {
		return nil, err
	}

	return wi, nil
}

func (wsdc *WorkloadServiceDependencyConfig) Create() (*tpapi.WorkloadServiceDependency, error) {
	// get workload instance by name
	workloadInstance, err := tpclient.GetWorkloadInstanceByName(
		wsdc.WorkloadInstanceName,
		"http://localhost:1323", "",
	)
	if err != nil {
		panic(err)
	}

	// construct workload service dependency object
	workloadServiceDependency := &tpapi.WorkloadServiceDependency{
		Name:               &wsdc.Name,
		UpstreamHost:       &wsdc.UpstreamHost,
		UpstreamPath:       &wsdc.UpstreamPath,
		WorkloadInstanceID: &workloadInstance.ID,
	}

	// create workload instance in API
	wsdJSON, err := json.Marshal(&workloadServiceDependency)
	if err != nil {
		panic(err)
	}
	wsd, err := tpclient.CreateWorkloadServiceDependency(wsdJSON, "http://localhost:1323", "")
	if err != nil {
		panic(err)
	}

	return wsd, nil
}