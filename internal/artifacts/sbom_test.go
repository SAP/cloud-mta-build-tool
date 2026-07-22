package artifacts

import (
	"os"
	"path/filepath"

	cdx "github.com/CycloneDX/cyclonedx-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testdataMta          = "testdata/mta"
	genSBomMergedBomPath = "gen-sbom-result/merged.bom.xml"
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
		source := testdataMta
		sbomFilePath := ""
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(filepath.Join(getTestPath("mta"), "mta.bom.xml"))).Should(Succeed())
	})
	It("Success - sbom-gen with abs source and relative sbom-file-path paramerter", func() {
		source := getTestPath("mta")
		sbomFilePath := genSBomMergedBomPath
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
		source := testdataMta
		sbomFilePath := genSBomMergedBomPath
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - sbom-gen with relative source and abs sbom-file-path paramerter", func() {
		source := testdataMta
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
		source := testdataMta
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "result_without_suffix")
		Ω(ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with json suffix sbom-file-name parameter", func() {
		source := testdataMta
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "result.json")

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .json is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with unknow suffix sbom-file-name parameter", func() {
		source := testdataMta
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "result.unknow")

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .unknow is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	/* It("Failure - sbom-gen with invalid sbom-file-path paramerter case 1", func() {
		source := testdataMta
		sbomFilePath := "gen-sbom-result>>?</merged.bom.xml"

		err := ExecuteProjectSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - sbom-gen with invalid sbom-file-path paramerter case 2", func() {
		source := testdataMta
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
		source := testdataMta
		sbomFilePath := genSBomMergedBomPath
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build with abs source and relative sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := genSBomMergedBomPath
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("mta", "gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build with relatvie source and abs sbom-file-path parameter", func() {
		source := testdataMta
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build with abs source and abs sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "merged.bom.xml")
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build without sbom-file-path parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := ""
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(Succeed())
	})
	It("Failure - build with invalid source paramerter case 1", func() {
		source := "testdata??/mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")

		err := ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with invalid source paramerter case 2", func() {
		source := "testdata/*??>mta"
		sbomFilePath := filepath.Join(getTestPath("gen-sbom-result"), "merged.bom.xml")

		err := ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Success - build without suffix sbom-file-name parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "result_without_suffix")
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(Succeed())
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with json suffix sbom-file-name parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "result.json")

		err := ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .json is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with unknow suffix sbom-file-name parameter", func() {
		source := getTestPath("mta")
		sbomFilePath := getTestPath("gen-sbom-result", "result.unknow")

		err := ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("sbom file type .unknow is not supported at present"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	/* It("Failure - build with invalid sbom-file-path paramerter case 1", func() {
		source := testdataMta
		sbomFilePath := "gen-sbom-result>>?</merged.bom.xml"

		err := ExecuteProjectBuildeSBomGenerate(source, sbomFilePath, os.Getwd)
		Ω(err).Should(HaveOccurred())
		//Ω(err.Error()).Should(ContainSubstring("The filename, directory name, or volume label syntax is incorrect"))
		Ω(os.RemoveAll(getTestPath("gen-sbom-result"))).Should(Succeed())
	})
	It("Failure - build with invalid sbom-file-path paramerter case 2", func() {
		source := testdataMta
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
		Ω(ExecuteProjectBuildeSBomGenerate(source, "", sbomFilePath, os.Getwd)).Should(HaveOccurred())
		Ω(os.RemoveAll(tmpSrcFolder)).Should(Succeed())
	})
})

var _ = Describe("deduplicateBOM", func() {
	ptr := func(s []string) *[]string { return &s }
	ptrComp := func(c []cdx.Component) *[]cdx.Component { return &c }
	ptrDep := func(d []cdx.Dependency) *[]cdx.Dependency { return &d }

	It("does nothing when BOM has no components or dependencies", func() {
		bom := &cdx.BOM{}
		deduplicateBOM(bom)
		Ω(bom.Components).Should(BeNil())
		Ω(bom.Dependencies).Should(BeNil())
	})

	It("keeps unique components unchanged", func() {
		bom := &cdx.BOM{
			Components: ptrComp([]cdx.Component{
				{BOMRef: "pkg:npm/a@1.0.0"},
				{BOMRef: "pkg:npm/b@2.0.0"},
			}),
		}
		deduplicateBOM(bom)
		Ω(*bom.Components).Should(HaveLen(2))
	})

	It("removes duplicate components with identical BOMRef", func() {
		bom := &cdx.BOM{
			Components: ptrComp([]cdx.Component{
				{BOMRef: "pkg:npm/lodash@4.17.21"},
				{BOMRef: "pkg:npm/lodash@4.17.21"},
				{BOMRef: "pkg:npm/express@4.18.0"},
			}),
		}
		deduplicateBOM(bom)
		Ω(*bom.Components).Should(HaveLen(2))
		Ω((*bom.Components)[0].BOMRef).Should(Equal("pkg:npm/lodash@4.17.21"))
		Ω((*bom.Components)[1].BOMRef).Should(Equal("pkg:npm/express@4.18.0"))
	})

	It("removes duplicate components using PackageURL as fallback when BOMRef is empty", func() {
		bom := &cdx.BOM{
			Components: ptrComp([]cdx.Component{
				{PackageURL: "pkg:npm/lodash@4.17.21"},
				{PackageURL: "pkg:npm/lodash@4.17.21"},
			}),
		}
		deduplicateBOM(bom)
		Ω(*bom.Components).Should(HaveLen(1))
	})

	It("keeps unique dependencies unchanged", func() {
		bom := &cdx.BOM{
			Dependencies: ptrDep([]cdx.Dependency{
				{Ref: "pkg:npm/a@1.0.0", Dependencies: ptr([]string{"pkg:npm/b@1.0.0"})},
				{Ref: "pkg:npm/b@1.0.0", Dependencies: ptr([]string{})},
			}),
		}
		deduplicateBOM(bom)
		Ω(*bom.Dependencies).Should(HaveLen(2))
	})

	It("removes duplicate dependency entries", func() {
		bom := &cdx.BOM{
			Dependencies: ptrDep([]cdx.Dependency{
				{Ref: "pkg:npm/lodash@4.17.21", Dependencies: ptr([]string{})},
				{Ref: "pkg:npm/express@4.18.0", Dependencies: ptr([]string{})},
				{Ref: "pkg:npm/lodash@4.17.21", Dependencies: ptr([]string{})},
			}),
		}
		deduplicateBOM(bom)
		Ω(*bom.Dependencies).Should(HaveLen(2))
		Ω((*bom.Dependencies)[0].Ref).Should(Equal("pkg:npm/lodash@4.17.21"))
		Ω((*bom.Dependencies)[1].Ref).Should(Equal("pkg:npm/express@4.18.0"))
	})

	It("removes duplicate entries within a dependsOn list", func() {
		bom := &cdx.BOM{
			Dependencies: ptrDep([]cdx.Dependency{
				{
					Ref: "pkg:npm/a@1.0.0",
					Dependencies: ptr([]string{
						"pkg:npm/lodash@4.17.21",
						"pkg:npm/lodash@4.17.21",
						"pkg:npm/express@4.18.0",
					}),
				},
			}),
		}
		deduplicateBOM(bom)
		deps := *bom.Dependencies
		Ω(deps).Should(HaveLen(1))
		Ω(*deps[0].Dependencies).Should(HaveLen(2))
		Ω(*deps[0].Dependencies).Should(ContainElements("pkg:npm/lodash@4.17.21", "pkg:npm/express@4.18.0"))
	})

	It("handles nil dependsOn list within a dependency", func() {
		bom := &cdx.BOM{
			Dependencies: ptrDep([]cdx.Dependency{
				{Ref: "pkg:npm/a@1.0.0", Dependencies: nil},
			}),
		}
		Ω(func() { deduplicateBOM(bom) }).ShouldNot(Panic())
		Ω((*bom.Dependencies)[0].Dependencies).Should(BeNil())
	})

	It("deduplicates both components and dependencies simultaneously", func() {
		bom := &cdx.BOM{
			Components: ptrComp([]cdx.Component{
				{BOMRef: "pkg:npm/shared@1.0.0"},
				{BOMRef: "pkg:npm/shared@1.0.0"},
				{BOMRef: "pkg:npm/unique@2.0.0"},
			}),
			Dependencies: ptrDep([]cdx.Dependency{
				{Ref: "pkg:npm/shared@1.0.0", Dependencies: ptr([]string{"pkg:npm/dep@1.0.0", "pkg:npm/dep@1.0.0"})},
				{Ref: "pkg:npm/shared@1.0.0", Dependencies: ptr([]string{"pkg:npm/dep@1.0.0"})},
				{Ref: "pkg:npm/unique@2.0.0", Dependencies: ptr([]string{})},
			}),
		}
		deduplicateBOM(bom)
		Ω(*bom.Components).Should(HaveLen(2))
		Ω(*bom.Dependencies).Should(HaveLen(2))
		Ω(*(*bom.Dependencies)[0].Dependencies).Should(HaveLen(1))
	})
})
