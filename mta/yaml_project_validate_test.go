package mta

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateYamlProject", func() {
	It("Sanity", func() {
		wd, _ := os.Getwd()
		ep := dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata", "testproject")}
		mta, _ := ReadMta(&ep)
		source, _ := ep.GetSource()
		issues := validateYamlProject(mta, source)
		Ω(issues[0].Msg).Should(Equal("Module <ui5app2> not found in project. Expected path: <ui5app2>"))
	})
})
