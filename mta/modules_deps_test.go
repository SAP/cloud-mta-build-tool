package mta

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModulesDeps", func() {

	It("Resolve dependencies - Valid case", func() {
		wd, _ := os.Getwd()
		mtaStr, _ := ReadMta(dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps.yaml"})
		actual, _ := mtaStr.GetModulesOrder()
		// last module depends on others
		Ω(actual[len(actual)-1]).Should(Equal("eb-uideployer"))
	})

	It("Resolve dependencies - cyclic dependencies", func() {
		wd, _ := os.Getwd()
		mtaStr, _ := ReadMta(dir.MtaLocationParameters{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps_cyclic_deps.yaml"})
		_, err := mtaStr.GetModulesOrder()
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("eb-ui-conf-eb"))
	})

})
