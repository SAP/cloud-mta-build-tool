package commands

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

var _ = Describe("Commands", func() {

	BeforeEach(func() {
		mtadCmdTrg = getTestPath("result")
		metaCmdTrg = getTestPath("result")
		mtarCmdTrg = getTestPath("result")
		packCmdTrg = getTestPath("result")
		buildModuleCmdTrg = getTestPath("result")
		cleanupCmdTrg = getTestPath("result")
		logs.Logger = logs.NewLogger()
		err := os.Mkdir(mtadCmdTrg, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(mtadCmdTrg)
		Ω(err).Should(Succeed())
	})

	var _ = Describe("Pack and cleanup commands", func() {
		It("Target file in opened status", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			// Target path has to be dir, but is currently created as file
			packCmdModule = "ui5app"
			packCmdSrc = getTestPath("mtahtml5")
			packCmdPlatform = "cf"
			ep := dir.Loc{SourcePath: packCmdSrc, TargetPath: packCmdTrg}
			createDirInTmpFolder("mtahtml5")
			createFileInTmpFolder("mtahtml5", "ui5app")
			Ω(packModuleCmd.RunE(nil, []string{})).Should(HaveOccurred())
			Ω(str.String()).Should(ContainSubstring(fmt.Sprintf(artifacts.PackFailedOnArchMsg, "ui5app")))
			// cleanup command used for test temp file removal
			cleanupCmdSrc = packCmdSrc
			Ω(cleanupCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
		})

	})

	var _ = Describe("Pack", func() {
		DescribeTable("Standard cases", func(projectPath string, match types.GomegaMatcher) {
			packCmdSrc = projectPath
			packCmdModule = "ui5app"
			Ω(packModuleCmd.RunE(nil, []string{})).Should(match)
		},
			Entry("SanityTest", getTestPath("mtahtml5"), Succeed()),
			Entry("Wrong path to project", getTestPath("mtahtml6"), HaveOccurred()),
		)
	})

	var _ = Describe("Build", func() {
		var config []byte

		BeforeEach(func() {
			config = make([]byte, len(commands.ModuleTypeConfig))
			copy(config, commands.ModuleTypeConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			commands.ModuleTypeConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  type:
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  type:
`)
		})

		AfterEach(func() {
			err := os.RemoveAll(getTestPath("result"))
			Ω(err).Should(Succeed())
			commands.ModuleTypeConfig = make([]byte, len(config))
			copy(commands.ModuleTypeConfig, config)
		})

		It("build Command", func() {
			buildModuleCmdModule = "node-js"
			buildModuleCmdSrc = getTestPath("mta")
			buildModuleCmdPlatform = "cf"
			ep := dir.Loc{SourcePath: buildModuleCmdSrc, TargetPath: buildModuleCmdTrg}
			Ω(buildModuleCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath(buildModuleCmdModule)).Should(BeAnExistingFile())
		})
	})
})
