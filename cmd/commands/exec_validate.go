package commands

import (
	"errors"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

func getValidationMode(validationFlag string) (bool, bool, error) {
	switch validationFlag {
	case "":
		return true, true, nil
	case "schema":
		return true, false, nil
	case "project":
		return false, true, nil
	}
	return false, false, errors.New("wrong argument of validation mode. Expected one of [all, schema, project]")
}

func validateMtaYaml(yamlPath string, yamlFilename string, validateSchema bool, validateProject bool) error {
	if validateProject || validateSchema {
		logs.Logger.Info("Starting MTA Yaml validation")
		yamlContent, err := mta.ReadMtaContent(yamlPath, yamlFilename)

		if err != nil {
			return errors.New("MTA validation failed. " + err.Error())
		}
		projectPath, err := dir.GetCurrentPath()
		if err != nil {
			return errors.New("MTA validation failed. " + err.Error())
		} else {
			issues := mta.Validate(yamlContent, projectPath, validateSchema, validateProject)
			valid := len(issues) == 0
			if valid {
				logs.Logger.Info("MTA Yaml is valid")
			} else {
				return errors.New("MTA Yaml is  invalid. Issues: \n" + issues.String())
			}
		}
	}
	return nil
}
