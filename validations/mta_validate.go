package validate

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/mta"
)

// GetValidationMode - convert validation mode flag to validation process flags
func GetValidationMode(validationFlag string) (bool, bool, error) {
	switch validationFlag {
	case "":
		return true, true, nil
	case "schema":
		return true, false, nil
	case "project":
		return false, true, nil
	}
	return false, false, fmt.Errorf("wrong validation mode <%v> (expected one of [all, schema, project])", validationFlag)
}

// MtaYaml - Validate MTA yaml
func MtaYaml(projectPath, mtaFilename string, validateSchema bool, validateProject bool) error {
	if validateProject || validateSchema {

		mtaPath := filepath.Join(projectPath, mtaFilename)
		// ParseFile MTA yaml content
		yamlContent, err := ioutil.ReadFile(mtaPath)

		if err != nil {
			return errors.Wrapf(err, "validation of %v failed when reading the file", mtaPath)
		}
		// validate mta content
		issues, err := validate(yamlContent, projectPath, validateSchema, validateProject)
		if len(issues) > 0 {
			return errors.Errorf("validation of %v failed with the following issues: \n%v %s", mtaPath, issues.String(), err)
		}
	}

	return nil
}

// validate validates an MTA schema.
func validate(yamlContent []byte, projectPath string, validateSchema bool, validateProject bool) (YamlValidationIssues, error) {
	//noinspection GoPreferNilSlice
	issues := []YamlValidationIssue{}
	if validateSchema {
		validations, schemaValidationLog := BuildValidationsFromSchemaText(schemaDef)
		if len(schemaValidationLog) > 0 {
			return schemaValidationLog, nil
		}
		yamlValidationLog, err := Yaml(yamlContent, validations...)
		if err != nil && len(yamlValidationLog) == 0 {
			yamlValidationLog = append(yamlValidationLog, []YamlValidationIssue{{Msg: "validation failed with error: " + err.Error()}}...)
		}
		issues = append(issues, yamlValidationLog...)

	}
	if validateProject {
		mtaStr := mta.MTA{}
		Unmarshal := yaml.Unmarshal
		err := Unmarshal(yamlContent, &mtaStr)
		if err != nil {
			return nil, errors.Wrap(err, "validation failed when unmarshalling mta")
		}
		projectIssues := validateYamlProject(&mtaStr, projectPath)
		issues = append(issues, projectIssues...)
	}
	return issues, nil
}
