package artifacts

import (
	"cloud-mta-build-tool/mta"
)

var simpleModulesList = []*mta.Module{
	{
		Name: "ui5",
		Type: "html5",
		Path: "ui5",
	}}

//type testWriter struct {
//	call       int
//	failOnCall int
//	writer     io.Writer
//}
//
//func (writer *testWriter) Write(p []byte) (n int, err error) {
//	writer.call++
//	if writer.call < writer.failOnCall {
//		return writer.writer.Write(p)
//	}
//	return 0, errors.New("error")
//}

//var _ = Describe("setManifestDesc", func() {
//	var _ = DescribeTable("Sanity", func(args []*mta.Module, expected string, modules []string) {
//		b := &bytes.Buffer{}
//		setManifestDesc(b, args, modules)
//		fmt.Println(b.String())
//		Ω(b.String()).Should(Equal(expected))
//	},
//		Entry("One module", simpleModulesList, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
//			"Name: ui5/data.zip\nMTA-Module: ui5\nContent-Type: application/zip",
//			[]string{}),
//		Entry(" Two modules", []*mta.Module{
//			{
//				Name: "ui6",
//				Type: "html5",
//				Path: "ui5",
//			},
//			{
//				Name: "ui4",
//				Type: "html5",
//				Path: "ui5",
//			}}, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
//			"Name: ui6/data.zip\nMTA-Module: ui6\nContent-Type: application/zip\n\n"+
//			"Name: ui4/data.zip\nMTA-Module: ui4\nContent-Type: application/zip",
//			[]string{}),
//		Entry(" multi module with filter of one", []*mta.Module{
//			{
//				Name: "ui6",
//				Type: "html5",
//				Path: "ui5",
//			},
//			{
//				Name: "ui4",
//				Type: "html5",
//				Path: "ui5",
//			}}, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
//			"Name: ui6/data.zip\nMTA-Module: ui6\nContent-Type: application/zip", []string{"ui6"}),
//	)
//
//	var _ = DescribeTable("Invalid cases", func(failOn int, modules []string) {
//		w := testWriter{
//			failOnCall: failOn,
//			call:       0,
//			writer:     &bytes.Buffer{},
//		}
//		Ω(setManifestDesc(&w, simpleModulesList, modules)).Should(HaveOccurred())
//	},
//		Entry("Fails on 1st line", 1, []string{}),
//		Entry("Fails on version line", 2, []string{}),
//		Entry("Fails on modules printing with empty modules list", 3, []string{}),
//		Entry("Fails on modules printing with not empty modules list", 3, []string{"ui5"}),
//	)
//
//	var _ = Describe("Failure", func() {
//		var config []byte
//
//		BeforeEach(func() {
//			config = make([]byte, len(version.VersionConfig))
//			copy(config, version.VersionConfig)
//			// Simplified commands configuration (performance purposes). removed "npm prune --production"
//			version.VersionConfig = []byte(`
//cli_version:["x"]
//`)
//		})
//
//		AfterEach(func() {
//			version.VersionConfig = make([]byte, len(config))
//			copy(version.VersionConfig, config)
//		})
//
//		It("Get version fails", func() {
//			Ω(setManifestDesc(os.Stdout, simpleModulesList, []string{})).Should(HaveOccurred())
//		})
//	})
//
//})
