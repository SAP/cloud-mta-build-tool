// By naming this file with _test suffix it is not measured
// in the coverage report, although we do end-up with a strange file name...
package validate

import (
	. "github.com/onsi/gomega"
)

func assertNoParsingErrors(err error) {
	Ω(err).Should(BeNil(), "Yaml Parsing Errors Detected: %v")
}

func assertNoValidationErrors(errors []YamlValidationIssue) {
	Ω(len(errors)).Should(Equal(0), "Validation issues detected: %v")
}

func expectSingleValidationError(actual []YamlValidationIssue, expectedMsg string) {
	numOfErrors := len(actual)
	Ω(numOfErrors).Should(Equal(1), "A single validation issue expected but found: <%d>", numOfErrors)

	actualMsg := actual[0].Msg
	Ω(actual[0].Msg).Should(Equal(expectedMsg), "expecting <%s>.\n\t but found <%s>.", expectedMsg, actualMsg)
}

func expectMultipleValidationError(actualIssues []YamlValidationIssue, expectedMsgs []string) {
	expectedNumOfErrors := len(expectedMsgs)
	actualNumOfErrors := len(actualIssues)
	Ω(actualNumOfErrors).Should(Equal(expectedNumOfErrors), "wrong number of issues found expected <%d> but found: <%d>", expectedNumOfErrors, actualNumOfErrors)

	for _, issue := range actualIssues {
		Ω(expectedMsgs).Should(ContainElement(issue.Msg))
	}
}

func expectSingleSchemaIssue(actual []YamlValidationIssue, expectedMsg string) {
	numOfErrors := len(actual)
	Ω(numOfErrors).Should(Equal(1), "a single validation issue expected but found: <%d>", numOfErrors)

	actualMsg := actual[0]
	Ω(actual[0].Msg).Should(Equal(expectedMsg), "expecting <%s>.\n\t but found <%s>.", expectedMsg, actualMsg)
}
