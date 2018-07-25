// By naming this file with _test suffix it is not measured
// in the coverage report, although we do end-up with a strange file name...
package mta_validate

import "testing"
import "github.com/stretchr/testify/assert"

func assertNoParsingErrors(err error, t *testing.T) {
	if err != nil {
		t.Error("Yaml Parsing Errors Detected")
	}
}

func assertNoValidationErrors(errors []YamlValidationIssue, t *testing.T) {
	if len(errors) != 0 {
		t.Error("Validation issues detected")
	}
}

func expectSingleValidationError(actual []YamlValidationIssue, expectedMsg string, t *testing.T) {
	numOfErrors := len(actual)
	if numOfErrors != 1 {
		t.Errorf("A single validation issue expected but found: <%d>", numOfErrors)
	}

	actualMsg := actual[0].msg
	if actual[0].msg != expectedMsg {
		t.Errorf("Expecting <%s>.\n\t But found <%s>.", expectedMsg, actualMsg)
	}
}

func expectMultipleValidationError(actualIssues []YamlValidationIssue, expectedMsgs []string, t *testing.T) {
	expectedNumOfErrors := len(expectedMsgs)
	actualNumOfErrors := len(actualIssues)
	if actualNumOfErrors != len(expectedMsgs) {
		t.Errorf("Wrong number of issues found expected <%d> but found: <%d>", expectedNumOfErrors, actualNumOfErrors)
	}

	var actualMsgs []string
	for _, issue := range actualIssues {
		actualMsgs = append(actualMsgs, issue.msg)
	}

	assert.Subset(t, actualMsgs, expectedMsgs)
	assert.Subset(t, expectedMsgs, actualMsgs)
}

func expectSingleSchemaIssue(actual []string, expectedMsg string, t *testing.T) {
	numOfErrors := len(actual)
	if numOfErrors != 1 {
		t.Errorf("A single validation issue expected but found: <%d>", numOfErrors)
	}

	actualMsg := actual[0]
	if actual[0] != expectedMsg {
		t.Errorf("Expecting <%s>.\n\t But found <%s>.", expectedMsg, actualMsg)
	}
}
