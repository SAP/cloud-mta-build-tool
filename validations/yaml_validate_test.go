package validate

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
	"github.com/smallfish/simpleyaml"
)

var _ = Describe("Yaml Validation", func() {

	DescribeTable("Valid Yaml", func(data string, validations ...YamlCheck) {
		validateIssues, parseErr := Yaml([]byte(data), validations...)

		assertNoParsingErrors(parseErr)
		assertNoValidationErrors(validateIssues)
	},
		Entry("matchesRegExp", `
firstName: Donald
lastName: duck
`, property("lastName", matchesRegExp("^[A-Za-z0-9_\\-\\.]+$"))),

		Entry("required", `
firstName: Donald
lastName: duck
`, property("firstName", required())),

		Entry("Type Is String", `
firstName: Donald
lastName: duck
`, property("firstName", typeIsNotMapArray())),

		Entry("Type Is Bool", `
name: bisli
registered: false
`, property("registered", typeIsBoolean())),

		Entry("Type Is Array", `
firstName:
   - 1
   - 2
   - 3
lastName: duck
`, property("firstName", typeIsArray())),

		Entry("sequenceFailFast", `
firstName: Hello
lastName: World
`, property("firstName", sequence(required(), matchesRegExp("^[A-Za-z0-9]+$")))),

		Entry("Type Is Map", `
firstName:
   - 1
   - 2
   - 3
lastName:
   a : 1
   b : 2
`, property("lastName", typeIsMap())),

		Entry("For Each", `
firstName: Hello
lastName: World
classes:
 - name: biology
   room: MR113

 - name: history
   room: MR225

`, property("classes", sequence(
			required(),
			typeIsArray(),
			forEach(
				property("name", required()),
				property("room", matchesRegExp("^MR[0-9]+$")))))),

		Entry("optional Exists", `
firstName: Donald
lastName: duck
`, property("firstName", optional(typeIsNotMapArray()))),

		Entry("optional Missing", `
lastName: duck
`, property("firstName", optional(typeIsNotMapArray()))),
	)

	DescribeTable("Invalid Yaml", func(data, message string, validations ...YamlCheck) {
		validateIssues, parseErr := Yaml([]byte(data), validations...)

		assertNoParsingErrors(parseErr)
		expectSingleValidationError(validateIssues, message)
	},
		Entry("matchesRegExp", `
firstName: Donald
lastName: duck
`, `the root.firstName property with the Donald value does not match the ^[0-9_\-\.]+$ pattern`,
			property("firstName", matchesRegExp("^[0-9_\\-\\.]+$"))),

		Entry("required", `
firstName: Donald
lastName: duck
`, `missing the age required property in root`,
			property("age", required())),

		Entry("required", `
firstName:
   - 1
   - 2
   - 3
lastName: duck
`, `the root.firstName property must be of the string type`,
			property("firstName", typeIsNotMapArray())),

		Entry("TypeIsBool", `
name: bamba
registered: 123
`, `the root.registered property must be of the boolean type`,
			property("registered", typeIsBoolean())),

		Entry("typeIsArray", `
firstName:
   - 1
   - 2
   - 3
lastName: duck
`, `the root.lastName property must be of the array type`,
			property("lastName", typeIsArray())),

		Entry("typeIsMap", `
firstName:
   - 1
   - 2
   - 3
lastName:
   a : 1
   b : 2
`, `the root.firstName property must be of the map type`,
			property("firstName", typeIsMap())),

		Entry("sequenceFailFast", `
firstName: Hello
lastName: World
`, `missing the missing required property in root`,
			property("missing", sequenceFailFast(
				required(),
				// This second validation should not be executed as sequence breaks early.
				matchesRegExp("^[0-9]+$")))),

		Entry("OptionalExists", `
firstName:
  - 1
  - 2
lastName: duck
`, `the root.firstName property must be of the string type`,
			property("firstName", optional(typeIsNotMapArray()))),
	)

	It("InvalidYamlHandling", func() {
		data := []byte(`
firstName: Donald
  lastName: duck # invalid indentation
		`)
		_, parseErr := Yaml(data, property("lastName", required()))
		Ω(parseErr).Should(HaveOccurred())
	})

	It("ForEachInValid", func() {
		data := []byte(`
firstName: Hello
lastName: World
classes:
 - name: biology
   room: oops

 - room: 225

`)
		validations := property("classes", sequence(
			required(),
			typeIsArray(),
			forEach(
				property("name", required()),
				property("room", matchesRegExp("^[0-9]+$")))))

		validateIssues, parseErr := Yaml(data, validations)

		assertNoParsingErrors(parseErr)
		expectMultipleValidationError(validateIssues,
			[]string{
				"the classes[0].room property with the oops value does not match the ^[0-9]+$ pattern",
				"missing the name required property in classes[1]"})
	})
})

var _ = DescribeTable("GetLiteralStringValue", func(data string, matcher GomegaMatcher) {
	y, _ := simpleyaml.NewYaml([]byte(data))
	value := getLiteralStringValue(y)
	Ω(value).Should(matcher)
},
	Entry("Invalid", `
  [a,b]
`, BeEmpty()),
	Entry("Valid", fmt.Sprintf("%g", 0.55), Equal("0.55")),
)
