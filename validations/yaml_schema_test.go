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
mapping: NotAMap`, `invalid .yaml file schema: the mapping node must be a map`),

		Entry("SchemaSequenceIssue", `
type: seq
sequence: NotASequence
`, `invalid .yaml file schema: the sequence node must be an array`),

		Entry("sequence One Item", `
type: seq
sequence:
- 1
- 2
`, `invalid .yaml file schema: the sequence node can have only one item`),

		Entry("required value not bool", `
type: map
mapping:
  firstName:  {required: 123}
`, `invalid .yaml file schema: the required node must be a boolean but found the 123 value`),

		Entry("sequence NestedTypeNotString", `
type: map
mapping:
  firstName:  {type: [1,2] }
`, `invalid .yaml file schema: the type node must be a string`),

		Entry("Pattern NotString", `
type: map
mapping:
  firstName:  {pattern: [1,2] }
`, `invalid .yaml file schema: the pattern node must be a string`),

		Entry("Pattern InvalidRegex", `
type: map
mapping:
  firstName:  {required: true, pattern: '/[a-zA-Z+/'}
`, "invalid .yaml file schema: the pattern node is invalid because: error parsing regexp: missing closing ]: `[a-zA-Z+`"),

		Entry("Enum NotString", `
type: enum
enums:
  duck : 1
  dog  : 2
`, `invalid .yaml file schema: enums values must be listed as an array`),

		Entry("Enum NoEnumsNode", `
type: enum
enumos:
  - duck
  - dog
`, `invalid .yaml file schema: enums values must be listed`),

		Entry("Enum ValueNotSimple", `
type: enum
enums:
  [duck, [dog, cat]]
`, `invalid .yaml file schema: enum values must be simple`),
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
`, "missing the age required property in the root .yaml node"),

		Entry("Enum", `
type: enum
enums:
   - duck
   - dog
   - cat
   - mouse
   - elephant
`, `bird`, "the bird value of the root enum property is invalid; expected one of the following: duck,dog,cat,mouse"),

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
`, "missing the name required property in the root[1] .yaml node"),

		Entry("Pattern", `
type: map
mapping:
   age:  {pattern: '/^[0-9]+$/'}
`, `
name: Bamba
age: NaN
`, "the NaN value of the root.age property does not match the ^[0-9]+$ pattern"),

		Entry("optional With Pattern", `
type: map
mapping:
   firstName:  {required: false, pattern: '/^[a-zA-Z]+$/'}
`, `
firstName: Donald123
lastName: duck
`, "the Donald123 value of the root.firstName property does not match the ^[a-zA-Z]+$ pattern"),

		Entry("Type Is Bool", `
type: map
mapping:
   isHappy:  {type: bool}
`, `
firstName: John
isHappy: 123
`, "the root.isHappy property must be a boolean"),
	)
})
