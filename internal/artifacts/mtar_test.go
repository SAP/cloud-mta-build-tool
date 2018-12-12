package artifacts

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fs"
)

var _ = Describe("Mtar", func() {
	var _ = Describe("Generate", func() {

		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		var _ = Describe("ExecuteGenMtar", func() {
			It("Sanity", func() {
				Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getTestPath("result"), "dev", "cf", os.Getwd)).Should(Succeed())
				Ω(ExecuteGenMtar(getTestPath("mtahtml5"), getTestPath("result"), "dev", os.Getwd)).Should(Succeed())
				Ω(getTestPath("result", "mtahtml5.mtar")).Should(BeAnExistingFile())
			})

			It("Fails on location initialization", func() {
				Ω(ExecuteGenMtar("", getTestPath("result"), "dev", func() (string, error) {
					return "", errors.New("err")
				})).Should(HaveOccurred())
			})

			It("Fails - wrong source", func() {
				Ω(ExecuteGenMtar(getTestPath("mtahtml6"), getTestPath("result"), "dev", os.Getwd)).Should(HaveOccurred())
			})
		})

		It("Generate Mtar - Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result")}
			Ω(generateMeta(&ep, &ep, false, "cf")).Should(Succeed())
			Ω(generateMtar(&ep, &ep)).Should(Succeed())
			Ω(getTestPath("result", "mtahtml5.mtar")).Should(BeAnExistingFile())
		})

		It("Generate Mtar - Invalid mta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaBroken.yaml"}
			Ω(generateMtar(&ep, &ep)).Should(HaveOccurred())
		})
		It("Generate Mtar - Mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaNotExists.yaml"}
			Ω(generateMtar(&ep, &ep)).Should(HaveOccurred())
		})

		var _ = Describe("Target Failures", func() {
			var _ = DescribeTable("Invalid location", func(loc *testMtarLoc) {
				ep := dir.Loc{}
				Ω(generateMtar(loc, &ep)).Should(HaveOccurred())
			},
				Entry("Fails on GetTargetTmpDir", &testMtarLoc{
					tmpDir:    "",
					targetDir: getTestPath("result"),
				}),
				Entry("Fails on GetTarget", &testMtarLoc{
					tmpDir:    getTestPath("result", "mtahtml5", "mtahtml5"),
					targetDir: "",
				}))
		})
	})
})

type testMtarLoc struct {
	tmpDir    string
	targetDir string
}

func (loc *testMtarLoc) GetTarget() string {
	return loc.targetDir
}
func (loc *testMtarLoc) GetTargetTmpDir() string {
	return loc.tmpDir
}
