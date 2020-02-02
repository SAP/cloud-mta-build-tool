package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
)

var _ = Describe("Generate commands call", func() {

	var ep dir.Loc

	BeforeEach(func() {
		metaCmdTrg = getTestPath("result")
		mtarCmdTrg = getTestPath("result")
	})

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})

	It("Generate Meta", func() {
		createDirInTmpFolder("mtahtml5", "ui5app2")
		createDirInTmpFolder("mtahtml5", "ui5app")
		createFileInTmpFolder("mtahtml5", "ui5app2", "data.zip")
		createFileInTmpFolder("mtahtml5", "xs-security.json")
		createFileInTmpFolder("mtahtml5", "ui5app", "data.zip")
		metaCmdSrc = getTestPath("mtahtml5")
		ep = dir.Loc{SourcePath: metaCmdSrc, TargetPath: metaCmdTrg}
		Ω(metaCmd.RunE(nil, []string{})).Should(Succeed())
		Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
	})
	It("Generate Mtar", func() {
		createDirInTmpFolder("mtahtml5", "ui5app2")
		createDirInTmpFolder("mtahtml5", "ui5app")
		createFileInTmpFolder("mtahtml5", "ui5app", "data.zip")
		createFileInTmpFolder("mtahtml5", "ui5app2", "data.zip")
		createFileInTmpFolder("mtahtml5", "xs-security.json")
		mtarCmdSrc = getTestPath("mtahtml5")
		mtarCmdMtarName = ""
		Ω(metaCmd.RunE(nil, []string{})).Should(Succeed())
		Ω(mtarCmd.RunE(nil, []string{})).Should(Succeed())
		Ω(getTestPath("result", "mtahtml5_0.0.1.mtar")).Should(BeAnExistingFile())
	})
})
