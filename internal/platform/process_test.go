package platform

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/mta"
)

var _ = Describe("Process", func() {

	var platforms = Platforms{[]Modules{
		{Name: "cf",
			Modules: []Properties{
				{NativeType: "html5", PlatformType: "javascript.nodejs"},
				{NativeType: "nodejs", PlatformType: "javascript.nodejs"},
				{NativeType: "java", PlatformType: "java.tomcat"},
				{NativeType: "hdb", PlatformType: "dbtype"},
			},
		},
		{Name: "neo",
			Modules: []Properties{
				{NativeType: "html5", PlatformType: "some.html"},
				{NativeType: "java", PlatformType: "java.tomcat"},
			},
		},
	}}

	It("Unmarshal", func() {
		var platformsCfg = []byte(`
platform:
 - name: cf
   modules:
   - native-type: html5
     platform-type: "javascript.nodejs"
   - native-type: nodejs
     platform-type: "javascript.nodejs"
   - native-type: java
     platform-type: "java.tomcat"
   - native-type: hdb
     platform-type: "dbtype"

 - name: neo
   modules:
   - native-type: html5
     platform-type: "some.html"
   - native-type: java
     platform-type: "java.tomcat"
`)
		Ω(Parse(platformsCfg)).Should(Equal(platforms))

	})

	var _ = DescribeTable("ConvertTypes", func(platform string) {
		schemaVersion := "2.0.0"
		mtaObj := mta.MTA{
			SchemaVersion: &schemaVersion,
			ID:            "mta_proj",
			Version:       "1.0.0",
			Modules: []*mta.Module{
				{Name: "htmlapp", Type: "html5", Path: "app"},
				{Name: "htmlapp2", Type: "html5", Path: "app2"},
				{Name: "java", Type: "java", Path: "app3"},
			},
		}
		mtaObjMap := make(map[string]mta.MTA)
		mtaObjMap["neo"] = mta.MTA{
			SchemaVersion: &schemaVersion,
			ID:            "mta_proj",
			Version:       "1.0.0",
			Modules: []*mta.Module{
				{Name: "htmlapp", Type: "some.html", Path: "app"},
				{Name: "htmlapp2", Type: "some.html", Path: "app2"},
				{Name: "java", Type: "java.tomcat", Path: "app3"},
			},
		}
		mtaObjMap["cf"] = mta.MTA{
			SchemaVersion: &schemaVersion,
			ID:            "mta_proj",
			Version:       "1.0.0",
			Modules: []*mta.Module{
				{Name: "htmlapp", Type: "javascript.nodejs", Path: "app"},
				{Name: "htmlapp2", Type: "javascript.nodejs", Path: "app2"},
				{Name: "java", Type: "java.tomcat", Path: "app3"},
			},
		}
		ConvertTypes(mtaObj, platforms, platform)
		Ω(mtaObj).Should(Equal(mtaObjMap[platform]))
	},
		Entry("Neo", "neo"),
		Entry("CF", "cf"),
	)

	It("platformConfig", func() {
		expected := Modules{Name: "cf",
			Modules: []Properties{
				{NativeType: "html5", PlatformType: "javascript.nodejs"},
				{NativeType: "nodejs", PlatformType: "javascript.nodejs"},
				{NativeType: "java", PlatformType: "java.tomcat"},
				{NativeType: "hdb", PlatformType: "dbtype"},
			},
		}
		Ω(platformConfig(platforms, "cf")).Should(Equal(expected))
	})
})
