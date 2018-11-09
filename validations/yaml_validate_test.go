package mta_validate

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
		validateIssues, parseErr := ValidateYaml([]byte(data), validations...)

		assertNoParsingErrors(parseErr)
		assertNoValidationErrors(validateIssues)
	},
		Entry("MatchesRegExp", `
firstName: Donald
lastName: duck
`, Property("lastName", MatchesRegExp("^[A-Za-z0-9_\\-\\.]+$"))),

		Entry("Required", `
firstName: Donald
lastName: duck
`, Property("firstName", Required())),

		Entry("Type Is String", `
firstName: Donald
lastName: duck
`, Property("firstName", typeIsNotMapArray())),

		Entry("Type Is Bool", `
name: bisli
registered: false
`, Property("registered", typeIsBoolean())),

		Entry("Type Is Array", `
firstName:
   - 1
   - 2
   - 3
lastName: duck
`, Property("firstName", TypeIsArray())),

		Entry("sequenceFailFast", `
firstName: Hello
lastName: World
`, Property("firstName", Sequence(Required(), MatchesRegExp("^[A-Za-z0-9]+$")))),

		Entry("Type Is Map", `
firstName:
   - 1
   - 2
   - 3
lastName:
   a : 1
   b : 2
`, Property("lastName", TypeIsMap())),

		Entry("For Each", `
firstName: Hello
lastName: World
classes:
 - name: biology
   room: MR113

 - name: history
   room: MR225

`, Property("classes", Sequence(
			Required(),
			TypeIsArray(),
			ForEach(
				Property("name", Required()),
				Property("room", MatchesRegExp("^MR[0-9]+$")))))),

		Entry("Optional Exists", `
firstName: Donald
lastName: duck
`, Property("firstName", Optional(typeIsNotMapArray()))),

		Entry("Optional Missing", `
lastName: duck
`, Property("firstName", Optional(typeIsNotMapArray()))),
	)

	DescribeTable("Invalid Yaml", func(data, message string, validations ...YamlCheck) {
		validateIssues, parseErr := ValidateYaml([]byte(data), validations...)

		assertNoParsingErrors(parseErr)
		expectSingleValidationError(validateIssues, message)
	},
		Entry("MatchesRegExp", `
firstName: Donald
lastName: duck
`, `Property <root.firstName> with value: <Donald> must match pattern: <^[0-9_\-\.]+$>`,
			Property("firstName", MatchesRegExp("^[0-9_\\-\\.]+$"))),

		Entry("Required", `
firstName: Donald
lastName: duck
`, `Missing Required Property <age> in <root>`,
			Property("age", Required())),

		Entry("Required", `
firstName:
   - 1
   - 2
   - 3
lastName: duck
`, `Property <root.firstName> must be of type <string>`,
			Property("firstName", typeIsNotMapArray())),

		Entry("TypeIsBool", `
name: bamba
registered: 123
`, `Property <root.registered> must be of type <Boolean>`,
			Property("registered", typeIsBoolean())),

		Entry("TypeIsArray", `
firstName:
   - 1
   - 2
   - 3
lastName: duck
`, `Property <root.lastName> must be of type <Array>`,
			Property("lastName", TypeIsArray())),

		Entry("TypeIsMap", `
firstName:
   - 1
   - 2
   - 3
lastName:
   a : 1
   b : 2
`, `Property <root.firstName> must be of type <Map>`,
			Property("firstName", TypeIsMap())),

		Entry("sequenceFailFast", `
firstName: Hello
lastName: World
`, `Missing Required Property <missing> in <root>`,
			Property("missing", sequenceFailFast(
				Required(),
				// This second validation should not be executed as sequence breaks early.
				MatchesRegExp("^[0-9]+$")))),

		Entry("OptionalExists", `
firstName:
  - 1
  - 2
lastName: duck
`, `Property <root.firstName> must be of type <string>`,
			Property("firstName", Optional(typeIsNotMapArray()))),
	)

	It("InvalidYamlHandling", func() {
		data := []byte(`
firstName: Donald
  lastName: duck # invalid indentation
		`)
		_, parseErr := ValidateYaml(data, Property("lastName", Required()))
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
		validations := Property("classes", Sequence(
			Required(),
			TypeIsArray(),
			ForEach(
				Property("name", Required()),
				Property("room", MatchesRegExp("^[0-9]+$")))))

		validateIssues, parseErr := ValidateYaml(data, validations)

		assertNoParsingErrors(parseErr)
		expectMultipleValidationError(validateIssues,
			[]string{
				"Property <classes[0].room> with value: <oops> must match pattern: <^[0-9]+$>",
				"Missing Required Property <name> in <classes[1]>"})
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
