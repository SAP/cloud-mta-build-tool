package mta

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateYamlProject", func() {
	It("Sanity", func() {
		wd, _ := os.Getwd()
		ep := Loc{SourcePath: filepath.Join(wd, "testdata", "testproject")}
		mta, _ := ParseFile(&ep)
		source, _ := ep.GetSource()
		issues := validateYamlProject(mta, source)
		Ω(issues[0].Msg).Should(Equal("Module <ui5app2> not found in project. Expected path: <ui5app2>"))
	})
})
