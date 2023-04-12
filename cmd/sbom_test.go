package commands

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mbt cli build and sbom gen", func() {
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
	It("Success - build and gen sbom with relatvie source and relative sbom-file-path parameter", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - build and gen sbom with abs source and relative sbom-file-path parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - build and gen sbom with relatvie source and abs sbom-file-path parameter", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - build and gen sbom with abs source and abs sbom-file-path parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - build and gen sbom with abs source and relative sbom-file-path paramerter with sbom file under project root", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "merged.bom.xml"))).Should(Succeed())
	})
	It("Success - build and gen sbom with relative source and relative sbom-file-path paramerter with sbom file under project root", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + "merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "merged.bom.xml"))).Should(Succeed())
	})
	It("Success - build and gen sbom without sbom-file-path parameter", func() {
		source := "\"" + getTestPath("mta") + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
	})
	It("Failure - build and gen sbom without mta.yaml", func() {
		source := "\"" + getTestPath("tmp") + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""
		Ω(os.MkdirAll(getTestPath("tmp"), os.ModePerm)).Should(Succeed())

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("tmp"))).Should(Succeed())
	})
	It("Failure - build and gen sbom with invalid sbom-file-path parameter case 1", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "sbom-gen-result>>?/merged.bom.xml" + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - build and gen sbom with invalid sbom-file-path parameter case 2", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "sbom-gen-result/**??merged.bom.xml" + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		// Notice: the merge sbom file name is invalidate, the error will raised from cyclondx-cli merge command
		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("mta", dir.MtarFolder))).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
})

var _ = Describe("mbt cli sbom-gen", func() {
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

	It("Success - sbom-gen with abs source and without sbom-file-path paramerter", func() {
		source := "\"" + getTestPath("mta") + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and without sbom-file-path paramerter", func() {
		source := "\"" + "testdata/mta" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with abs source and relative sbom-file-path paramerter", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with abs source and relative sbom-file-path paramerter with sbom file under project root", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + "merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "merged.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with abs source and abs sbom-file-path paramerter", func() {
		source := "\"" + getTestPath("mta") + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and relative sbom-file-path paramerter", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and relative sbom-file-path paramerter with sbom file under project root", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + "merged.bom.xml" + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "merged.bom.xml"))).Should(Succeed())
	})

	It("Success - sbom-gen with relative source and abs sbom-file-path paramerter", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen without mta.yaml", func() {
		source := "\"" + getTestPath("tmp") + "\""
		sbom_file_path := "\"" + "sbom-gen-result/merged.bom.xml" + "\""
		Ω(os.MkdirAll(getTestPath("tmp"), os.ModePerm)).Should(Succeed())

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)

		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("tmp"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate source paramerter case 1", func() {
		source := "\"" + "testdata??>/mta" + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)
		cmd.Stdout = &stdout
		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate source paramerter case 2", func() {
		source := "\"" + "testdata/***mta" + "\""
		sbom_file_path := "\"" + getTestPath("mta", "sbom-gen-result", "merged.bom.xml") + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)
		cmd.Stdout = &stdout
		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate sbom-file-path paramerter case 1", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + "sbom-gen-result??/merged.bom.xml" + "\""
		var stdout bytes.Buffer

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)
		cmd.Stdout = &stdout
		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(stdout.String()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate sbom-file-path paramerter case 2", func() {
		source := "\"" + "testdata/mta" + "\""
		sbom_file_path := "\"" + "sbom-gen-result/>>>merged.bom.xml" + "\""

		cmd := exec.Command("bash", "-c", mbtCmdCLI+" sbom-gen"+" --source "+source+" --sbom-file-path "+sbom_file_path)
		// Notice: the merge sbom file name is invalidate, the error will raised from cyclondx-cli merge command
		Ω(cmd.Run()).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
})

var _ = Describe("project sbom gen command", func() {
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
		projectSBomGenCmdSBOMPath = "sbom-gen-result/merged.bom.xml"
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success -sbom-gen with abs source and abs sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = getTestPath("mta")
		projectSBomGenCmdSBOMPath = filepath.Join(getTestPath("sbom-gen-result"), "merged.bom.xml")
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("sbom-gen-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and relative sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = "sbom-gen-result/merged.bom.xml"
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "sbom-gen-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and abs sbom-file-path paramerter", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = filepath.Join(getTestPath("sbom-gen-result"), "merged.bom.xml")
		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())
		projectSBomGenCmdSrc = tmpSrcFolder
		sbomFolderName := getTestPath("sbom-gen-result")
		sbomFileName := "merged.bom.xml"
		projectSBomGenCmdSBOMPath = filepath.Join(sbomFolderName, sbomFileName)

		Ω(projectSBomGenCommand.RunE(nil, []string{})).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate source paramerter case 1", func() {
		projectSBomGenCmdSrc = "testdata/>><>mta"
		projectSBomGenCmdSBOMPath = filepath.Join(getTestPath("sbom-gen-result"), "merged.bom.xml")

		err := projectSBomGenCommand.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate source paramerter case 2", func() {
		projectSBomGenCmdSrc = "testdata??/mta"
		projectSBomGenCmdSBOMPath = filepath.Join(getTestPath("sbom-gen-result"), "merged.bom.xml")

		err := projectSBomGenCommand.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate sbom-file-path paramerter case 1", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = "sbom-gen-result>>/merged.bom.xml"

		err := projectSBomGenCommand.RunE(nil, []string{})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("sbom-gen-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalidate sbom-file-path paramerter case 2", func() {
		projectSBomGenCmdSrc = "testdata/mta"
		projectSBomGenCmdSBOMPath = "sbom-gen-result/???merged.bom.xml"
		err := projectSBomGenCommand.RunE(nil, []string{})
		// Notice: the merge sbom file name is invalidate, the error will raised from cyclondx-cli merge command
		Ω(err).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("sbom-gen-result"))).Should(Succeed())
	})
})
