package validate

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}

var _ = Describe("Validation", func() {
	var _ = DescribeTable("GetValidationMode", func(flag string, expectedValidateSchema, expectedValidateProject, expectedSuccess bool) {
		res1, res2, err := GetValidationMode(flag)
		Ω(res1).Should(Equal(expectedValidateSchema))
		Ω(res2).Should(Equal(expectedValidateProject))
		Ω(err == nil).Should(Equal(expectedSuccess))
	},
		Entry("all", "", true, true, true),
		Entry("schema", "schema", true, false, true),
		Entry("project", "project", false, true, true),
		Entry("invalid", "value", false, false, false),
	)

	var _ = DescribeTable("ValidateMtaYaml", func(projectRelPath string, validateSchema, validateProject, expectedSuccess bool) {
		err := ValidateMtaYaml(getTestPath(projectRelPath), "mta.yaml", validateSchema, validateProject)
		Ω(err == nil).Should(Equal(expectedSuccess))
	},
		Entry("invalid path to yaml - all", "ui5app1", true, true, false),
		Entry("invalid path to yaml - schema", "ui5app1", true, false, false),
		Entry("invalid path to yaml - project", "ui5app1", false, true, false),
		Entry("invalid path to yaml - nothing to validate", "ui5app1", false, false, true),
		Entry("valid schema", "mtahtml5", true, false, true),
		Entry("invalid project - no ui5app2 path", "mtahtml5", false, true, false),
	)

})
