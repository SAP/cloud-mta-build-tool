package validate

import (
	"fmt"

	"github.com/SAP/cloud-mta/mta"
)

// isNameUnique - validates the global uniqueness of the names of modules, provided services and resources
func isNameUnique(mta *mta.MTA, source string) []YamlValidationIssue {
	var issues []YamlValidationIssue
	// map: name -> object kind (module, provided services or resource)
	names := make(map[string]string)
	for _, module := range mta.Modules {
		// validate module name
		issues = validateNameUniqueness(names, module.Name, "module", issues)
		for _, provide := range module.Provides {
			// validate name of provided service
			issues = validateNameUniqueness(names, provide.Name, "provided property set", issues)
		}
	}
	for _, resource := range mta.Resources {
		// validate resource name
		issues = validateNameUniqueness(names, resource.Name, "resource", issues)
	}
	return issues
}

// validateNameUniqueness - validate that name not defined already (not exists in the 'names' map)
func validateNameUniqueness(names map[string]string, name string,
	objectName string, issues []YamlValidationIssue) []YamlValidationIssue {
	result := issues
	// try to find name in the global map
	prevObjectName, ok := names[name]
	// name found -> add issue
	if ok {
		var article string
		if objectName == prevObjectName {
			article = "another"
		} else {
			article = "a"
		}
		result = appendIssue(result,
			fmt.Sprintf(`the "%s" %s name is not unique; %s %s was found with the same name`,
				name, objectName, article, prevObjectName))
	} else {
		// name not found -> add it to the global map
		names[name] = objectName
	}
	return result
}
