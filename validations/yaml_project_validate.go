package validate

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"cloud-mta-build-tool/mta"
)

type yamlProjectCheck func(mta *mta.MTA, path string) []YamlValidationIssue

// validateModules - Validate the MTA file
func validateModules(mta *mta.MTA, projectPath string) []YamlValidationIssue {
	//noinspection GoPreferNilSlice
	issues := []YamlValidationIssue{}
	for _, module := range mta.Modules {
		modulePath := module.Path
		if modulePath == "" {
			modulePath = module.Name
		}
		dirName := filepath.Join(projectPath, modulePath)
		_, err := ioutil.ReadDir(dirName)
		if err != nil {
			issues = append(issues, []YamlValidationIssue{{Msg: fmt.Sprintf("Module <%s> not found in project. Expected path: <%s>", module.Name, modulePath)}}...)
		}
	}

	return issues
}

// validateYamlProject - Validate the MTA file
func validateYamlProject(mta *mta.MTA, path string) []YamlValidationIssue {
	validations := []yamlProjectCheck{validateModules}
	//noinspection GoPreferNilSlice
	issues := []YamlValidationIssue{}
	for _, validation := range validations {
		validationIssues := validation(mta, path)
		issues = append(issues, validationIssues...)

	}
	return issues
}
