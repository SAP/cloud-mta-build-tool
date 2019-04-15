package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/smallfish/simpleyaml"
)

// YamlValidationIssue - specific issue
type YamlValidationIssue struct {
	Msg string
}

// YamlValidationIssues - list of issue's
type YamlValidationIssues []YamlValidationIssue

func (issues YamlValidationIssues) String() string {
	var messages []string
	for _, issue := range issues {
		messages = append(messages, issue.Msg)
	}
	return strings.Join(messages, "\n")
}

// YamlCheck - validation check function type
type YamlCheck func(y *simpleyaml.Yaml, path []string) YamlValidationIssues

// DSL method to execute validations on a sub node(property) of a YAML tree.
// Can be nested to check properties farther and farther down the tree.
func property(propName string, checks ...YamlCheck) YamlCheck {
	return func(y *simpleyaml.Yaml, path []string) YamlValidationIssues {
		var issues YamlValidationIssues
		yProp := y.Get(propName)

		// Will perform all the validations without stopping
		for _, check := range checks {
			newIssues := check(yProp, append(path, propName))
			issues = append(issues, newIssues...)
		}

		return issues
	}
}

// DSL method to execute validations in order and break early as soon as the first one fails
// This is very useful if a certain validation cannot be executed without the previous ones succeeding.
// For example: matching vs a regExp should not be performed for a property that is not a string.
func sequence(
	checks ...YamlCheck) YamlCheck {

	return sequenceInternal(false, checks...)
}

// DSL method to execute validations in order and break early as soon as the first one fails
// This is very useful if a certain validation cannot be executed without the previous ones succeeding.
// For example: matching vs a regExp should not be performed for a property that is not a string.
func sequenceFailFast(
	checks ...YamlCheck) YamlCheck {

	return sequenceInternal(true, checks...)
}

func sequenceInternal(failfast bool,
	checks ...YamlCheck) YamlCheck {

	return func(y *simpleyaml.Yaml, path []string) YamlValidationIssues {
		var issues YamlValidationIssues

		for _, check := range checks {
			newIssues := check(y, path)
			// Only perform the next validation, if the previous one succeeded
			if len(newIssues) > 0 {
				issues = append(issues, newIssues...)
				if failfast {
					break
				}
			}
		}

		return issues
	}
}

// DSL method to iterate over a YAML array items
func forEach(checks ...YamlCheck) YamlCheck {

	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		arrSize, _ := yProp.GetArraySize()

		var issues YamlValidationIssues

		validation := sequence(checks...)

		for i := 0; i < arrSize; i++ {
			yElem := yProp.GetIndex(i)
			elemErrors := validation(yElem, append(path, fmt.Sprintf("[%d]", i)))
			issues = append(issues, elemErrors...)
		}

		return issues
	}
}

// DSL method to ensure a property exists.
// Note that this has no context, the property being checked is provided externally
// via the "property" DSL method.
func required() YamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		if !yProp.IsFound() {
			return []YamlValidationIssue{{Msg: fmt.Sprintf(`missing the "%s" required property in the %s .yaml node`,
				last(path),
				buildPathString(dropRight(path)))}}
		}

		return []YamlValidationIssue{}
	}
}

// DSL method that will only perform validations if the property exists
// Useful to avoid executing validations on none mandatory properties which are not present.
func optional(checks ...YamlCheck) YamlCheck {
	return func(y *simpleyaml.Yaml, path []string) YamlValidationIssues {
		var issues YamlValidationIssues

		// If an optional property is not found
		// no sense in executing the validations.
		if !y.IsFound() {
			return issues
		}

		for _, check := range checks {
			newIssues := check(y, path)
			issues = append(issues, newIssues...)
		}

		return issues
	}
}

func typeIsNotMapArray() YamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {

		if yProp.IsMap() || yProp.IsArray() {
			return []YamlValidationIssue{{Msg: fmt.Sprintf(`the "%s" property must be a string`, buildPathString(path))}}
		}

		return []YamlValidationIssue{}
	}
}

