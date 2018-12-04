package validate

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/mta"
)

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}

var _ = Describe("MTA tests", func() {

	var _ = Describe("Parsing", func() {
		It("Modules parsing - sanity", func() {
			var moduleSrv = mta.Module{
				Name: "srv",
				Type: "java",
				Path: "srv",
				Requires: []mta.Requires{
					{
						Name: "db",
						Properties: map[string]interface{}{
							"JBP_CONFIG_RESOURCE_CONFIGURATION": `[tomcat/webapps/ROOT/META-INF/context.xml: {"service_name_for_DefaultDB" : "~{hdi-container-name}"}]`,
						},
					},
				},
				Provides: []mta.Provides{
					{
						Name:       "srv_api",
						Properties: map[string]interface{}{"url": "${default-url}"},
					},
				},
				Parameters: map[string]interface{}{"memory": "512M"},
				Properties: map[string]interface{}{
					"VSCODE_JAVA_DEBUG_LOG_LEVEL": "ALL",
					"APPC_LOG_LEVEL":              "info",
				},
			}
			var moduleUI = mta.Module{
				Name: "ui",
				Type: "html5",
				Path: "ui",
				Requires: []mta.Requires{
					{
						Name:  "srv_api",
						Group: "destinations",
						Properties: map[string]interface{}{
							"forwardAuthToken": true,
							"strictSSL":        false,
							"name":             "srv_api",
							"url":              "~{url}",
						},
					},
				},
				BuildParams: map[string]interface{}{"builder": "grunt"},
				Parameters:  map[string]interface{}{"disk-quota": "256M", "memory": "256M"},
			}
			var modules = []*mta.Module{&moduleSrv, &moduleUI}
			mtaFile, _ := ioutil.ReadFile("./testdata/mta.yaml")
			// Unmarshal file
			oMta := &mta.MTA{}
			Ω(yaml.Unmarshal(mtaFile, oMta)).Should(Succeed())
			Ω(oMta.Modules).Should(HaveLen(2))
			Ω(oMta.GetModules()).Should(Equal(modules))

		})

	})

	var _ = Describe("Validation", func() {
		var _ = DescribeTable("getValidationMode", func(flag string, expectedValidateSchema, expectedValidateProject, expectedSuccess bool) {
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

		var _ = DescribeTable("validateMtaYaml", func(projectRelPath string, validateSchema, validateProject, expectedSuccess bool) {
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

})
