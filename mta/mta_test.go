package mta

import (
	"io/ioutil"
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

var _ = Describe("MTA tests", func() {

	var _ = DescribeTable("Validation", func(locationSource, mtaFilename string, issuesNumber int, validateProject bool) {
		ep := Loc{SourcePath: locationSource, MtaFilename: mtaFilename}
		yamlContent, _ := ReadMtaContent(&ep)
		source, _ := ep.GetSource()
		issues, _ := Validate(yamlContent, source, true, validateProject)
		Ω(len(issues)).Should(Equal(issuesNumber))
	},

		Entry("Validate All", getTestPath("testproject"), "mta.yaml", 1, true),
		Entry("Validate Schema", getTestPath(), "mta_multiapps.yaml", 0, false),
	)

	var _ = Describe("ReadMtaYaml", func() {
		It("Sanity", func() {
			res, resErr := ReadMtaYaml(&Loc{SourcePath: getTestPath("testproject")})
			Ω(res).ShouldNot(BeNil())
			Ω(resErr).Should(BeNil())
		})
	})

	var _ = Describe("GetModulesNames", func() {
		It("Sanity", func() {
			mta := &MTA{Modules: []*Modules{{Name: "someproj-db"}, {Name: "someproj-java"}}}
			Ω(mta.GetModulesNames()).Should(Equal([]string{"someproj-db", "someproj-java"}))
		})
	})

	var _ = Describe("Parsing", func() {
		It("Modules parsing - sanity", func() {
			var moduleSrv = Modules{
				Name: "srv",
				Type: "java",
				Path: "srv",
				Requires: []Requires{
					{
						Name: "db",
						Properties: Properties{
							"JBP_CONFIG_RESOURCE_CONFIGURATION": `[tomcat/webapps/ROOT/META-INF/context.xml: {"service_name_for_DefaultDB" : "~{hdi-container-name}"}]`,
						},
					},
				},
				Provides: []Provides{
					{
						Name:       "srv_api",
						Properties: Properties{"url": "${default-url}"},
					},
				},
				Parameters: Parameters{"memory": "512M"},
				Properties: Properties{
					"VSCODE_JAVA_DEBUG_LOG_LEVEL": "ALL",
					"APPC_LOG_LEVEL":              "info",
				},
			}
			var moduleUI = Modules{
				Name: "ui",
				Type: "html5",
				Path: "ui",
				Requires: []Requires{
					{
						Name:  "srv_api",
						Group: "destinations",
						Properties: Properties{
							"forwardAuthToken": true,
							"strictSSL":        false,
							"name":             "srv_api",
							"url":              "~{url}",
						},
					},
				},
				BuildParams: buildParameters{Builder: "grunt"},
				Parameters:  Parameters{"disk-quota": "256M", "memory": "256M"},
			}
			var modules = []*Modules{&moduleSrv, &moduleUI}
			mtaFile, _ := ioutil.ReadFile("./testdata/mta.yaml")
			// Parse file
			oMta := &MTA{}
			Ω(oMta.Parse(mtaFile)).Should(Succeed())
			Ω(oMta.Modules).Should(HaveLen(2))
			Ω(oMta.GetModules()).Should(Equal(modules))

		})

		It("BrokenMta", func() {
			mtaContent, _ := ioutil.ReadFile("./testdata/mtaWithBrokenProperties.yaml")
			oMta := &MTA{}
			Ω(oMta.Parse(mtaContent)).Should(HaveOccurred())
		})

		It("Full MTA Parsing - Sanity", func() {
			schemaVersion := "2.0.0"
			expected := MTA{
				SchemaVersion: &schemaVersion,
				ID:            "cloud.samples.someproj",
				Version:       "1.0.0",
				Parameters:    Parameters{"deploy_mode": "html5-repo"},
				Modules: []*Modules{
					{
						Name: "someproj-db",
						Type: "hdb",
						Path: "db",
						Requires: []Requires{
							{Name: "someproj-hdi-container"},
							{Name: "someproj-logging"},
						},
						Parameters: Parameters{"disk-quota": "256M", "memory": "256M"},
					},
					{
						Name:       "someproj-java",
						Type:       "java",
						Path:       "srv",
						Parameters: Parameters{"memory": "512M", "disk-quota": "256M"},
						Provides: []Provides{
							{
								Name:       "java",
								Properties: Properties{"url": "${default-url}"},
							},
						},
						Requires: []Requires{
							{
								Name: "someproj-hdi-container",
								Properties: Properties{
									"JBP_CONFIG_RESOURCE_CONFIGURATION": "[tomcat/webapps/ROOT/META-INF/context.xml: " +
										"{\"service_name_for_DefaultDB\" : \"~{hdi-container-name}\"}]",
								},
							},
							{Name: "someproj-logging"},
						},
						BuildParams: buildParameters{
							Requires: []BuildRequires{{Name: "someproj-db", TargetPath: ""}},
						},
					},
					{
						Name: "someproj-catalog-ui",
						Type: "html5",
						Path: "someproj-someprojCatalog",
						Parameters: Parameters{
							"memory":     "256M",
							"disk-quota": "256M",
						},
						Requires: []Requires{
							{
								Name:  "java",
								Group: "destinations",
								Properties: Properties{
									"name": "someproj-backend",
									"url":  "~{url}",
								},
							},
							{Name: "someproj-logging"},
						},
						BuildParams: buildParameters{
							Builder: "grunt",
							Requires: []BuildRequires{
								{
									Name:       "someproj-java",
									TargetPath: "",
								},
							},
						},
					},
					{
						Name: "someproj-uideployer",
						Type: "content",
						Parameters: Parameters{
							"memory":     "256M",
							"disk-quota": "256M",
						},
						Requires: []Requires{{Name: "someproj-apprepo-dt"}},
						BuildParams: buildParameters{
							Builder: "grunt",
							Type:    "content",
							Requires: []BuildRequires{
								{Name: "someproj-catalog-ui"},
							},
						},
					},
					{
						Name: "someproj",
						Type: "approuter.nodejs",
						Path: "approuter",
						Parameters: Parameters{
							"memory":     "256M",
							"disk-quota": "256M",
						},
						Requires: []Requires{
							{
								Name:  "java",
								Group: "destinations",
								Properties: Properties{
									"name": "someproj-backend",
									"url":  "~{url}",
								},
							},
							{Name: "someproj-apprepo-rt"},
							{Name: "someproj-logging"},
						},
					},
				},
				Resources: []*Resources{
					{
						Name:       "someproj-hdi-container",
						Properties: Properties{"hdi-container-name": "${service-name}"},
						Type:       "container",
					},
					{
						Name: "someproj-apprepo-rt",
						Type: "org.cloudfoundry.managed-service",
						Parameters: Parameters{
							"service":      "html5-apps-repo",
							"service-plan": "app-runtime",
						},
					},
					{
						Name: "someproj-apprepo-dt",
						Type: "org.cloudfoundry.managed-service",
						Parameters: Parameters{
							"service":      "html5-apps-repo",
							"service-plan": "app-host",
						},
					},
					{
						Name: "someproj-logging",
						Type: "org.cloudfoundry.managed-service",
						Parameters: Parameters{
							"service":      "application-logs",
							"service-plan": "lite",
						},
					},
				},
			}

			mtaContent, _ := ioutil.ReadFile("./testdata/mta2.yaml")

			actual := &MTA{}
			Ω(actual.Parse(mtaContent)).Should(Succeed())
			Ω(expected).Should(Equal(*actual))
		})
	})

	var _ = Describe("Get methods on MTA", func() {
		modules := []*Modules{
			{
				Name: "someproj-db",
				Type: "hdb",
				Path: "db",
				Requires: []Requires{
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
				Parameters: Parameters{
					"memory":     "512M",
					"disk-quota": "256M",
				},
			},
		}
		schemaVersion := "0.0.2"
		mta := &MTA{
			SchemaVersion: &schemaVersion,
			ID:            "MTA",
			Version:       "1.1.1",
			Modules:       modules,
			Resources: []*Resources{
				{
					Name: "someproj-hdi-container",
					Properties: Properties{
						"hdi-container-name": "${service-name}",
					},
					Type: "container",
				},
				{
					Name: "someproj-apprepo-rt",
					Type: "org.cloudfoundry.managed-service",
					Parameters: Parameters{
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

})
