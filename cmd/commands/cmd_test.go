package commands

import "testing"

func Test_getValidationMode(t *testing.T) {
	validModeTest := []struct{
		n string //input
		expected [2] bool //expected result
	}{
		{"", [2]bool{true, true}},
		{"schema", [2]bool{true, false}},
		{"project", [2]bool{false, true}},
	}

	for _, elem := range validModeTest {
		first, second:= getValidationMode(elem.n)
		if first != elem.expected[0] || second != elem.expected[1] {
			t.Errorf("expected output (%v,%v), actual (%v,%v)", elem.expected[0], elem.expected[1], first, second)
		}
	}
}

/*func Test_validateMtaYaml(t *testing.T) {
	validateMtaTest := []struct {
		yamlPath string //input
		yamlFileName string //input
		validateSchema bool //input
		validateProject bool //input
		expected error //expected result
	}{
		{"ui5app", "mta.yaml", true, true, nil},
		{"ui5app", "mta.yaml", true, false, error()},
		{"ui5app", "mta.yaml", false, true, error()},
	}

	for _, elem := range validateMtaTest {
		err := getValidationMode(elem.yamlPath, elem.yamlFileName, elem.validateSchema, elem.validateProject)
		if err != elem.expected {
			t.Errorf("expected output %v, actual %v", elem.expected, err)
		}
	}
}
*/