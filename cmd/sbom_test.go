package commands

import (
	"os"
	"os/exec"
	"path/filepath"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build and SBom Gen", func() {
	BeforeEach(func() {
		mbtCmdCLI = getBuildCmdCli()
	})

	AfterEach(func() {
		mbtCmdCLI = ""
		buildCmdSrc = ""
		buildCmdTrg = ""
		buildCmdPlatform = ""
		buildCmdKeepMakefile = false
	})

	It("Success - build with relative sbom-file-path parameter", func() {
		buildCmdSrc = getTestPath("mta")
		sbomFolderName := "sbom-gen-result"
		sbomFileName := "merged.bom.xml"
		sbomTarget := filepath.Join(getTestPath("mta"), sbomFolderName)
		buildCmdSBomFilePath := sbomFolderName + "/" + sbomFileName

		source := "\"" + buildCmdSrc + "\""
		target := filepath.Join(buildCmdSrc, dir.MtarFolder)
		sbomfilepath := "\"" + buildCmdSBomFilePath + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbomfilepath)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(target)).Should(Succeed())
		Ω(os.RemoveAll(sbomTarget)).Should(Succeed())
	})
	It("Success - build with abs sbom-file-path parameter", func() {
		buildCmdSrc = getTestPath("mta")
		sbomPath := getTestPath("sbom-gen-result")
		sbomFileName := "merged.bom.xml"
		buildCmdSBomFilePath := sbomPath + "/" + sbomFileName

		target := filepath.Join(buildCmdSrc, dir.MtarFolder)
		source := "\"" + buildCmdSrc + "\""
		sbomfilepath := "\"" + buildCmdSBomFilePath + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbomfilepath)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(target)).Should(Succeed())
		Ω(os.RemoveAll(sbomPath)).Should(Succeed())
	})
})

var _ = Describe("Gen SBom", func() {
	BeforeEach(func() {
	})
	AfterEach(func() {
		projectBuildSBomGenCmdSrc = ""
		projectBuildSBomGenCmdSBOMPath = ""
	})
	It("Success - gen sbom with relatvie source and relative sbom-file-path parameter", func() {
		projectBuildSBomGenCmdSrc = "testdata/mta"
		sbomFolderName := "Gen-SBom-Result"
		sbomFileName := "merged.bom.xml"
		sbomTarget := filepath.Join(getTestPath("mta"), sbomFolderName)
		projectBuildSBomGenCmdSBOMPath = sbomFolderName + "/" + sbomFileName

		Ω(projectBuildSBomGenCommand.RunE(nil, []string{})).Should(Succeed())

		Ω(os.RemoveAll(sbomTarget)).Should(Succeed())
	})
	It("Success - gen sbom with abs source and relative sbom-file-path parameter", func() {
		projectBuildSBomGenCmdSrc = getTestPath("mta")
		sbomFolderName := "Gen-SBom-Result"
		sbomFileName := "merged.bom.xml"
		sbomTarget := filepath.Join(getTestPath("mta"), sbomFolderName)
		projectBuildSBomGenCmdSBOMPath = sbomFolderName + "/" + sbomFileName

		Ω(projectBuildSBomGenCommand.RunE(nil, []string{})).Should(Succeed())

		Ω(os.RemoveAll(sbomTarget)).Should(Succeed())
	})
	It("Success - gen sbom with relatvie source and abs sbom-file-path parameter", func() {
		projectBuildSBomGenCmdSrc = "testdata/mta"
		sbomFolderName := getTestPath("Gen-SBom-Result")
		sbomFileName := "merged.bom.xml"
		projectBuildSBomGenCmdSBOMPath = filepath.Join(sbomFolderName, sbomFileName)

		Ω(projectBuildSBomGenCommand.RunE(nil, []string{})).Should(Succeed())

		Ω(os.RemoveAll(sbomFolderName)).Should(Succeed())
	})
	It("Success - gen sbom with abs source and abs sbom-file-path parameter", func() {
		projectBuildSBomGenCmdSrc = getTestPath("mta")
		sbomFolderName := getTestPath("Gen-SBom-Result")
		sbomFileName := "merged.bom.xml"
		projectBuildSBomGenCmdSBOMPath = filepath.Join(sbomFolderName, sbomFileName)

		Ω(projectBuildSBomGenCommand.RunE(nil, []string{})).Should(Succeed())

		Ω(os.RemoveAll(sbomFolderName)).Should(Succeed())
	})
	It("Success - gen sbom without sbom-file-path parameter", func() {
		projectBuildSBomGenCmdSrc = getTestPath("mta")
		projectBuildSBomGenCmdSBOMPath = ""

		Ω(projectBuildSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
	})
	It("Failure - gen sbom without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())
		projectBuildSBomGenCmdSrc = tmpSrcFolder
		sbomFolderName := getTestPath("Gen-SBom-Result")
		sbomFileName := "merged.bom.xml"
		projectBuildSBomGenCmdSBOMPath = filepath.Join(sbomFolderName, sbomFileName)

		Ω(projectBuildSBomGenCommand.RunE(nil, []string{})).Should(HaveOccurred())

		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
})

var _ = Describe("SBom Gen", func() {
	BeforeEach(func() {
	})

	AfterEach(func() {
		projectSBomGenCmdSrc = ""
		projectSBomGenCmdSBOMPath = ""
	})
	It("Success - sbom-gen with abs source and without sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = getTestPath("mta")
		projectSBomGenCmdSBOMPath = ""
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(projectSBomGenCmdSrc, "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and without sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = ""
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(getTestPath("mta"), "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with abs source and relative sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = getTestPath("mta")
		projectSBomGenCmdSBOMPath = "Gen-SBom-Result/merged.bom.xml"
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "Gen-SBom-Result"))).Should(Succeed())
	})
	It("Success -sbom-gen with abs source and abs sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = getTestPath("mta")
		projectSBomGenCmdSBOMPath = filepath.Join(getTestPath("Gen-SBom-Result"), "merged.bom.xml")
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("Gen-SBom-Result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and relative sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = "Gen-SBom-Result/merged.bom.xml"
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "Gen-SBom-Result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and abs sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = filepath.Join(getTestPath("Gen-SBom-Result"), "merged.bom.xml")
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("Gen-SBom-Result"))).Should(Succeed())
	})
	It("Failure - sbom-gen without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())
		projectSBomGenCmdSrc = tmpSrcFolder
		sbomFolderName := getTestPath("Gen-SBom-Result")
		sbomFileName := "merged.bom.xml"
		projectSBomGenCmdSBOMPath = filepath.Join(sbomFolderName, sbomFileName)

		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
})
