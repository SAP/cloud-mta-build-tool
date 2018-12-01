package artifacts

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/mta"
)

var _ = Describe("Manifest", func() {
	var _ = DescribeTable("setManifetDesc", func(args []*mta.Module, expected string, modules []string) {
		b := &bytes.Buffer{}
		setManifetDesc(b, args, modules)
		fmt.Println(b.String())
		Ω(b.String()).Should(Equal(expected))
	},
		Entry("One module", []*mta.Module{
			{
				Name: "ui5",
				Type: "html5",
				Path: "ui5",
			}}, "manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n"+
			"Name: ui5/data.zip\nMTA-Module: ui5\nContent-Type: application/zip",
			[]string{}),
		Entry(" Two modules", []*mta.Module{
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
		Entry(" multi module with filter of one", []*mta.Module{
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

})
