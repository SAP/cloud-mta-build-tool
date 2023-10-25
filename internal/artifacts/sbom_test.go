package artifacts

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mbt sbom-gen command", func() {
	BeforeEach(func() {

	})
	AfterEach(func() {

	})

	It("Success - sbom-gen with abs source and without sbom-file-path paramerter", func() {
		source := getTestPath("mta")
		sbomFilePath := ""
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(getTestPath("mta"), "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and without sbom-file-path paramerter", func() {
		source := "testdata/mta"
		sbomFilePath := ""
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(getTestPath("mta"), "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with abs source and relative sbom-file-path paramerter", func() {
		source := getTestPath("mta")
		sbomFilePath := "gen-sbom-result/merged.bom.xml"
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(getTestPath("mta", "gen-sbom-result")))).Should(Succeed())

	})
	It("Success - sbom-gen with abs source and abs sbom-file-path paramerter", func() {
		source := getTestPath("mta")
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(getTestPath("gen-sbom-result")))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and relative sbom-file-path paramerter", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result/merged.bom.xml"
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and abs sbom-file-path paramerter", func() {
		source := "testdata/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalid source paramerter case 1", func() {
		source := "testdata??/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())

	})
	It("Failure - sbom-gen with invalid source paramerter case 2", func() {
		source := "testdata/*??>mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - sbom-gen without suffix sbom-file-name paramerter", func() {
		source := "testdata/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "result_without_suffix")
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with json suffix sbom-file-name parameter", func() {
		source := "testdata/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "result.json")

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .json is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with unknow suffix sbom-file-name parameter", func() {
		source := "testdata/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "result.unknow")

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .unknow is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	/* It("Failure - sbom-gen with invalid sbom-file-path paramerter case 1", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result>>?</merged.bom.xml"

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalid sbom-file-path paramerter case 2", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result/<<*merged.bom.xml"

		// Notice: the merge sbom file name is invalid, the error will raised from cyclondx-cli merge command
		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	}) */
	It("Failure - sbom-gen without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())
		source := tmpSrcFolder
		sbomFolderName := getTestPath("gen-sbom-result")
		sbomFileName := "merged.bom.xml"
		sbomFilePath := filepath.Join(sbomFolderName, sbomFileName)

		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
})

var _ = Describe("mbt build with sbom gen command", func() {
	BeforeEach(func() {
	})
	AfterEach(func() {
	})
	It("Success - build with relatvie source and relative sbom-file-path parameter", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result/merged.bom.xml"
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build with abs source and relative sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := "gen-sbom-result/merged.bom.xml"
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build with relatvie source and abs sbom-file-path parameter", func() {
		source := "testdata/mta"
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build with abs source and abs sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build without sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := ""
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
	})
	It("Failure - build with invalid source paramerter case 1", func() {
		source := "testdata??/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")

		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with invalid source paramerter case 2", func() {
		source := "testdata/*??>mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")

		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build without suffix sbom-file-name parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "result_without_suffix")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with json suffix sbom-file-name parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "result.json")

		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .json is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with unknow suffix sbom-file-name parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "result.unknow")

		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .unknow is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	/* It("Failure - build with invalid sbom-file-path paramerter case 1", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result>>?</merged.bom.xml"

		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with invalid sbom-file-path paramerter case 2", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result/<<*merged.bom.xml"

		// Notice: the merge sbom file name is invalid, the error will raised from cyclondx-cli merge command
		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	}) */
	It("Failure - build without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())

		source := tmpSrcFolder
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
})
