// TODO: HO
// 1. Position information for schema issues (blocked)
// 2. Path information for schema issues.
// 3. Validate that regExp patterns are valid
// 4. regExp patterns also implicitly require a NotMapSequence validation
// 5. TypeValidations (Bool / Enum)

// TODO: Shahar
// 1. Comments
// 2. extract error prefix

package mta_validate

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"strings"
)

func BuildValidationsFromSchemaText(yaml []byte) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	y, parseError := simpleyaml.NewYaml(yaml)
	if parseError != nil {
		schemaIssues = append(schemaIssues, parseError.Error())
		return validations, schemaIssues
	}

	return BuildValidationsFromSchema(y)
}

func BuildValidationsFromSchema(schema *simpleyaml.Yaml) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	typeNode := schema.Get("type")
	typeNodeValue, _ := typeNode.String()
	if typeNode.IsFound() && (typeNodeValue == "map" || typeNodeValue == "seq") {
		switch typeNodeValue {
		case "map":
			mappingNode := schema.Get("mapping")
			if !mappingNode.IsMap() {
				schemaIssues = append(schemaIssues, "YAML Schema Error: <mapping> node must be a map")
				return validations, schemaIssues
			}
			newValidations, newSchemaIssues := BuildValidationsFromMap(mappingNode)
			schemaIssues = append(schemaIssues, newSchemaIssues...)
			validations = append(validations, newValidations...)
		case "seq":
			sequenceNode := schema.Get("sequence")
			if !sequenceNode.IsArray() {
				schemaIssues = append(schemaIssues, "YAML Schema Error: <sequence> node must be an array")
				return validations, schemaIssues
			}

			seqSize, _ := sequenceNode.GetArraySize()
			if seqSize > 1 {
				schemaIssues = append(schemaIssues, "YAML Schema Error: <sequence> node can only have one item")
				return validations, schemaIssues
			}

			sequenceItemNode := sequenceNode.GetIndex(0)
			sequenceValidations, newSchemaIssues := BuildValidationsFromSequence(sequenceItemNode)
			schemaIssues = append(schemaIssues, newSchemaIssues...)
			validations = append(validations, sequenceValidations...)
		}
	} else {
		return buildLeafValidations(schema)
	}

	return validations, schemaIssues
}

func BuildValidationsFromMap(y *simpleyaml.Yaml) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	validations = append(validations, TypeIsMap())

	keys, _ := y.GetMapKeys()
	for _, key := range keys {
		value := y.Get(key)
		propInnerValidations, newSchemaIssues := BuildValidationsFromSchema(value)
		schemaIssues = append(schemaIssues, newSchemaIssues...)
		propWrapperValidation := Property(key, propInnerValidations...)
		validations = append(validations, propWrapperValidation)
	}
	return validations, schemaIssues
}

func BuildValidationsFromSequence(y *simpleyaml.Yaml) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	sequenceInnerValidations, newIssues := BuildValidationsFromSchema(y)
	validations = append(validations, SequenceFailFast(TypeIsArray(), ForEach(sequenceInnerValidations...)))
	schemaIssues = append(schemaIssues, newIssues...)

	return validations, newIssues
}

func buildLeafValidations(y *simpleyaml.Yaml) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	validations, schemaIssues = invokeLeafValidation(y, validations, schemaIssues, buildTypeValidation)

	validations, schemaIssues = invokeLeafValidation(y, validations, schemaIssues, buildPatternValidation)

	// Special handling is needed for "Optional" and "Required" Validations, must be invoked last
	// and receive all previous built validations as wrapping will be done here.
	return buildOptionalOrRequiredValidation(y, validations, schemaIssues)

}

func buildOptionalOrRequiredValidation(y *simpleyaml.Yaml, validations []YamlCheck, schemaIssues []string) ([]YamlCheck, []string) {
	requiredNode := y.Get("required")
	if requiredNode.IsFound() {
		requiredValue := getLiteralStringValue(requiredNode)
		if requiredValue != "true" && requiredValue != "false" {
			schemaIssues = append(schemaIssues,
				fmt.Sprintf("YAML Schema Error: <required> node must be a boolean but found <%s>", requiredValue))
			return validations, schemaIssues
		}

		if requiredValue == "true" {
			// The required check must be performed first in the sequence.
			validationsWithRequiredFirst := append([]YamlCheck{Required()}, validations...)
			validations = []YamlCheck{SequenceFailFast(validationsWithRequiredFirst...)}
		} else {
			// An Optional wraps all our other validations.
			validations = []YamlCheck{Optional(validations...)}
		}
	}
	return validations, schemaIssues
}

func buildTypeValidation(y *simpleyaml.Yaml) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	typeNode := y.Get("type")
	if typeNode.IsFound() {
		typeValue, stringErr := typeNode.String()
		if stringErr != nil {
			schemaIssues = append(schemaIssues, "YAML Schema Error: <type> node must be a string")
			return validations, schemaIssues
		}
		if typeValue == "bool" {
			// TODO: TBD
		} else if typeValue == "enum" {
			// TODO: TBD
		}
	}
	return validations, schemaIssues
}

func buildPatternValidation(y *simpleyaml.Yaml) ([]YamlCheck, []string) {
	var validations []YamlCheck
	var schemaIssues []string

	patternNode := y.Get("pattern")
	if patternNode.IsFound() {
		patternValue, err := patternNode.String()
		if err != nil {
			schemaIssues = append(schemaIssues, "YAML Schema Error: <pattern> node must be a string")
			return validations, schemaIssues
		}
		// TODO: we must validate: NOT MAP/SEQ
		// TODO: validate that the pattern is valid
		patternWithoutSlashes := strings.TrimSuffix(strings.TrimPrefix(patternValue, "/"), "/")
		validations = append(validations, MatchesRegExp(patternWithoutSlashes))
	}
	return validations, schemaIssues
}

func invokeLeafValidation(y *simpleyaml.Yaml, validations []YamlCheck, schemaIsssues []string, leafBuilder func(y *simpleyaml.Yaml) ([]YamlCheck, []string)) ([]YamlCheck, []string) {
	newValidations, newSchemaIssues := leafBuilder(y)
	validations = append(validations, newValidations...)
	schemaIsssues = append(schemaIsssues, newSchemaIssues...)

	return validations, schemaIsssues
}
