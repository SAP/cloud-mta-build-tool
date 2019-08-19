package artifacts

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Mtar", func() {
	var _ = Describe("Generate", func() {

		AfterEach(func() {
			os.RemoveAll(getResultPath())
		})

		var _ = Describe("ExecuteGenMtar", func() {
			It("Sanity, target provided", func() {
				createMtahtml5TmpFolder()
				Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", nil, "cf", os.Getwd)).Should(Succeed())
				Ω(ExecuteGenMtar(getTestPath("mtahtml5"), getResultPath(), "true", "dev", nil, "", os.Getwd)).Should(Succeed())
				Ω(getTestPath("result", "mtahtml5_0.0.1.mtar")).Should(BeAnExistingFile())
			})

			It("Sanity, target not provided", func() {
				createMtahtml5TmpFolder()
				Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", nil, "cf", os.Getwd)).Should(Succeed())
				Ω(ExecuteGenMtar(getTestPath("mtahtml5"), getResultPath(), "false", "dev", nil, "", os.Getwd)).Should(Succeed())
				Ω(getTestPath("result", "mta_archives", "mtahtml5_0.0.1.mtar")).Should(BeAnExistingFile())
			})

			It("Fails on location initialization", func() {
				Ω(ExecuteGenMtar("", getResultPath(), "true", "dev", nil, "", func() (string, error) {
					return "", errors.New("err")
				})).Should(HaveOccurred())
			})

			It("Fails - wrong source", func() {
				Ω(ExecuteGenMtar(getTestPath("mtahtml6"), getResultPath(), "true", "dev", nil, "", os.Getwd)).Should(HaveOccurred())
			})
		})

		It("Generate Mtar - Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			createMtahtml5TmpFolder()
			Ω(generateMeta(&ep, &ep, false, "cf", true)).Should(Succeed())
			mtarPath, err := generateMtar(&ep, &ep, &ep, true, "")
			Ω(err).Should(Succeed())
			Ω(mtarPath).Should(BeAnExistingFile())
		})

		It("Generate Mtar - Fails on wrong source", func() {
			ep := dir.Loc{SourcePath: getTestPath("not_existing"), TargetPath: getResultPath()}
			ep1 := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			_, err := generateMtar(&ep, &ep1, &ep1, true, "")
			Ω(err).Should(HaveOccurred())
		})

		It("Generate Mtar - Invalid mta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(), MtaFilename: "mtaBroken.yaml"}
			_, err := generateMtar(&ep, &ep, &ep, true, "")
			Ω(err).Should(HaveOccurred())
		})
		It("Generate Mtar - Mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(), MtaFilename: "mtaNotExists.yaml"}
			_, err := generateMtar(&ep, &ep, &ep, true, "")
			Ω(err).Should(HaveOccurred())
		})
		It("Generate Mtar - results file exists, folder results can't be created ", func() {
			file, _ := os.Create(getTestPath("result"))
			defer file.Close()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			_, err := generateMtar(&ep, &ep, &ep, true, "")
			Ω(err).Should(HaveOccurred())
		})
		DescribeTable("isTargetProvided", func(target, provided string, expected bool) {
			Ω(isTargetProvided(target, provided)).Should(Equal(expected))
		},
			Entry("Sanity", "", "true", true),
			Entry("Wrong provided value", "", "xx", false),
			Entry("Empty provided value, target path provided", "path", "", true),
			Entry("Empty provided value, no target path provided", "", "", false),
		)

		var _ = Describe("Target Failures", func() {
			var _ = DescribeTable("Invalid location", func(loc *testMtarLoc) {
				ep := dir.Loc{}
				_, err := generateMtar(loc, &ep, &ep, true, "")
				Ω(err).Should(HaveOccurred())
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

	var _ = DescribeTable("getMtarFileName", func(mtarName, expected string) {
		m := mta.MTA{ID: "proj", Version: "0.1.5"}
		Ω(getMtarFileName(&m, mtarName)).Should(Equal(expected))
	},
		Entry("default mtar", "", "proj_0.1.5.mtar"),
		Entry("default supporting make file", "*", "proj_0.1.5.mtar"),
		Entry("default supporting make file", "abc", "abc.mtar"),
		Entry("default supporting make file", "abc.zip", "abc.zip"),
	)
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
