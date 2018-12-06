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
			Ω(GenerateMtar(&ep, &ep)).Should(Succeed())
			mtarPath := getTestPath("result", "mtahtml5.mtar")
			Ω(mtarPath).Should(BeAnExistingFile())
		})

		It("Generate Mtar - Invalid mta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaBroken.yaml"}
			Ω(GenerateMtar(&ep, &ep)).Should(HaveOccurred())
		})
		It("Generate Mtar - Mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaNotExists.yaml"}
			Ω(GenerateMtar(&ep, &ep)).Should(HaveOccurred())
		})

		var _ = Describe("Target Failures", func() {
			var _ = DescribeTable("Invalid location", func(loc *testMtarLoc) {
				ep := dir.Loc{}
				Ω(GenerateMtar(loc, &ep)).Should(HaveOccurred())
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

func (loc *testMtarLoc) GetTarget() (string, error) {
	if loc.targetDir == "" {
		return "", errors.New("err")
	}
	return loc.targetDir, nil
}
func (loc *testMtarLoc) GetTargetTmpDir() (string, error) {
	if loc.tmpDir == "" {
		return "", errors.New("err")
	}
	return loc.tmpDir, nil
}
