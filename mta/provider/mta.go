package provider

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/constants"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
)

// MTA - provide full mta object
func MTA(source string) (mta.MTA, error) {

	// Read the MTA
	yamlFile, err := ioutil.ReadFile(source + constants.PathSep + "mta.yaml")
	m := mta.MTA{}
	if err != nil {
		logs.Logger.Error("Not able to read the mta file ")
		return m, err
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &m)
	if err != nil {
		logs.Logger.Error("Not able to unmarshal the mta file ")
		return m, err
	}
	return m, err
}

// Provide list of modules
func modules(mta mta.MTA) []string {
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
