package mta_validate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func assertNoSchemaIssues(errors []YamlValidationIssue) {
	Ω(len(errors)).Should(Equal(0), "Schema issues detected")
}

var _ = Describe("Schema tests Issues", func() {

	var _ = DescribeTable("Schema issues",
		func(schema string, message string) {
			_, schemaIssues := BuildValidationsFromSchemaText([]byte(schema))
			expectSingleSchemaIssue(schemaIssues, message)
		},
		Entry("Parsing", `
type: map
# bad indentation
 mapping:
   firstName:  {required: true}`, `unmarshal []byte to yaml failed: yaml: line 3: did not find expected key`),

		Entry("Mapping", `
type: map
mapping: NotAMap`, `YAML Schema Error: <mapping> node must be a map`),

		Entry("SchemaSequenceIssue", `
type: seq
sequence: NotASequence
`, `YAML Schema Error: <sequence> node must be an array`),

		Entry("Sequence One Item", `
type: seq
sequence:
- 1
- 2
`, `YAML Schema Error: <sequence> node can only have one item`),

		Entry("Required value not bool", `
type: map
mapping:
  firstName:  {required: 123}
`, `YAML Schema Error: <required> node must be a boolean but found <123>`),

		Entry("Sequence NestedTypeNotString", `
type: map
mapping:
  firstName:  {type: [1,2] }
`, `YAML Schema Error: <type> node must be a string`),

		Entry("Pattern NotString", `
type: map
mapping:
  firstName:  {pattern: [1,2] }
`, `YAML Schema Error: <pattern> node must be a string`),

		Entry("Pattern InvalidRegex", `
type: map
mapping:
  firstName:  {required: true, pattern: '/[a-zA-Z+/'}
`, "YAML Schema Error: <pattern> node not valid: error parsing regexp: missing closing ]: `[a-zA-Z+`"),

		Entry("Enum NotString", `
type: enum
enums:
  duck : 1
  dog  : 2
`, `YAML Schema Error: enums values must be listed as array`),

		Entry("Enum NoEnumsNode", `
type: enum
enumos:
  - duck
  - dog
`, `YAML Schema Error: enums values must be listed`),

		Entry("Enum ValueNotSimple", `
type: enum
enums:
  [duck, [dog, cat]]
`, `YAML Schema Error: enum values must be simple`),
	)

	var _ = DescribeTable("Valid input",
		func(schema, input string) {
			schemaValidations, schemaIssues := BuildValidationsFromSchemaText([]byte(schema))
			assertNoSchemaIssues(schemaIssues)
			validateIssues, parseErr := ValidateYaml([]byte(input), schemaValidations...)
			assertNoParsingErrors(parseErr)
			assertNoValidationErrors(validateIssues)
		},
		Entry("Required", `
type: map
mapping:
 firstName:  {required: true}
`, `
firstName: Donald
lastName: duck`),
		Entry("Enum value", `
type: enum
enums:
  - duck
  - dog
`, `duck`),
		Entry("Sequence", `
type: seq
sequence:
- type: map
  mapping:
    name: {required: true}
`, `
- name: Donald
  lastName: duck

- name: Bugs
  lastName: Bunny

`),
		Entry("Pattern", `
type: map
mapping:
   firstName:  {required: true, pattern: '/^[a-zA-Z]+$/'}
`, `
firstName: Donald
lastName: duck
`),
		Entry("Optional", `
type: map
mapping:
   firstName:  {required: false, pattern: '/^[a-zA-Z]+$/'}
`, `
lastName: duck
`),
		Entry("Type Is Bool", `
type: map
mapping:
   isHappy:  {type: bool}
`, `
firstName: Tim
isHappy: false
`),
	)

	var _ = DescribeTable("Invalid input",
		func(schema, input, message string) {
			schemaValidations, schemaIssues := BuildValidationsFromSchemaText([]byte(schema))
			assertNoSchemaIssues(schemaIssues)
			validateIssues, parseErr := ValidateYaml([]byte(input), schemaValidations...)
			assertNoParsingErrors(parseErr)
			expectSingleValidationError(validateIssues, message)
		},
		Entry("Required", `
type: map
mapping:
   age:  {required: true}
`, `
firstName: Donald
lastName: duck
`, "Missing Required Property <age> in <root>"),

		Entry("Enum", `
type: enum
enums:
   - duck
   - dog
   - cat
   - mouse
   - elephant
`, `bird`, "Enum property <root> has invalid value. Expecting one of [duck,dog,cat,mouse]"),

		Entry("Sequence", `
type: seq
sequence:
- type: map
  mapping:
    name: {required: true}
`, `
- name: Donald
  lastName: duck

- age: 80
  lastName: Bunny
`, "Missing Required Property <name> in <root[1]>"),

		Entry("Pattern", `
type: map
mapping:
   age:  {pattern: '/^[0-9]+$/'}
`, `
name: Bamba
age: NaN
`, "Property <root.age> with value: <NaN> must match pattern: <^[0-9]+$>"),

		Entry("Optional With Pattern", `
type: map
mapping:
   firstName:  {required: false, pattern: '/^[a-zA-Z]+$/'}
`, `
firstName: Donald123
lastName: duck
`, "Property <root.firstName> with value: <Donald123> must match pattern: <^[a-zA-Z]+$>"),

		Entry("Type Is Bool", `
type: map
mapping:
   isHappy:  {type: bool}
`, `
firstName: John
isHappy: 123
`, "Property <root.isHappy> must be of type <Boolean>"),
	)
})
