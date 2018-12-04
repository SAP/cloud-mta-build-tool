package mta

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Mta", func() {
	var _ = Describe("MTA tests", func() {

		var _ = Describe("Parsing", func() {
			It("Modules parsing - sanity", func() {
				var moduleSrv = Module{
					Name: "srv",
					Type: "java",
					Path: "srv",
					Requires: []Requires{
						{
							Name: "db",
							Properties: map[string]interface{}{
								"JBP_CONFIG_RESOURCE_CONFIGURATION": `[tomcat/webapps/ROOT/META-INF/context.xml: {"service_name_for_DefaultDB" : "~{hdi-container-name}"}]`,
							},
						},
					},
					Provides: []Provides{
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
				var moduleUI = Module{
					Name: "ui",
					Type: "html5",
					Path: "ui",
					Requires: []Requires{
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
				var modules = []*Module{&moduleSrv, &moduleUI}
				mtaFile, _ := ioutil.ReadFile("./testdata/mta.yaml")
				// Unmarshal file
				oMta := &MTA{}
				Ω(yaml.Unmarshal(mtaFile, oMta)).Should(Succeed())
				Ω(oMta.Modules).Should(HaveLen(2))
				Ω(oMta.GetModules()).Should(Equal(modules))

			})

		})

		var _ = Describe("Get methods on MTA", func() {
			modules := []*Module{
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
					Parameters: map[string]interface{}{
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
				Resources: []*Resource{
					{
						Name: "someproj-hdi-container",
						Properties: map[string]interface{}{
							"hdi-container-name": "${service-name}",
						},
						Type: "container",
					},
					{
						Name: "someproj-apprepo-rt",
						Type: "org.cloudfoundry.managed-service",
						Parameters: map[string]interface{}{
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

	var _ = Describe("Unmarshal", func() {
		It("Sanity", func() {
			wd, err := os.Getwd()
			Ω(err).Should(Succeed())
			content, err := ioutil.ReadFile(filepath.Join(wd, "testdata", "mta.yaml"))
			Ω(err).Should(Succeed())
			m, err := Unmarshal(content)
			Ω(err).Should(Succeed())
			Ω(len(m.Modules)).Should(Equal(2))
		})
	})

	var _ = Describe("UnmarshalExt", func() {
		It("Sanity", func() {
			wd, err := os.Getwd()
			Ω(err).Should(Succeed())
			content, err := ioutil.ReadFile(filepath.Join(wd, "testdata", "mta.yaml"))
			Ω(err).Should(Succeed())
			m, err := UnmarshalExt(content)
			Ω(err).Should(Succeed())
			Ω(len(m.Modules)).Should(Equal(2))
		})
	})

	var _ = Describe("extendMap", func() {
		var m1 map[string]interface{}
		var m2 map[string]interface{}
		var m3 map[string]interface{}
		var m4 map[string]interface{}

		BeforeEach(func() {
			m1 = make(map[string]interface{})
			m2 = make(map[string]interface{})
			m3 = make(map[string]interface{})
			m4 = nil
			m1["a"] = "aa"
			m1["b"] = "xx"
			m2["b"] = "bb"
			m3["c"] = "cc"
		})

		var _ = DescribeTable("Sanity", func(m *map[string]interface{}, e *map[string]interface{}, ln int, key string, value interface{}) {
			extendMap(m, e)
			Ω(len(*m)).Should(Equal(ln))

			if value != nil {
				Ω((*m)[key]).Should(Equal(value))
			} else {
				Ω((*m)[key]).Should(BeNil())
			}
		},
			Entry("overwrite", &m1, &m2, 2, "b", "bb"),
			Entry("add", &m1, &m3, 3, "c", "cc"),
			Entry("res equals ext", &m4, &m1, 2, "b", "xx"),
			Entry("nothing to add", &m1, &m4, 2, "b", "xx"),
			Entry("both nil", &m4, &m4, 0, "b", nil),
		)
	})

	var _ = Describe("MergeMtaAndExt", func() {
		It("Sanity", func() {
			moduleA := Module{
				Name: "modA",
				Properties: map[string]interface{}{
					"a": "aa",
					"b": "xx",
				},
			}
			moduleB := Module{
				Name: "modB",
				Properties: map[string]interface{}{
					"b": "yy",
				},
			}
			moduleAExt := ModuleExt{
				Name: "modA",
				Properties: map[string]interface{}{
					"a": "aa",
					"b": "bb",
				},
			}
			mta := MTA{
				Modules: []*Module{&moduleA, &moduleB},
			}
			mtaExt := MTAExt{
				Modules: []*ModuleExt{&moduleAExt},
			}
			Merge(&mta, &mtaExt)
			m, err := mta.GetModuleByName("modA")
			Ω(err).Should(Succeed())
			Ω(m.Properties["b"]).Should(Equal("bb"))
		})
	})

})
