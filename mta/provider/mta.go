package provider

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/constants"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta/models"
)

// MTA - provide full mta object
func MTA(source string) (models.MTA, error) {

	// Read the MTA
	yamlFile, err := ioutil.ReadFile(source + constants.PathSep + "mta.yaml")
	mta := models.MTA{}
	if err != nil {
		logs.Logger.Error("Not able to read the mta file ")
		return mta, err
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &mta)
	if err != nil {
		logs.Logger.Error("Not able to unmarshal the mta file ")
		return mta, err
	}
	return mta, err
}

// Provide list of modules
func modules(mta models.MTA) []string {
	var mNames []string
	for _, mod := range mta.Modules {
		mNames = append(mNames, mod.Name)
	}
	return mNames
}

// GetModulesNames - get list of modules names
func GetModulesNames(source string) []string {

	mta, e := MTA(source)
	if e != nil {
		logs.Logger.Error(e)
	}
	return modules(mta)
}
