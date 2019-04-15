package validate

import (
	"fmt"
	"github.com/SAP/cloud-mta/mta"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetValidationMode converts validation mode flags to validation process flags.
func GetValidationMode(validationFlag string) (bool, bool, error) {
	switch validationFlag {
	case "schema":
		return true, false, nil
	case "semantic":
		return true, true, nil
	case "":
		return true, true, nil
	}
	return false, false,
		fmt.Errorf(`the "%s" validation mode is incorrect; expected one of the following: schema, semantic`,
			validationFlag)
}

// MtaYaml validates an MTA.yaml file.
func MtaYaml(projectPath, mtaFilename string,
	validateSchema, validateSemantic, strict bool, exclude string) (warning string, err error) {
	if validateSemantic || validateSchema {

		mtaPath := filepath.Join(projectPath, mtaFilename)
		// ParseFile contains MTA yaml content.
		yamlContent, e := ioutil.ReadFile(mtaPath)

		if e != nil {
			return "", errors.Wrapf(e, `could not read the "%v" file; the validation failed`, mtaPath)
		}
		// Validates MTA content.
		errIssues, warnIssues := validate(yamlContent, projectPath,
			validateSchema, validateSemantic, strict, exclude, yaml.Unmarshal)
		if len(errIssues) > 0 {
			return warnIssues.String(), errors.Errorf(`the "%v" file is not valid: `+"\n%v",
				mtaPath, errIssues.String())
		}
		return warnIssues.String(), nil
	}

	return "", nil
}

// validate - validates the MTA descriptor
func validate(yamlContent []byte, projectPath string,
	validateSchema, validateSemantic, strict bool, exclude string,
	unmarshal func(mtaContent []byte, mtaStr interface{}) error) (errIssues YamlValidationIssues, warnIssues YamlValidationIssues) {

	mtaStr := mta.MTA{}

	err := yaml.UnmarshalStrict(yamlContent, &mtaStr)
	if strict && err != nil {
		errIssues = appendIssue(errIssues, err.Error())
	} else if err != nil {
		warnIssues = appendIssue(warnIssues, err.Error())
		err = unmarshal(yamlContent, &mtaStr)
		if err != nil {
			errIssues = appendIssue(errIssues, err.Error())
			return errIssues, nil
		}
	}

	if validateSchema {
		validations, schemaValidationLog := buildValidationsFromSchemaText(schemaDef)
		if len(schemaValidationLog) > 0 {
			errIssues = append(errIssues, schemaValidationLog...)
			return errIssues, warnIssues
		}
		errIssues = append(errIssues, runSchemaValidations(yamlContent, validations...)...)
	}

	if validateSemantic {
		errIssues = append(errIssues, runSemanticValidations(&mtaStr, projectPath, exclude)...)
	}
	return errIssues, warnIssues
}
