package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"cloud-mta-build-tool/internal/commands"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
)

var _ = Describe("Commands", func() {

	BeforeEach(func() {
		mtadCmdTrg = getTestPath("result")
		metaCmdTrg = getTestPath("result")
		mtarCmdTrg = getTestPath("result")
		packCmdTrg = getTestPath("result")
		buildCmdTrg = getTestPath("result")
		cleanupCmdTrg = getTestPath("result")
		logs.Logger = logs.NewLogger()
		err := os.Mkdir(mtadCmdTrg, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		os.RemoveAll(mtadCmdTrg)
	})

	var _ = Describe("Pack and cleanup commands", func() {
		It("Target file in opened status", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			// Target path has to be dir, but is currently created and opened as file
			packCmdModule = "ui5app"
			packCmdSrc = getTestPath("mtahtml5")
			packCmdPlatform = "cf"
			ep := dir.Loc{SourcePath: packCmdSrc, TargetPath: packCmdTrg}
			targetTmpDir := ep.GetTargetTmpDir()
			err := os.MkdirAll(targetTmpDir, os.ModePerm)
			if err != nil {
				logs.Logger.Error(err)
			}
			f, _ := os.Create(filepath.Join(targetTmpDir, "ui5app"))
			Ω(packModuleCmd.RunE(nil, []string{})).Should(HaveOccurred())
			fmt.Println(str.String())
			Ω(str.String()).Should(ContainSubstring("the packing of the ui5app module failed when creating"))

			f.Close()
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
			config = make([]byte, len(commands.CommandsConfig))
			copy(config, commands.CommandsConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			commands.CommandsConfig = []byte(`
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
			os.RemoveAll(getTestPath("result"))
			commands.CommandsConfig = make([]byte, len(config))
			copy(commands.CommandsConfig, config)
		})

		It("build Command", func() {
			buildCmdModule = "node-js"
			buildCmdSrc = getTestPath("mta")
			ep := dir.Loc{SourcePath: buildCmdSrc, TargetPath: buildCmdTrg}
			Ω(buildModuleCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath(buildCmdModule)).Should(BeAnExistingFile())
		})
	})
})
