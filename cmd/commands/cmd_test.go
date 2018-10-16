package commands

import (
	"testing"
)

func Test_getValidationMode(t *testing.T) {
	validModeTest := []struct{
		n string //input
		expected [2] bool //expected result
		isNilErr bool //expected error
	}{
		{"", [2]bool{true, true}, true},
		{"schema", [2]bool{true, false}, true},
		{"project", [2]bool{false, true}, true},
		{"value", [2]bool{false, false}, false},
	}

	for _, elem := range validModeTest {
		res1, res2, err := getValidationMode(elem.n)
		isNilErr := err == nil
		if res1 != elem.expected[0] || res2 != elem.expected[1] || isNilErr != elem.isNilErr {
			t.Errorf("expected output (%v,%v), actual (%v,%v)", elem.expected[0], elem.expected[1], res1, res2)
		}
	}
}

func Test_validateMtaYaml(t *testing.T) {
	validateMtaTest := []struct {
		yamlPath string //input
		yamlFileName string //input
		validateSchema bool //input
		validateProject bool //input
		isNilErr bool //expected error
	}{
		{"ui5app", "mta.yaml", true, true, false},
		{"ui5app", "mta.yaml", true, false, false},
		{"ui5app", "mta.yaml", false, true, false},
		{"ui5app", "mta.yaml", false, false, true},
	}

	for _, elem := range validateMtaTest {
		err := validateMtaYaml(elem.yamlPath, elem.yamlFileName, elem.validateSchema, elem.validateProject)
		isNilErr := err == nil
		if isNilErr != elem.isNilErr {
			t.Errorf("expected error isNill %v, actual isNill %v", elem.isNilErr, isNilErr)
		}
	}
}
