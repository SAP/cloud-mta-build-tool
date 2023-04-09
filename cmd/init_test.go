package commands

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {

	BeforeEach(func() {
		Ω(os.MkdirAll(getTestPath("result"), os.ModePerm)).Should(Succeed())
	})
	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})
	It("Sanity", func() {
		initCmdSrc = getTestPath("mta")
		initCmdTrg = getTestPath("result")
		initCmd.Run(nil, []string{})
		Ω(getTestPath("result", "Makefile.mta")).Should(BeAnExistingFile())
	})
})

var _ = Describe("Build", func() {
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
	/* It("Success - build with abs source parameter", func() {
		buildCmdSrc = getTestPath("mta")
		source := "\"" + buildCmdSrc + "\""
		target := filepath.Join(buildCmdSrc, dir.MtarFolder)

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(target)).Should(Succeed())
	})
	It("Success - build with relative source parameter", func() {
		// Notice: relative source path is relative to os.Getwd()
		buildCmdSrc := "testdata/mta"

		source := "\"" + buildCmdSrc + "\""
		target := filepath.Join(buildCmdSrc, dir.MtarFolder)

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(target)).Should(Succeed())
	})
	It("Success - build with source and abs target parameter", func() {
		buildCmdSrc = getTestPath("mta")
		buildCmdTrg = getTestPath("result_for_sbom")
		// Ω(os.MkdirAll(buildCmdTrg, os.ModePerm)).Should(Succeed())

		source := "\"" + buildCmdSrc + "\""
		target := "\"" + buildCmdTrg + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(buildCmdTrg)).Should(Succeed())
	})
	It("Success - build with source and relative target path parameter", func() {
		// Notice: target parameter is relative to source parameter
		buildCmdSrc = getTestPath("mta")
		buildCmdTrg = "testdata/result_for_sbom"
		targetPath := filepath.Join(getTestPath("mta"), buildCmdTrg)

		source := "\"" + buildCmdSrc + "\""
		target := "\"" + buildCmdTrg + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(targetPath)).Should(Succeed())
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
	It("Failure - build without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+tmpSrcFolder)

		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
	It("Failure - build with invalidate platform parameter", func() {
		buildCmdSrc = getTestPath("mta")
		buildCmdPlatform = "xxx"
		source := "\"" + buildCmdSrc + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --platform "+buildCmdPlatform)

		Ω(cmd.Run()).Should(HaveOccurred())
	}) */
})
