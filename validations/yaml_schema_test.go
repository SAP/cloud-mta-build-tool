package validate

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
mapping: NotAMap`, `invalid yaml schema: <mapping> node must be a map`),

		Entry("SchemaSequenceIssue", `
type: seq
sequence: NotASequence
`, `invalid yaml schema: <sequence> node must be an array`),

		Entry("sequence One Item", `
type: seq
sequence:
- 1
- 2
`, `invalid yaml schema: <sequence> node can only have one item`),

		Entry("required value not bool", `
type: map
mapping:
  firstName:  {required: 123}
`, `invalid yaml schema: <required> node must be a boolean but found <123>`),

		Entry("sequence NestedTypeNotString", `
type: map
mapping:
  firstName:  {type: [1,2] }
`, `invalid yaml schema: <type> node must be a string`),

		Entry("Pattern NotString", `
type: map
mapping:
  firstName:  {pattern: [1,2] }
`, `invalid yaml schema: <pattern> node must be a string`),

		Entry("Pattern InvalidRegex", `
type: map
mapping:
  firstName:  {required: true, pattern: '/[a-zA-Z+/'}
`, "invalid yaml schema: <pattern> node not valid: error parsing regexp: missing closing ]: `[a-zA-Z+`"),

		Entry("Enum NotString", `
type: enum
enums:
  duck : 1
  dog  : 2
`, `invalid yaml schema: enums values must be listed as array`),

		Entry("Enum NoEnumsNode", `
type: enum
enumos:
  - duck
  - dog
`, `invalid yaml schema: enums values must be listed`),

		Entry("Enum ValueNotSimple", `
type: enum
enums:
  [duck, [dog, cat]]
`, `invalid yaml schema: enum values must be simple`),
	)

	var _ = DescribeTable("Valid input",
		func(schema, input string) {
			schemaValidations, schemaIssues := BuildValidationsFromSchemaText([]byte(schema))
			assertNoSchemaIssues(schemaIssues)
			validateIssues, parseErr := Yaml([]byte(input), schemaValidations...)
			assertNoParsingErrors(parseErr)
			assertNoValidationErrors(validateIssues)
		},
		Entry("required", `
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
		Entry("sequence", `
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
		Entry("optional", `
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
			validateIssues, parseErr := Yaml([]byte(input), schemaValidations...)
			assertNoParsingErrors(parseErr)
			expectSingleValidationError(validateIssues, message)
		},
		Entry("required", `
type: map
mapping:
   age:  {required: true}
`, `
firstName: Donald
lastName: duck
`, "missing required property <age> in <root>"),

		Entry("Enum", `
type: enum
enums:
   - duck
   - dog
   - cat
   - mouse
   - elephant
`, `bird`, "enum property <root> has invalid value. expecting one of [duck,dog,cat,mouse]"),

		Entry("sequence", `
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
`, "missing required property <name> in <root[1]>"),

		Entry("Pattern", `
type: map
mapping:
   age:  {pattern: '/^[0-9]+$/'}
`, `
name: Bamba
age: NaN
`, "property <root.age> with value: <NaN> must match pattern: <^[0-9]+$>"),

		Entry("optional With Pattern", `
type: map
mapping:
   firstName:  {required: false, pattern: '/^[a-zA-Z]+$/'}
`, `
firstName: Donald123
lastName: duck
`, "property <root.firstName> with value: <Donald123> must match pattern: <^[a-zA-Z]+$>"),

		Entry("Type Is Bool", `
type: map
mapping:
   isHappy:  {type: bool}
`, `
firstName: John
isHappy: 123
`, "property <root.isHappy> must be of type <Boolean>"),
	)
})
