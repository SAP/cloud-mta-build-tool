package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/SAP/cloud-mta/mta"
)

// ifRequiredDefined - validates that required property sets are defined in modules, provided sections or resources
func ifRequiredDefined(mta *mta.MTA, source string) []YamlValidationIssue {
	var issues []YamlValidationIssue

	// init set of all provided property sets
	provided := make(map[string]map[string]interface{})

	for _, module := range mta.Modules {
		// add module to provided property sets
		provided[module.Name] = module.Properties
		// add all property sets provided by module
		for _, prov := range module.Provides {
			provided[prov.Name] = prov.Properties
		}
	}

	// add resources to provided property sets
	for _, resource := range mta.Resources {
		provided[resource.Name] = resource.Properties
	}

	for _, module := range mta.Modules {
		issues = append(issues, checkComponent(provided, module, "module")...)
		for _, moduleProvides := range module.Provides {
			issues = append(issues, checkComponent(provided, &moduleProvides, "provided property set of the "+module.Name+" module")...)
		}
	}

	for _, resource := range mta.Resources {
		issues = append(issues, checkComponent(provided, resource, "resource")...)
	}
	return issues
}

func structFieldToMap(str interface{}, field string) map[string]interface{} {
	v := reflect.ValueOf(str).Elem().FieldByName(field)
	if v.IsValid() {
		mapValue, ok := v.Addr().Interface().(*map[string]interface{})
		if ok {
			return *mapValue
		}
	}
	return nil
}

//
func structFieldToRequires(str interface{}) []mta.Requires {
	v := reflect.ValueOf(str).Elem().FieldByName("Requires")
	if v.IsValid() {
		mapValue, ok := v.Addr().Interface().(*[]mta.Requires)
		if ok {
			return *mapValue
		}
	}
	return []mta.Requires{}
}

func structFieldToString(str interface{}) string {
	v := reflect.ValueOf(str).Elem().FieldByName("Name")
	strValue := v.Addr().Interface().(*string)
	return *strValue
}

func checkComponent(provided map[string]map[string]interface{}, component interface{}, compDesc string) []YamlValidationIssue {
	var issues []YamlValidationIssue

	compName := structFieldToString(component)
	issues = append(issues,
		checkRequiredProperties(provided, "", structFieldToMap(component, "Properties"),
			fmt.Sprintf(`"%s" %s`, compName, compDesc), "property")...)
	issues = append(issues,
		checkRequiredProperties(provided, "", structFieldToMap(component, "Parameters"),
			fmt.Sprintf(`"%s" %s`, compName, compDesc), "parameter")...)
	issues = append(issues,
		checkRequiredProperties(provided, "", structFieldToMap(component, "BuildParams"),
			fmt.Sprintf(`"%s" %s`, compName, compDesc), "build parameter")...)
	// check that each required by resource property set was provided in mta.yaml
	for _, requires := range structFieldToRequires(component) {
		if _, contains := provided[requires.Name]; !contains {
			issues = appendIssue(issues,
				fmt.Sprintf(`the "%s" property set required by the "%s" %s is not defined`,
					requires.Name, compName, compDesc))
		}
		// check that each property of resource is resolved
		issues = append(issues,
			checkRequiredProperties(provided, requires.Name, requires.Properties,
				fmt.Sprintf(`"%s" %s`, compName, compDesc), "property")...)
		// check that each parameter of resource is resolved
		issues = append(issues,
			checkRequiredProperties(provided, requires.Name, requires.Parameters,
				fmt.Sprintf(`"%s" %s`, compName, compDesc), "parameter")...)
	}
	return issues
}

func checkRequiredProperties(providedProps map[string]map[string]interface{}, requiredPropSet string,
	requiredEntities map[string]interface{}, requiringObject, entityKind string) []YamlValidationIssue {

	var issues []YamlValidationIssue
	if requiredEntities == nil {
		return nil
	}
	for entityName, entityValue := range requiredEntities {
		issues = append(issues, checkValue(providedProps, entityName, entityKind, requiredPropSet, requiringObject, entityValue)...)
	}
	return issues
}

func checkValue(providedProps map[string]map[string]interface{},
	entityName, entityKind, propSet, requiringObject string, entityValue interface{}) []YamlValidationIssue {
	var issues []YamlValidationIssue
	propValueStr, ok := entityValue.(string)
	if ok {
		// property is simple - check if it can be resolved
		issues = checkStringEntityValue(providedProps, entityName, propValueStr, entityKind, propSet, requiringObject)
	} else {
		propValueMap, ok := entityValue.(map[interface{}]interface{})
		if ok {
			// property is a map
			for key, value := range propValueMap {
				// check every sub property
				issues = append(issues, checkValue(providedProps, entityName+"."+key.(string), entityKind, propSet, requiringObject, value)...)
			}
		}
	}
	return issues
}

func checkStringEntityValue(providedProps map[string]map[string]interface{},
	entityName, entityValue, entityKind, propSet, requiringObject string) []YamlValidationIssue {
	var issues []YamlValidationIssue
	r := regexp.MustCompile(`~{[^{}]+}`)
	// find all placeholders
	matches := r.FindAllString(entityValue, -1)
	for _, match := range matches {
		// get placeholder pure name, removing tilda and brackets
		requiredProp := strings.TrimPrefix(strings.TrimSuffix(match, "}"), "~{")
		// if property set was not provided it has to be presented in placeholder
		if propSet == "" {
			// split placeholder to property set and property name
			requiredPropArr := strings.SplitN(requiredProp, "/", 2)
			if len(requiredPropArr) != 2 {
				// no property set provided
				issues = appendIssue(issues,
					fmt.Sprintf(`the "%s" %s of the %s is unresolved; the "%s" property is not provided`,
						entityName, entityKind, requiringObject, requiredProp))
			} else {
				// check existence of property if property set
				issues = appendIssue(issues,
					checkRequiredProperty(providedProps, entityName, entityKind, requiredPropArr[0], requiredPropArr[1], requiringObject))
			}
		} else {
			// check existence of property if property set
			issues = appendIssue(issues,
				checkRequiredProperty(providedProps, entityName, entityKind, propSet, requiredProp, requiringObject))
		}

	}
	return issues
}

func checkRequiredProperty(providedProps map[string]map[string]interface{}, entityName, entityKind,
	requiredSet, requiredProp, requiringObject string) string {
	providedSet, ok := providedProps[requiredSet]
	if ok {
		_, ok = providedSet[requiredProp]
	}
	if !ok {
		return fmt.Sprintf(`the "%s" %s of the %s is unresolved; the "%s/%s" property is not provided`,
			entityName, entityKind, requiringObject, requiredSet, requiredProp)
	}
	return ""
}
