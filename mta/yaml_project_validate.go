package mta

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"cloud-mta-build-tool/validations"
)

type YamlProjectCheck func(mta MTA, path string) []mta_validate.YamlValidationIssue

// validateModules - Validate the MTA file
func validateModules(mta MTA, projectPath string) []mta_validate.YamlValidationIssue {
	issues := []mta_validate.YamlValidationIssue{}
	for _, module := range mta.Modules {
		modulePath := module.Path
		if modulePath == "" {
			modulePath = module.Name
		}
		dirName := filepath.Join(projectPath, modulePath)
		_, err := ioutil.ReadDir(dirName)
		if err != nil {
			issues = append(issues, []mta_validate.YamlValidationIssue{{fmt.Sprintf("Module <%s> not found in project. Expected path: <%s>", module.Name, modulePath)}}...)
		}
	}

	return issues
}

// ValidateYamlProject - Validate the MTA file
func ValidateYamlProject(mta MTA, path string) []mta_validate.YamlValidationIssue {
	validations := []YamlProjectCheck{validateModules}
	issues := []mta_validate.YamlValidationIssue{}
	for _, validation := range validations {
		validationIssues := validation(mta, path)
		issues = append(issues, validationIssues...)

	}
	return issues
}
