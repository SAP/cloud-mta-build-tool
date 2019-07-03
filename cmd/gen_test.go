package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
	"path/filepath"
)

var _ = Describe("Generate commands call", func() {

	var ep dir.Loc

	BeforeEach(func() {
		mtadCmdTrg = getTestPath("result")
		metaCmdTrg = getTestPath("result")
		mtarCmdTrg = getTestPath("result")
		err := dir.CreateDirIfNotExist(mtadCmdTrg)
		Ω(err).Should(Succeed())
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
	It("Generate Mtad - Sanity", func() {
		createDirInTmpFolder("mtahtml5", "ui5app")
		createDirInTmpFolder("mtahtml5", "ui5app2")
		mtadCmdSrc = getTestPath("mtahtml5")
		mtadCmdPlatform = "cf"
		Ω(mtadCmd.RunE(nil, []string{})).Should(Succeed())
		Ω(filepath.Join(getTestPath("result"), "mtad.yaml")).Should(BeAnExistingFile())
	})
	It("Generate Mtad - Invalid source", func() {
		mtadCmdSrc = getTestPath("mtahtml6")
		Ω(mtadCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})
	It("Generate Mtad - Invalid platform configuration", func() {
		mtadCmdSrc = getTestPath("mtahtml5")
		config := platform.PlatformConfig
		platform.PlatformConfig = []byte("wrong config")
		mtadCmdPlatform = "cf"
		Ω(mtadCmd.RunE(nil, []string{})).Should(HaveOccurred())
		platform.PlatformConfig = config
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
