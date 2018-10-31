package mta

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Desc tests", func() {

	var _ = DescribeTable("setManifetDesc", func(args []*Modules, expected string, modules []string) {
		b := &bytes.Buffer{}
		setManifetDesc(b, args, modules)
		fmt.Println(b.String())
		Ω(b.String()).Should(Equal(expected))
	},
		Entry("One module", []*Modules{
			{
				Name: "ui5",
				Type: "html5",
				Path: "ui5",
			}}, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
			"Name: ui5/data.zip\nMTA-Module: ui5\nContent-Type: application/zip",
			[]string{}),
		Entry(" Two modules", []*Modules{
			{
				Name: "ui6",
				Type: "html5",
				Path: "ui5",
			},
			{
				Name: "ui4",
				Type: "html5",
				Path: "ui5",
			}}, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
			"Name: ui6/data.zip\nMTA-Module: ui6\nContent-Type: application/zip\n\n"+
			"Name: ui4/data.zip\nMTA-Module: ui4\nContent-Type: application/zip",
			[]string{}),
		Entry(" multi module with filter of one", []*Modules{
			{
				Name: "ui6",
				Type: "html5",
				Path: "ui5",
			},
			{
				Name: "ui4",
				Type: "html5",
				Path: "ui5",
			}}, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
			"Name: ui6/data.zip\nMTA-Module: ui6\nContent-Type: application/zip", []string{"ui6"}),
	)

	var _ = Describe("GenMetaInf", func() {
		wd, _ := os.Getwd()
		ep := dir.EndPoints{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}

		AfterEach(func() {
			os.RemoveAll(ep.GetTarget())
		})

		It("Sanity", func() {
			var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)
			mta := MTA{}
			yaml.Unmarshal(mtaSingleModule, &mta)
			GenMetaInfo(ep, mta, []string{"htmlapp"}, func(mtaStr MTA) {})
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})
	})
})
