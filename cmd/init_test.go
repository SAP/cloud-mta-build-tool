package commands

import (
	"bytes"
	"os"
	"os/exec"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
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
	})
	It("Success - build with abs source parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
	})
	It("Success - build with relative source parameter", func() {
		// Notice: relative source path is relative to os.Getwd()
		source := "\"" + "testdata/mta" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
	})
	It("Success - build with abs source and abs target parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		target := "\"" + getTestPath("mtar_result") + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mtar_result"))).Should(Succeed())
	})
	It("Success - build with abs source and relative target path parameter", func() {
		// Notice: target parameter is relative to source parameter
		source := "\"" + getTestPath("mta") + "\""
		target := "\"" + "mtar_result" + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	})
	It("Success - build with relative source and abs target path parameter", func() {
		source := "\"" + "testdata/mta" + "\""
		target := "\"" + getTestPath("mta", "mtar_result") + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	})
	It("Success - build with relative source and relative target path parameter", func() {
		// Notice: target parameter is relative to source parameter
		source := "\"" + "testdata/mta" + "\""
		target := "\"" + "mtar_result" + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	})
	It("Failure - build with invalid source parameter case 1", func() {
		// Notice: target parameter is relative to source parameter
		source := "\"" + "testdata/**??mta" + "\""
		target := "\"" + "mtar_result" + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(HaveOccurred())
		// Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	})
	It("Failure - build with invalid source parameter case 2", func() {
		// Notice: target parameter is relative to source parameter
		source := "\"" + "testdata<></mta" + "\""
		target := "\"" + "mtar_result" + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(HaveOccurred())
		// Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	})
	/* It("Failure - build with invalid target parameter case 1", func() {
		// Notice: target parameter is relative to source parameter
		source := "\"" + "testdata/mta" + "\""
		target := "\"" + "mtar_result<>/tmp" + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(HaveOccurred())
		//Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	})
	It("Failure - build with invalid target parameter case 2", func() {
		// Notice: target parameter is relative to source parameter
		source := "\"" + "testdata/mta" + "\""
		target := "\"" + "mtar_result/??*tmp" + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --target "+target)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(HaveOccurred())
		//Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "mtar_result"))).Should(Succeed())
	}) */
	It("Success - build with relative sbom-file-path parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - build with abs sbom-file-path parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - build without mta.yaml", func() {
		source := "\"" + getTestPath("tmp") + "\""
		Ω(os.MkdirAll(getTestPath("tmp"), os.ModePerm)).Should(Succeed())
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)

		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("tmp"))).Should(Succeed())
	})
	It("Failure - build with invalid platform parameter", func() {
		platform := "\"" + "xxx" + "\""
		source := "\"" + getTestPath("mta") + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --platform "+platform)

		Ω(cmd.Run()).Should(HaveOccurred())
	})
})
