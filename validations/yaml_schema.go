// TODO: HO
// 1. Position information for schema issues (blocked)
// 2. Path information for schema issues.
// 4. regExp patterns also implicitly require a NotMapSequence validation

package mta_validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/smallfish/simpleyaml"
)

// BuildValidationsFromSchemaText Entry point that accepts a Yaml Schema as text and produces YAML validation functions
// and the schema issues detected.
func BuildValidationsFromSchemaText(yaml []byte) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	y, parseError := simpleyaml.NewYaml(yaml)
	if parseError != nil {
		schemaIssues = appendIssue(schemaIssues, parseError.Error())
		return validations, schemaIssues
	}

	return buildValidationsFromSchema(y)
}

// Internal YAML validation builder
// Will be called recursively and traverse the schema structure.
func buildValidationsFromSchema(schema *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	typeNode := schema.Get("type")
	typeNodeValue, _ := typeNode.String()
	if typeNode.IsFound() && (typeNodeValue == "map" || typeNodeValue == "seq") {
		switch typeNodeValue {
		// type: map
		// mapping:
		//   firstName:  {required: true}
		//   ...
		case "map":
			mappingNode := schema.Get("mapping")
			if !mappingNode.IsMap() {
				schemaIssues = appendIssue(schemaIssues, "YAML Schema Error: <mapping> node must be a map")
				return validations, schemaIssues
			}
			newValidations, newSchemaIssues := buildValidationsFromMap(mappingNode)
			schemaIssues = append(schemaIssues, newSchemaIssues...)
			validations = append(validations, newValidations...)
			// type: seq
			// sequence:
			//  - type: map
			//  mapping:
			//    name: {required: true}
			//    ...
		case "seq":
			sequenceNode := schema.Get("sequence")
			if !sequenceNode.IsArray() {
				schemaIssues = appendIssue(schemaIssues, "YAML Schema Error: <sequence> node must be an array")
				return validations, schemaIssues
			}

			seqSize, _ := sequenceNode.GetArraySize()
			if seqSize > 1 {
				schemaIssues = appendIssue(schemaIssues, "YAML Schema Error: <sequence> node can only have one item")
				return validations, schemaIssues
			}

			// A sequence schema node has exactly only 1 element
			sequenceItemNode := sequenceNode.GetIndex(0)
			sequenceValidations, newSchemaIssues := buildValidationsFromSequence(sequenceItemNode)
			schemaIssues = append(schemaIssues, newSchemaIssues...)
			validations = append(validations, sequenceValidations...)
		}
		// {required: true, pattern: /^[a-z]+$/}
	} else {
		return buildLeafValidations(schema)
	}

	return validations, schemaIssues
}

func appendIssue(issues []YamlValidationIssue, issue string) []YamlValidationIssue {
	return append(issues, []YamlValidationIssue{{issue}}...)
}

// Create Validations for a mapping
// each key's inner validations will be wrapping in a "property" validation
func buildValidationsFromMap(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	validations = append(validations, typeIsMap())

	keys, _ := y.GetMapKeys()
	for _, key := range keys {
		value := y.Get(key)
		propInnerValidations, newSchemaIssues := buildValidationsFromSchema(value)
		schemaIssues = append(schemaIssues, newSchemaIssues...)
		propWrapperValidation := property(key, propInnerValidations...)
		validations = append(validations, propWrapperValidation)
	}
	return validations, schemaIssues
}

// Creates validations for a sequence
// Will wrap the nested checks with a "typeIsArray" check and iterate over the elements
// using "forEach"
func buildValidationsFromSequence(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	sequenceInnerValidations, newIssues := buildValidationsFromSchema(y)
	seqChecksWrapper := sequenceFailFast(typeIsArray(), forEach(sequenceInnerValidations...))
	validations = append(validations, seqChecksWrapper)
	schemaIssues = append(schemaIssues, newIssues...)

	return validations, newIssues
}

// Will create the "edge" nodes validations, these are specific checks
// for a specific path at the end of the YAML Schema
// e.g: {required: true, pattern: /^[a-zA-Z]$/}
func buildLeafValidations(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	validations, schemaIssues = invokeLeafValidation(y, validations, schemaIssues, buildTypeValidation)

	validations, schemaIssues = invokeLeafValidation(y, validations, schemaIssues, buildPatternValidation)

	// Special handling is needed for "optional" and "required" Validations, must be invoked last
	// and receive all previous built validations as extra wrapping will be done here.
	return buildOptionalOrRequiredValidation(y, validations, schemaIssues)

}

