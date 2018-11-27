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

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
)

var _ = Describe("Commands", func() {

	BeforeEach(func() {
		targetMtadFlag = getTestPath("result")
		targetMetaFlag = getTestPath("result")
		targetMtarFlag = getTestPath("result")
		targetPackFlag = getTestPath("result")
		targetBModuleFlag = getTestPath("result")
		targetCleanupFlag = getTestPath("result")
		logs.Logger = logs.NewLogger()
		err := os.Mkdir(targetMtadFlag, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		os.RemoveAll(targetMtadFlag)
	})

	var _ = Describe("Pack and cleanup commands", func() {
		It("Target file in opened status", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			// Target path has to be dir, but is currently created and opened as file
			pPackModuleFlag = "ui5app"
			sourcePackFlag = getTestPath("mtahtml5")
			ep := dir.Loc{SourcePath: sourcePackFlag, TargetPath: targetPackFlag}
			targetTmpDir, _ := ep.GetTargetTmpDir()
			err := os.MkdirAll(targetTmpDir, os.ModePerm)
			if err != nil {
				logs.Logger.Error(err)
			}
			f, _ := os.Create(filepath.Join(targetTmpDir, "ui5app"))
			Ω(packCmd.RunE(nil, []string{})).Should(HaveOccurred())
			fmt.Println(str.String())
			Ω(str.String()).Should(ContainSubstring("Pack of module ui5app failed on making directory"))

			f.Close()
			// cleanup command used for test temp file removal
			sourceCleanupFlag = sourcePackFlag
			Ω(cleanupCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
		})

	})

	var _ = Describe("Generate commands call", func() {

		var ep dir.Loc

		It("Generate Meta", func() {
			sourceMetaFlag = getTestPath("mtahtml5")
			ep = dir.Loc{SourcePath: sourceMetaFlag, TargetPath: targetMetaFlag}
			Ω(genMetaCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})
		It("Generate Mtad", func() {
			sourceMtadFlag = getTestPath("mtahtml5")
			ep = dir.Loc{SourcePath: sourceMtadFlag, TargetPath: targetMtadFlag}
			err := genMtadCmd.RunE(nil, []string{})
			if err != nil {
				fmt.Println(err)
			}
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})
		It("Generate Mtar", func() {
			sourceMtarFlag = getTestPath("mtahtml5")
			ep = dir.Loc{SourcePath: sourceMtarFlag, TargetPath: targetMtarFlag}
			Ω(genMetaCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(genMtarCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(getTestPath("result", "mtahtml5.mtar")).Should(BeAnExistingFile())
		})
	})

	var _ = Describe("Validate", func() {
		It("Invalid yaml path", func() {
			var str bytes.Buffer
			sourceValidateFlag = getTestPath("mta1")
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			validateCmd.RunE(nil, []string{})

			Ω(str.String()).Should(ContainSubstring("Error reading the MTA file"))
		})
	})

	var _ = Describe("Pack", func() {
		DescribeTable("Standard cases", func(projectPath string, match types.GomegaMatcher) {
			sourcePackFlag = projectPath
			pPackModuleFlag = "ui5app"
			Ω(packCmd.RunE(nil, []string{})).Should(match)
		},
			Entry("SanityTest", getTestPath("mtahtml5"), Succeed()),
			Entry("Wrong path to project", getTestPath("mtahtml6"), HaveOccurred()),
		)
	})

	var _ = Describe("Build", func() {
		var config []byte

		BeforeEach(func() {
			config = make([]byte, len(builders.CommandsConfig))
			copy(config, builders.CommandsConfig)
			// Simplified commands configuration (performance purposes). removed "npm prune --production"
			builders.CommandsConfig = []byte(`
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
			builders.CommandsConfig = make([]byte, len(config))
			copy(builders.CommandsConfig, config)
		})

		It("build Command", func() {
			pBuildModuleNameFlag = "node-js"
			sourceBModuleFlag = getTestPath("mta")
			ep := dir.Loc{SourcePath: sourceBModuleFlag, TargetPath: targetBModuleFlag}
			Ω(bModuleCmd.RunE(nil, []string{})).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath(pBuildModuleNameFlag)).Should(BeAnExistingFile())
		})
	})

})

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}
