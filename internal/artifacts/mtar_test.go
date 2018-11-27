package artifacts

import (
	"fmt"
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

		It("Generate Mtar", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			err := GenerateMeta(&ep)
			if err != nil {
				fmt.Println(err)
			}
			err = GenerateMtar(&ep)
			if err != nil {
				fmt.Println(err)
			}
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
})
