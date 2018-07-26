package mta_validate

import (
	"testing"
)

func assertNoSchemaIssues(errors []string, t *testing.T) {
	if len(errors) != 0 {
		t.Error("Schema issues detected")
	}
}

func TestSchemaParseIssue(t *testing.T) {
	var schema = []byte(`
type: map
# bad indentation
 mapping:
   firstName:  {required: true}
`)

	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	if len(schemaIssues) != 1 {
		t.Error("Expected a single Schema issue")
	}
}

func TestSchemaMappingIssue(t *testing.T) {
	var schema = []byte(`
type: map
mapping: NotAMap
`)
	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	expectSingleSchemaIssue(schemaIssues,
		`YAML Schema Error: <mapping> node must be a map`, t)
}

func TestSchemaSequenceIssue(t *testing.T) {
	var schema = []byte(`
type: seq
sequence: NotASequence
`)
	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	expectSingleSchemaIssue(schemaIssues,
		`YAML Schema Error: <sequence> node must be an array`, t)
}

func TestSchemaSequenceOneItemIssue(t *testing.T) {
	var schema = []byte(`
type: seq
sequence: 
- 1
- 2
`)
	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	expectSingleSchemaIssue(schemaIssues,
		`YAML Schema Error: <sequence> node can only have one item`, t)
}

func TestSchemaRequiredNotBoolIssue(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {required: 123}
`)
	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	expectSingleSchemaIssue(schemaIssues,
		`YAML Schema Error: <required> node must be a boolean but found <123>`, t)
}

func TestSchemaNestedTypeNotStringIssue(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {type: [1,2] }
`)
	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	expectSingleSchemaIssue(schemaIssues,
		`YAML Schema Error: <type> node must be a string`, t)
}

func TestSchemaPatternNotStringIssue(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {pattern: [1,2] }
`)
	_, schemaIssues := BuildValidationsFromSchemaText(schema)
	expectSingleSchemaIssue(schemaIssues,
		`YAML Schema Error: <pattern> node must be a string`, t)
}

func TestSchemaRequiredValid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {required: true}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestSchemaRequiredInvalid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   age:  {required: true}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Missing Required Property <age> in <root>`,
		t)
}

func TestSchemaSequenceValid(t *testing.T) {
	var schema = []byte(`
type: seq
sequence:
- type: map
  mapping:
    name: {required: true}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
- name: Donald
  lastName: duck

- name: Bugs
  lastName: Bunny

`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestSchemaSequenceInvalid(t *testing.T) {
	var schema = []byte(`
type: seq
sequence:
- type: map
  mapping:
    name: {required: true}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
- name: Donald
  lastName: duck

- age: 80
  lastName: Bunny
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Missing Required Property <name> in <root[1]>`,
		t)
}

func TestSchemaPatternValid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {required: true}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestSchemaPatternInvalid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   age:  {pattern: '/^[0-9]+$/'}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
name: Bamba
age: NaN
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.age> with value: <NaN> must match pattern: <^[0-9]+$>`,
		t)
}

func TestSchemaOptionalValid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {required: false, pattern: '/^[a-zA-Z]+$/'}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: Donald
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestSchemaOptionalInvalid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   firstName:  {required: false, pattern: '/^[a-zA-Z]+$/'}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: Donald123
lastName: duck
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.firstName> with value: <Donald123> must match pattern: <^[a-zA-Z]+$>`,
		t)
}

func TestSchemaTypeIsBoolValid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   isHappy:  {type: bool}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: Tim
isHappy: false
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	assertNoValidationErrors(validateIssues, t)
}

func TestSchemaTypeIsBoolInvalid(t *testing.T) {
	var schema = []byte(`
type: map
mapping:
   isHappy:  {type: bool}
`)

	schemaValidations, schemaIssues := BuildValidationsFromSchemaText(schema)
	assertNoSchemaIssues(schemaIssues, t)

	input := []byte(`
firstName: John
isHappy: 123
`)
	validateIssues, parseErr := ValidateYaml(input, schemaValidations...)
	assertNoParsingErrors(parseErr, t)
	expectSingleValidationError(validateIssues,
		`Property <root.isHappy> must be of type <Boolean>`,
		t)
}
