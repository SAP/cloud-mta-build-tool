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

	var _ = DescribeTable("Validation", func(locationSource, mtaFilename string, issuesNumber int, validateProject bool) {
		yamlContent, _ := ioutil.ReadFile(filepath.Join(locationSource, mtaFilename))
		issues, _ := validate(yamlContent, locationSource, true, validateProject)
		Ω(len(issues)).Should(Equal(issuesNumber))
	},

		Entry("Validate All", getTestPath("testproject"), "mta.yaml", 1, true),
		Entry("Validate Schema", getTestPath(), "mta_multiapps.yaml", 0, false),
	)

	var _ = Describe("Parsing", func() {
		It("Modules parsing - sanity", func() {
			var moduleSrv = mta.Module{
				Name: "srv",
				Type: "java",
				Path: "srv",
				Requires: []mta.Requires{
					{
						Name: "db",
						Properties: mta.Properties{
							"JBP_CONFIG_RESOURCE_CONFIGURATION": `[tomcat/webapps/ROOT/META-INF/context.xml: {"service_name_for_DefaultDB" : "~{hdi-container-name}"}]`,
						},
					},
				},
				Provides: []mta.Provides{
					{
						Name:       "srv_api",
						Properties: mta.Properties{"url": "${default-url}"},
					},
				},
				Parameters: mta.Parameters{"memory": "512M"},
				Properties: mta.Properties{
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
						Properties: mta.Properties{
							"forwardAuthToken": true,
							"strictSSL":        false,
							"name":             "srv_api",
							"url":              "~{url}",
						},
					},
				},
				BuildParams: mta.BuildParameters{Builder: "grunt"},
				Parameters:  mta.Parameters{"disk-quota": "256M", "memory": "256M"},
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

	var _ = Describe("Get methods on MTA", func() {
		modules := []*mta.Module{
			{
				Name: "someproj-db",
				Type: "hdb",
				Path: "db",
				Requires: []mta.Requires{
					{
						Name: "someproj-hdi-container",
					},
					{
						Name: "someproj-logging",
					},
				},
			},
			{
				Name: "someproj-java",
				Type: "java",
				Path: "srv",
				Parameters: mta.Parameters{
					"memory":     "512M",
					"disk-quota": "256M",
				},
			},
		}
		schemaVersion := "0.0.2"
		mta := &mta.MTA{
			SchemaVersion: &schemaVersion,
			ID:            "MTA",
			Version:       "1.1.1",
			Modules:       modules,
			Resources: []*mta.Resource{
				{
					Name: "someproj-hdi-container",
					Properties: mta.Properties{
						"hdi-container-name": "${service-name}",
					},
					Type: "container",
				},
				{
					Name: "someproj-apprepo-rt",
					Type: "org.cloudfoundry.managed-service",
					Parameters: mta.Parameters{
						"service":      "html5-apps-repo",
						"service-plan": "app-runtime",
					},
				},
			}}
		It("GetModules", func() {
			Ω(mta.GetModules()).Should(Equal(modules))
		})
		It("GetResourceByName - Sanity", func() {
			Ω(mta.GetResourceByName("someproj-hdi-container")).Should(Equal(mta.Resources[0]))
		})
		It("GetResourceByName - Negative", func() {
			_, err := mta.GetResourceByName("")
			Ω(err).Should(HaveOccurred())
		})
		It("GetResources - Sanity ", func() {
			Ω(mta.GetResources()).Should(Equal(mta.Resources))
		})
		It("GetModuleByName - Sanity ", func() {
			Ω(mta.GetModuleByName("someproj-db")).Should(Equal(modules[0]))
		})
		It("GetModuleByName - Negative ", func() {
			_, err := mta.GetModuleByName("foo")
			Ω(err).Should(HaveOccurred())
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
