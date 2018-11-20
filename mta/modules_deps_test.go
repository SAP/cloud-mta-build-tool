package mta

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModulesDeps", func() {

	It("Resolve dependencies - Valid case", func() {
		wd, _ := os.Getwd()
		mtaStr, _ := ReadMta(&Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps.yaml"})
		actual, _ := mtaStr.getModulesOrder()
		// last module depends on others
		Ω(actual[len(actual)-1]).Should(Equal("eb-uideployer"))
	})

	It("Resolve dependencies - cyclic dependencies", func() {
		wd, _ := os.Getwd()
		mtaStr, _ := ReadMta(&Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps_cyclic_deps.yaml"})
		_, err := mtaStr.getModulesOrder()
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("eb-ui-conf-eb"))
	})

})
