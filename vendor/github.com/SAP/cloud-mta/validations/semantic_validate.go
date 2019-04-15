package validate

import (
	"github.com/SAP/cloud-mta/mta"
	"strings"
)

type checkSemantic func(mta *mta.MTA, source string) []YamlValidationIssue

// runSemanticValidations - runs semantic validations
func runSemanticValidations(mtaStr *mta.MTA, source string, exclude string) []YamlValidationIssue {
	var issues []YamlValidationIssue

	validations := getSemanticValidations(exclude)
	for _, validation := range validations {
		validationIssues := validation(mtaStr, source)
		issues = append(issues, validationIssues...)

	}
	return issues
}

// getSemanticValidations - gets list of all semantic validations minus excludes validations
func getSemanticValidations(exclude string) []checkSemantic {
	var validations []checkSemantic
	if !strings.Contains(exclude, "paths") {
		validations = append(validations, ifModulePathExists)
	}
	if !strings.Contains(exclude, "names") {
		validations = append(validations, isNameUnique)
	}
	if !strings.Contains(exclude, "requires") {
		validations = append(validations, ifRequiredDefined)
	}
	return validations
}
