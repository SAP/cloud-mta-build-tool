package mta_validate

import (
	"fmt"
	"testing"

	"github.com/smallfish/simpleyaml"
)
import "github.com/stretchr/testify/assert"

func TestInvalidYamlHandling(t *testing.T) {
	data := []byte(`
firstName: Donald
  lastName: duck # invalid indentation
`)
	_, parseErr := ValidateYaml(data, Property("lastName",
		Required()))

	assert.NotNil(t, parseErr)
}

func TestMatchesRegExpValid(t *testing.T) {
	data := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("lastName",
		MatchesRegExp("^[A-Za-z0-9_\\-\\.]+$")))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestMatchesRegExpInvalid(t *testing.T) {
	data := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data,
		Property("firstName",
			MatchesRegExp("^[0-9_\\-\\.]+$")))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.firstName> with value: <Donald> must match pattern: <^[0-9_\-\.]+$>`,
		t)
}

func TestRequiredValid(t *testing.T) {
	data := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("firstName",
		Required()))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestRequiredInvalid(t *testing.T) {
	data := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("age",
		Required()))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Missing Required Property <age> in <root>`,
		t)
}

func TestTypeIsStringValid(t *testing.T) {
	data := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("firstName",
		TypeIsNotMapArray()))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestTypeIsStringInvalid(t *testing.T) {
	data := []byte(`
firstName: 
   - 1
   - 2
   - 3
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("firstName",
		TypeIsNotMapArray()))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.firstName> must be of type <string>`,
		t)
}

func TestTypeIsBoolValid(t *testing.T) {
	data := []byte(`
name: bisli
registered: false
`)
	validateIssues, parseErr := ValidateYaml(data, Property("registered",
		TypeIsBoolean()))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestTypeIsBoolInvalid(t *testing.T) {
	data := []byte(`
name: bamba
registered: 123
`)
	validateIssues, parseErr := ValidateYaml(data, Property("registered",
		TypeIsBoolean()))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.registered> must be of type <Boolean>`,
		t)
}

func TestTypeIsArrayValid(t *testing.T) {
	data := []byte(`
firstName: 
   - 1
   - 2
   - 3
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("firstName",
		TypeIsArray()))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestTypeIsArrayInvalid(t *testing.T) {
	data := []byte(`
firstName: 
   - 1
   - 2
   - 3
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data, Property("lastName",
		TypeIsArray()))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.lastName> must be of type <Array>`,
		t)
}

func TestTypeIsMapValid(t *testing.T) {
	data := []byte(`
firstName: 
   - 1
   - 2
   - 3
lastName: 
   a : 1
   b : 2
`)
	validateIssues, parseErr := ValidateYaml(data, Property("lastName",
		TypeIsMap()))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestTypeIsMapInvalid(t *testing.T) {
	data := []byte(`
firstName: 
   - 1
   - 2
   - 3
lastName: 
   a : 1
   b : 2
`)
	validateIssues, parseErr := ValidateYaml(data, Property("firstName",
		TypeIsMap()))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.firstName> must be of type <Map>`,
		t)
}

func TestSequenceFailFastValid(t *testing.T) {
	data := []byte(`
firstName: Hello
lastName: World
`)
	sequence := Property("firstName",
		Sequence(
			Required(),
			MatchesRegExp("^[A-Za-z0-9]+$")))
	validateIssues, parseErr := ValidateYaml(data, sequence)

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestSequenceFailFastInValid(t *testing.T) {
	data := []byte(`
firstName: Hello
lastName: World
`)
	sequence := Property("missing",
		SequenceFailFast(
			Required(),
			// This second validation should not be executed as sequence breaks early.
			MatchesRegExp("^[0-9]+$")))
	validateIssues, parseErr := ValidateYaml(data, sequence)

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Missing Required Property <missing> in <root>`,
		t)
}

func TestDorEachValid(t *testing.T) {
	data := []byte(`
firstName: Hello
lastName: World
classes:
 - name: biology  
   room: MR113

 - name: history
   room: MR225

`)
	validations := Property("classes", Sequence(
		Required(),
		TypeIsArray(),
		ForEach(
			Property("name",
				Required()),
			Property("room",
				MatchesRegExp("^MR[0-9]+$")))))

	validateIssues, parseErr := ValidateYaml(data, validations)

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)

}

func TestDorEachInValid(t *testing.T) {
	data := []byte(`
firstName: Hello
lastName: World
classes:
 - name: biology  
   room: oops

 - room: 225

`)
	validations := Property("classes", Sequence(
		Required(),
		TypeIsArray(),
		ForEach(
			Property("name",
				Required()),
			Property("room",
				MatchesRegExp("^[0-9]+$")))))

	validateIssues, parseErr := ValidateYaml(data, validations)

	assertNoParsingErrors(parseErr, t)
	expectMultipleValidationError(validateIssues,
		[]string{
			"Property <classes[0].room> with value: <oops> must match pattern: <^[0-9]+$>",
			"Missing Required Property <name> in <classes[1]>"},
		t)
}

func TestOptionalExistsValid(t *testing.T) {
	data := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data,
		Property("firstName",
			Optional(
				TypeIsNotMapArray())))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestOptionalMissingValid(t *testing.T) {
	data := []byte(`
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data,
		Property("firstName",
			Optional(
				TypeIsNotMapArray())))

	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestOptionalExistsInValid(t *testing.T) {
	data := []byte(`
firstName: 
  - 1
  - 2
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(data,
		Property("firstName",
			Optional(
				TypeIsNotMapArray())))

	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.firstName> must be of type <string>`,
		t)
}

func TestGetLiteralStringValueInvalid(t *testing.T) {
	data := []byte(`
  [a,b]
`)
	y, _ := simpleyaml.NewYaml(data)
	value := getLiteralStringValue(y)
	assert.Empty(t, value)

}

func TestGetLiteralStringValueFloat(t *testing.T) {
	str := fmt.Sprintf("%g", 0.55)
	data := []byte(str)
	y, _ := simpleyaml.NewYaml(data)
	value := getLiteralStringValue(y)
	assert.Equal(t, str, value)

}
