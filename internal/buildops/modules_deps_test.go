package buildops

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("ModulesDeps", func() {

	It("Resolve dependencies - Valid case", func() {
		wd, _ := os.Getwd()
		mtaStr, _ := dir.ParseFile(&dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps.yaml"})
		actual, _ := getModulesOrder(mtaStr)
		// last module depends on others
		Ω(actual[len(actual)-1]).Should(Equal("eb-uideployer"))
	})

	It("Resolve dependencies - cyclic dependencies", func() {
		wd, _ := os.Getwd()
		mtaStr, _ := dir.ParseFile(&dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps_cyclic_deps.yaml"})
		_, err := getModulesOrder(mtaStr)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("eb-ui-conf-eb"))
	})

	var _ = Describe("GetModulesNames", func() {
		It("Sanity", func() {
			mtaStr := &mta.MTA{Modules: []*mta.Module{{Name: "someproj-db"}, {Name: "someproj-java"}}}
			Ω(GetModulesNames(mtaStr)).Should(Equal([]string{"someproj-db", "someproj-java"}))
		})
	})

})
