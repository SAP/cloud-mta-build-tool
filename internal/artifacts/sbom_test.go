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
	It("Success - gen sbom with relatvie source and relative sbom-file-path parameter", func() {
		source := "testdata/mta"
		sbomFilePath := "gen-sbom-result/merged.bom.xml"
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - gen sbom with abs source and relative sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := "gen-sbom-result/merged.bom.xml"
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - gen sbom with relatvie source and abs sbom-file-path parameter", func() {
		source := "testdata/mta"
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - gen sbom with abs source and abs sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - gen sbom without sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := ""
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
	})
	It("Failure - gen sbom without mta.yaml", func() {
		tmpSrcFolder := getTestPath("tmp")
		Ω(os.MkdirAll(tmpSrcFolder, os.ModePerm)).Should(Succeed())

		source := tmpSrcFolder
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
})

var _ = Describe("validate path", func() {
	BeforeEach(func() {
	})
	AfterEach(func() {
	})

	/* It("Success - validate source and sbomFilePath", func() {
		source := "/c/windows/"
		sbomFilePath := "./test.txt"
		Ω(validatePath(source, sbomFilePath)).Should(Succeed())
	})

	It("Failure - invalidate source and sbomFilePath", func() {
		source := "  /??  "
		sbomFilePath := "./test.txt"
		Ω(validatePath(source, sbomFilePath)).Should(Succeed())
	}) */
})
