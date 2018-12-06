package artifacts

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
)

var _ = Describe("Mtar", func() {
	var _ = Describe("Generate", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		It("Generate Mtar - Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			Ω(GenerateMeta(&ep, "cf")).Should(Succeed())
			Ω(GenerateMtar(&ep)).Should(Succeed())
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})

		It("Generate Mtar - Invalid mta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaBroken.yaml"}
			Ω(GenerateMtar(&ep)).Should(HaveOccurred())
		})
		It("Generate Mtar - Mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaNotExists.yaml"}
			Ω(GenerateMtar(&ep)).Should(HaveOccurred())
		})

		var _ = Describe("Failures on GetWorkingDirectory", func() {
			AfterEach(func() {
				dir.GetWorkingDirectory = dir.OsGetWd
			})

			var _ = DescribeTable("Invalid location", func(failOnCall int) {
				call := 0
				dir.GetWorkingDirectory = func() (string, error) {
					if call >= failOnCall {
						return "", errors.New("error")
					}
					call++
					return getTestPath("mtahtml5"), nil
				}
				ep := dir.Loc{}
				Ω(GenerateMtar(&ep)).Should(HaveOccurred())
			},
				Entry("Fails on GetTargetTmpDir", 1),
				Entry("Fails on GetTarget", 3))
		})
	})
})
