// TODO: Implement additional validations
// 1. Unique:
// 4. Allowed Properties.
// 5. TypeIsNotMapOrSet

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
	s := ""
	for _, issue := range issues {
		s = s + issue.Msg
	}
	return s
}

type yamlCheck func(y *simpleyaml.Yaml, path []string) YamlValidationIssues

// DSL method to execute validations on a sub node(property) of a YAML tree.
// Can be nested to check properties farther and farther down the tree.
func property(propName string, checks ...yamlCheck) yamlCheck {
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
	checks ...yamlCheck) yamlCheck {

	return sequenceInternal(false, checks...)
}

// DSL method to execute validations in order and break early as soon as the first one fails
// This is very useful if a certain validation cannot be executed without the previous ones succeeding.
// For example: matching vs a regExp should not be performed for a property that is not a string.
func sequenceFailFast(
	checks ...yamlCheck) yamlCheck {

	return sequenceInternal(true, checks...)
}

func sequenceInternal(failfast bool,
	checks ...yamlCheck) yamlCheck {

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
func forEach(checks ...yamlCheck) yamlCheck {

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
func required() yamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		if !yProp.IsFound() {
			return []YamlValidationIssue{{Msg: fmt.Sprintf("Missing required property <%s> in <%s>",
				last(path),
				buildPathString(dropRight(path)))}}
		}

		return []YamlValidationIssue{}
	}
}

// DSL method that will only perform validations if the property exists
// Useful to avoid executing validations on none mandatory properties which are not present.
func optional(checks ...yamlCheck) yamlCheck {
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

func typeIsNotMapArray() yamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {

		if yProp.IsMap() || yProp.IsArray() {
			return []YamlValidationIssue{{Msg: fmt.Sprintf("property <%s> must be of type <string>", buildPathString(path))}}
		}

		return []YamlValidationIssue{}
	}
}

func typeIsArray() yamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {

		if yProp.IsFound() {
			_, err := yProp.Array()

			if err != nil {
				return []YamlValidationIssue{{Msg: fmt.Sprintf("property <%s> must be of type <Array>", buildPathString(path))}}
			}
		}

		return []YamlValidationIssue{}
	}
}

func typeIsMap() yamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {

		if yProp.IsFound() {
			_, err := yProp.Map()

			if err != nil {
				return []YamlValidationIssue{{Msg: fmt.Sprintf("property <%s> must be of type <Map>", buildPathString(path))}}
			}
		}

		return []YamlValidationIssue{}
	}
}

func typeIsBoolean() yamlCheck {
	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		if yProp.IsFound() {
			_, err := yProp.Bool()

			if err != nil {
				return []YamlValidationIssue{{Msg: fmt.Sprintf("property <%s> must be of type <Boolean>", buildPathString(path))}}
			}
		}

		return []YamlValidationIssue{}
	}
}

func matchesRegExp(pattern string) yamlCheck {
	regExp, _ := regexp.Compile(pattern)

	return func(yProp *simpleyaml.Yaml, path []string) YamlValidationIssues {
		strValue := getLiteralStringValue(yProp)

		if !regExp.MatchString(strValue) {
			return []YamlValidationIssue{
				{Msg: fmt.Sprintf("property <%s> with value: <%s> must match pattern: <%s>", buildPathString(path), strValue, pattern)}}
		}

		return []YamlValidationIssue{}
	}
}

// Validates that value matches to one of defined enums values
func matchesEnumValues(enumValues []string) yamlCheck {
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
			return []YamlValidationIssue{{Msg: fmt.Sprintf("Enum property <%s> has invalid value. Expecting one of [%s]",
				buildPathString(path), expectedSubset)}}
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

// ValidateYaml - Given a YAML text and a set of validations will execute them and will return relevant issue slice
// And an "err" object in case of a parsing error.
func ValidateYaml(yaml []byte, validations ...yamlCheck) ([]YamlValidationIssue, error) {
	var issues []YamlValidationIssue

	y, parseError := simpleyaml.NewYaml(yaml)
	if parseError != nil {
		return issues, parseError
	}

	for _, validation := range validations {
		issues = append(issues, validation(y, []string{})...)
	}

	return issues, nil
}