func buildOptionalOrRequiredValidation(y *simpleyaml.Yaml, validations []yamlCheck, schemaIssues []YamlValidationIssue) ([]yamlCheck, []YamlValidationIssue) {
	requiredNode := y.Get("required")
	if requiredNode.IsFound() {
		requiredValue := getLiteralStringValue(requiredNode)
		if requiredValue != "true" && requiredValue != "false" {
			schemaIssues = appendIssue(schemaIssues,
				fmt.Sprintf("YAML Schema Error: <required> node must be a boolean but found <%s>", requiredValue))
			return validations, schemaIssues
		}

		// When a "required" check is needed there is no need to perform additional validations
		// if the property is missing, thus we use "sequenceFailFast"
		if requiredValue == "true" {
			// The required check must be performed first in the sequence.
			validationsWithRequiredFirst := append([]yamlCheck{required()}, validations...)
			validations = []yamlCheck{sequenceFailFast(validationsWithRequiredFirst...)}
			// When "required" is false, we must only perform additional validations
			// if the property actually exists, thus we use "optional"
		} else {
			// An optional wraps all our other validations.
			validations = []yamlCheck{optional(validations...)}
		}
	}
	return validations, schemaIssues
}

func buildTypeValidation(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	typeNode := y.Get("type")
	if typeNode.IsFound() {
		typeValue, stringErr := typeNode.String()
		if stringErr != nil {
			schemaIssues = appendIssue(schemaIssues, "YAML Schema Error: <type> node must be a string")
			return validations, schemaIssues
		}
		if typeValue == "bool" {
			validations = append(validations, typeIsBoolean())
		} else if typeValue == "enum" {
			enumValidations, enumSchemaIssues := buildEnumValidation(y)
			validations = append(validations, enumValidations...)
			schemaIssues = append(schemaIssues, enumSchemaIssues...)
		}
	}
	return validations, schemaIssues
}

func buildEnumValidation(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	enumsNode := y.Get("enums")
	if !enumsNode.IsFound() {
		return []yamlCheck{}, []YamlValidationIssue{{"YAML Schema Error: enums values must be listed"}}
	}
	if !enumsNode.IsArray() {
		return []yamlCheck{}, []YamlValidationIssue{{"YAML Schema Error: enums values must be listed as array"}}
	}

	enumsNumber, _ := enumsNode.GetArraySize()

	enumValues := []string{}
	for i := 0; i < enumsNumber; i++ {
		enumNode := enumsNode.GetIndex(i)
		if enumNode.IsArray() || enumNode.IsMap() {
			return []yamlCheck{}, []YamlValidationIssue{{"YAML Schema Error: enum values must be simple"}}
		} else {
			enumValue := getLiteralStringValue(enumNode)
			enumValues = append(enumValues, enumValue)
		}
	}

	return []yamlCheck{matchesEnumValues(enumValues)}, []YamlValidationIssue{}
}

func buildPatternValidation(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue) {
	var validations []yamlCheck
	var schemaIssues []YamlValidationIssue

	patternNode := y.Get("pattern")
	if patternNode.IsFound() {
		patternValue, err := patternNode.String()
		if err != nil {
			schemaIssues = append(schemaIssues, YamlValidationIssue{"YAML Schema Error: <pattern> node must be a string"})
			return validations, schemaIssues
		}
		// TODO: we must validate: NOT MAP/SEQ
		patternWithoutSlashes := strings.TrimSuffix(strings.TrimPrefix(patternValue, "/"), "/")
		_, err = regexp.Compile(patternWithoutSlashes)
		if err != nil {
			schemaIssues = append(schemaIssues, YamlValidationIssue{"YAML Schema Error: <pattern> node not valid: " + err.Error()})
		} else {
			validations = append(validations, matchesRegExp(patternWithoutSlashes))
		}
	}
	return validations, schemaIssues
}

// Utility to reduce verbosity
func invokeLeafValidation(y *simpleyaml.Yaml, validations []yamlCheck, schemaIsssues []YamlValidationIssue, leafBuilder func(y *simpleyaml.Yaml) ([]yamlCheck, []YamlValidationIssue)) ([]yamlCheck, []YamlValidationIssue) {
	newValidations, newSchemaIssues := leafBuilder(y)
	validations = append(validations, newValidations...)
	schemaIsssues = append(schemaIsssues, newSchemaIssues...)

	return validations, schemaIsssues
}
