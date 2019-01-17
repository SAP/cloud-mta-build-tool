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
			os.RemoveAll(getResultPath())
		})

		var _ = Describe("ExecuteGenMtar", func() {
			It("Sanity", func() {
				os.MkdirAll(getTestPath("result", "mtahtml5", "testapp"), os.ModePerm)
				os.MkdirAll(getTestPath("result", "mtahtml5", "ui5app2"), os.ModePerm)
				Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", "cf", true, os.Getwd)).Should(Succeed())
				Ω(ExecuteGenMtar(getTestPath("mtahtml5"), getResultPath(), "dev", os.Getwd)).Should(Succeed())
				Ω(getTestPath("result", "mtahtml5.mtar")).Should(BeAnExistingFile())
			})

			It("Fails on location initialization", func() {
				Ω(ExecuteGenMtar("", getResultPath(), "dev", func() (string, error) {
					return "", errors.New("err")
				})).Should(HaveOccurred())
			})

			It("Fails - wrong source", func() {
				Ω(ExecuteGenMtar(getTestPath("mtahtml6"), getResultPath(), "dev", os.Getwd)).Should(HaveOccurred())
			})
		})

		It("Generate Mtar - Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, nil, false, "cf", true)).Should(Succeed())
			Ω(generateMtar(&ep, &ep)).Should(Succeed())
			Ω(getTestPath("result", "mtahtml5.mtar")).Should(BeAnExistingFile())
		})

		It("Generate Mtar - Fails on wrong source", func() {
			ep := dir.Loc{SourcePath: getTestPath("not_existing"), TargetPath: getResultPath()}
			ep1 := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMtar(&ep, &ep1)).Should(HaveOccurred())
		})

		It("Generate Mtar - Invalid mta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(), MtaFilename: "mtaBroken.yaml"}
			Ω(generateMtar(&ep, &ep)).Should(HaveOccurred())
		})
		It("Generate Mtar - Mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(), MtaFilename: "mtaNotExists.yaml"}
			Ω(generateMtar(&ep, &ep)).Should(HaveOccurred())
		})

		var _ = Describe("Target Failures", func() {
			var _ = DescribeTable("Invalid location", func(loc *testMtarLoc) {
				ep := dir.Loc{}
				Ω(generateMtar(loc, &ep)).Should(HaveOccurred())
			},
				Entry("Fails on GetTargetTmpDir", &testMtarLoc{
					tmpDir:    "",
					targetDir: getResultPath(),
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
