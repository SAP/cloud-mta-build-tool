package mta

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
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
				var moduleUI = Module{
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
					BuildParams: BuildParameters{Builder: "grunt"},
					Parameters:  Parameters{"disk-quota": "256M", "memory": "256M"},
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
				Resources: []*Resource{
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
})
