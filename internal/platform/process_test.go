package platform

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Process", func() {

	var platforms = Platforms{[]Modules{
		{Name: "cf",
			Modules: []Properties{
				{NativeType: "html5", PlatformType: "javascript.nodejs"},
				{NativeType: "nodejs", PlatformType: "javascript.nodejs"},
				{NativeType: "java", PlatformType: "java.tomcat"},
				{NativeType: "hdb", PlatformType: "dbtype"},
				{NativeType: "hdb", PlatformType: "dbtype2", Parameters: map[string]string{"value": "2"}},
				{NativeType: "java", PlatformType: "java.tomee", Properties: map[string]string{"TARGET_RUNTIME": "tomee"}},
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
   - native-type: hdb
     platform-type: "dbtype2"
     parameters:
       value: 2
   - native-type: java
     platform-type: "java.tomee"
     properties:
       TARGET_RUNTIME: tomee

 - name: neo
   modules:
   - native-type: html5
     platform-type: "some.html"
   - native-type: java
     platform-type: "java.tomcat"
`)
		Ω(Unmarshal(platformsCfg)).Should(Equal(platforms))

	})

	It("Unmarshal - wrong elements", func() {
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
     platform-typex: "dbtype"
`)
		_, err := Unmarshal(platformsCfg)
		Ω(err).Should(HaveOccurred())

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
				{Name: "java2", Type: "java", Path: "app4", Properties: map[string]interface{}{"TARGET_RUNTIME": "tomee"}},
				{Name: "hdb2", Type: "hdb", Path: "app5", Parameters: map[string]interface{}{"value": "2"}},
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
				{Name: "java2", Type: "java.tomcat", Path: "app4", Properties: map[string]interface{}{"TARGET_RUNTIME": "tomee"}},
				{Name: "hdb2", Type: "hdb", Path: "app5", Parameters: map[string]interface{}{"value": "2"}},
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
				{Name: "java2", Type: "java.tomee", Path: "app4", Properties: map[string]interface{}{"TARGET_RUNTIME": "tomee"}},
				{Name: "hdb2", Type: "dbtype2", Path: "app5", Parameters: map[string]interface{}{"value": "2"}},
			},
		}
		ConvertTypes(mtaObj, platforms, platform)
		Ω(mtaObj).Should(Equal(mtaObjMap[platform]))
	},
		Entry("Neo", "neo"),
		Entry("CF", "cf"),
	)

	It("ConvertTypes returns the most accurate type when it appears before the less accurate type", func() {
		schemaVersion := "2.0.0"
		mtaObj := mta.MTA{
			SchemaVersion: &schemaVersion,
			ID:            "mta_proj",
			Version:       "1.0.0",
			Modules: []*mta.Module{
				{Name: "java", Type: "java", Path: "app4", Properties: map[string]interface{}{"TARGET_RUNTIME": "tomee"}},
			},
		}
		expected := mta.MTA{
			SchemaVersion: &schemaVersion,
			ID:            "mta_proj",
			Version:       "1.0.0",
			Modules: []*mta.Module{
				{Name: "java", Type: "java.tomee", Path: "app4", Properties: map[string]interface{}{"TARGET_RUNTIME": "tomee"}},
			},
		}
		platforms := Platforms{[]Modules{
			{Name: "cf",
				Modules: []Properties{
					{NativeType: "java", PlatformType: "java.tomee", Properties: map[string]string{"TARGET_RUNTIME": "tomee"}},
					{NativeType: "java", PlatformType: "java.tomcat"},
				},
			},
		}}
		ConvertTypes(mtaObj, platforms, "cf")
		Ω(mtaObj).Should(Equal(expected))
	})

	It("platformConfig", func() {
		expected := Modules{Name: "cf",
			Modules: []Properties{
				{NativeType: "html5", PlatformType: "javascript.nodejs"},
				{NativeType: "nodejs", PlatformType: "javascript.nodejs"},
				{NativeType: "java", PlatformType: "java.tomcat"},
				{NativeType: "hdb", PlatformType: "dbtype"},
				{NativeType: "hdb", PlatformType: "dbtype2", Parameters: map[string]string{"value": "2"}},
				{NativeType: "java", PlatformType: "java.tomee", Properties: map[string]string{"TARGET_RUNTIME": "tomee"}},
			},
		}
		Ω(platformConfig(platforms, "cf")).Should(Equal(expected))
	})

	Describe("satisfiesModuleConfig", func() {
		It("returns ok=false when module type is mismatched", func() {
			m := mta.Module{Name: "htmlapp", Type: "javascript.nodejs", Path: "app"}
			config := Properties{NativeType: "abcd", PlatformType: "abcd"}
			ok, acc := satisfiesModuleConfig(&m, &config)
			Ω(ok).Should(BeFalse())
			Ω(acc).Should(BeNumerically("<", 0))
		})

		It("returns ok=true when module type matches and no other conditions", func() {
			m := mta.Module{Type: "htmlapp"}
			config := Properties{NativeType: "htmlapp", PlatformType: "abcd"}
			ok, acc := satisfiesModuleConfig(&m, &config)
			Ω(ok).Should(BeTrue())
			Ω(acc).Should(BeNumerically(">=", 0))
		})

		DescribeTable("returns ok=false when module type matches and property or parameter doesn't match", func(c Properties) {
			m := mta.Module{
				Type:       "a",
				Properties: map[string]interface{}{"a": "b"},
				Parameters: map[string]interface{}{"a": "b"},
			}
			c.NativeType = "a"
			ok, acc := satisfiesModuleConfig(&m, &c)
			Ω(ok).Should(BeFalse())
			Ω(acc).Should(BeNumerically("<", 0))
		},
			Entry("property doesn't match", Properties{Properties: map[string]string{"a": "a"}}),
			Entry("property doesn't exist", Properties{Properties: map[string]string{"b": "a"}}),
			Entry("parameter doesn't match", Properties{Parameters: map[string]string{"a": "a"}}),
			Entry("parameter doesn't exist", Properties{Parameters: map[string]string{"b": "a"}}),
			Entry("One parameter matches and one doesn't", Properties{Parameters: map[string]string{"a": "b", "b": "a"}}),
			Entry("One property matches and one doesn't", Properties{Properties: map[string]string{"a": "b", "b": "a"}}),
			Entry("Parameter matches and property doesn't", Properties{Parameters: map[string]string{"a": "b"}, Properties: map[string]string{"a": "a"}}),
			Entry("Property matches and parameter doesn't", Properties{Parameters: map[string]string{"b": "b"}, Properties: map[string]string{"a": "b"}}),
		)

		It("returns ok=true and accuracy is higher the more conditions there are", func() {
			configs := []struct {
				desc   string
				config Properties
			}{
				{"no properties or parameters", Properties{}},
				{"1 property", Properties{Properties: map[string]string{"a": "a"}}},
				{"2 properties", Properties{Properties: map[string]string{"a": "a", "b": "b"}}},
				{"3 parameters", Properties{Parameters: map[string]string{"c": "c", "d": "d", "e": "e"}}},
				{"3 parameters and a property", Properties{Parameters: map[string]string{"c": "c", "d": "d", "e": "e"}, Properties: map[string]string{"c": "c"}}},
			}
			m := mta.Module{
				Type:       "a",
				Properties: map[string]interface{}{"a": "a", "b": "b", "c": "c"},
				Parameters: map[string]interface{}{"c": "c", "d": "d", "e": "e"},
			}
			acc := -1
			prevDesc := "initial value"
			for _, c := range configs {
				c.config.NativeType = "a"
				ok, moduleAcc := satisfiesModuleConfig(&m, &c.config)
				Ω(ok).Should(BeTrue(), "module did not satisfy config with "+c.desc)
				Ω(moduleAcc).Should(BeNumerically(">", acc), c.desc+" has lower accuracy than "+prevDesc)
				acc = moduleAcc
				prevDesc = c.desc
			}
		})
	})
})
