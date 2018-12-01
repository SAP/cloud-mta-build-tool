package artifacts

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fsys"
)

var _ = Describe("Mtar", func() {
	var _ = Describe("Generate", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		It("Generate Mtar - Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			Ω(GenerateMeta(&ep)).Should(Succeed())
			Ω(GenerateMtar(&ep)).Should(Succeed())
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
})