func typeIsArray() YamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {

		if yProp.IsFound() {
			_, err := yProp.Array()

			if err != nil {
				return []YamlValidationIssue{{Msg: fmt.Sprintf(`the "%s" property must be an array`, buildPathString(path))}}
			}
		}

		return []YamlValidationIssue{}
	}
}

func typeIsMap() YamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {

		if yProp.IsFound() {
			_, err := yProp.Map()

			if err != nil {
				return []YamlValidationIssue{{Msg: fmt.Sprintf(`the "%s" property must be a map`, buildPathString(path))}}
			}
		}

		return []YamlValidationIssue{}
	}
}

func typeIsBoolean() YamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		if yProp.IsFound() {
			_, err := yProp.Bool()

			if err != nil {
				return []YamlValidationIssue{{Msg: fmt.Sprintf(`the "%s" property must be a boolean`, buildPathString(path))}}
			}
		}

		return []YamlValidationIssue{}
	}
}

func matchesRegExp(pattern string) YamlCheck {
	regExp, _ := regexp.Compile(pattern)

	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		strValue := getLiteralStringValue(yProp)

		if !regExp.MatchString(strValue) {
			return []YamlValidationIssue{
				{Msg: fmt.Sprintf(`the "%s" value of the "%s" property does not match the "%s" pattern`, strValue, buildPathString(path), pattern)}}
		}

		return []YamlValidationIssue{}
	}
}

// Validates that value matches to one of defined enums values
func matchesEnumValues(enumValues []string) YamlCheck {
	expectedSubset := ""
	i := 0
	for _, enumValue := range enumValues {
		i++
		if i > 4 {
			break
		}
		if i > 1 {
			expectedSubset = expectedSubset + ","
		}
		expectedSubset = expectedSubset + enumValue
	}

	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		value := getLiteralStringValue(yProp)
		found := false
		for _, enumValue := range enumValues {
			if enumValue == value {
				found = true
				break
			}
		}
		if !found {
			return []YamlValidationIssue{{Msg: fmt.Sprintf(`the "%s" value of the "%s" enum property is invalid; expected one of the following: %s`,
				value, buildPathString(path), expectedSubset)}}
		}

		return []YamlValidationIssue{}
	}
}

func prettifyPath(path string) string {
	wrongIdxSyntax, _ := regexp.Compile("\\.\\[")

	return wrongIdxSyntax.ReplaceAllString(path, "[")
}

func buildPathString(path []string) string {
	if len(path) == 0 {
		return "root"
	}

	if len(path) == 1 {
		return buildPathString(append([]string{"root"}, path...))
	}
	pathStr := strings.Join(append(path), ".")

	prettyPathStr := prettifyPath(pathStr)

	return prettyPathStr
}

func last(sl []string) string {
	return sl[len(sl)-1]
}

func dropRight(sl []string) []string {
	return sl[:len(sl)-1]
}

func getLiteralStringValue(y *simpleyaml.Yaml) string {
	strVal, strErr := y.String()

	if strErr == nil {
		return strVal
	}

	boolVal, boolErr := y.Bool()
	if boolErr == nil {
		return fmt.Sprintf("%t", boolVal)
	}

	IntVal, IntErr := y.Int()
	if IntErr == nil {
		return fmt.Sprintf("%d", IntVal)
	}

	FloatVal, FloatErr := y.Float()
	if FloatErr == nil {
		return fmt.Sprintf("%g", FloatVal)
	}

	return ""
}

// runSchemaValidations - Given a YAML text and a set of validations will execute them and will return relevant issue slice
func runSchemaValidations(yaml []byte, validations ...YamlCheck) []YamlValidationIssue {
	var issues []YamlValidationIssue

	y, parseError := simpleyaml.NewYaml(yaml)
	if parseError != nil {
		issues = appendIssue(issues, parseError.Error())
		return issues
	}

	for _, validation := range validations {
		issues = append(issues, validation(y, []string{})...)
	}

	return issues
}
