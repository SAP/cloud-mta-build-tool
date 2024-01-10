package dir

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModuleLocation", func() {
	It("GetTarget", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("test"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		loc := ModuleLocation(projectLoc, false)
		Ω(loc.GetTarget()).Should(Equal(getPath("test")))
	})

	It("GetTargetTmpRoot, target path calculated", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata", ".test_mta_build_tmp", "module"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc, false)
		Ω(moduleLoc.GetTargetTmpRoot()).Should(Equal(getPath("testdata", ".test_mta_build_tmp")))
	})

	It("GetTargetTmpRoot, target path defined", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc, true)
		Ω(moduleLoc.GetTargetTmpRoot()).Should(Equal(getPath("testdata")))
	})

	It("GetSourceModuleDir", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc, false)
		Ω(moduleLoc.GetSourceModuleDir("path1")).Should(Equal((getPath("testdata", "path1"))))
	})

	It("GetSourceModuleArtifactRelPath", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc, false)
		relPath, err := moduleLoc.GetSourceModuleArtifactRelPath("path1", getPath("testdata", "path1", "data.zip"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(Equal(""))
	})

	It("GetTargetModuleDir", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc, false)
		Ω(moduleLoc.GetTargetModuleDir("module1")).Should(Equal(getPath("testdata")))
	})

	It("ParseFile", func() {
		ep := Loc{SourcePath: getPath("testdata", "testext")}
		moduleLoc := ModuleLocation(&ep, false)
		mta, err := moduleLoc.ParseFile()
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(Succeed())

		module1, err := mta.GetModuleByName("ui5app")
		Ω(err).Should(Succeed())
		Ω(module1.Properties).Should(BeNil())

		module2, err := mta.GetModuleByName("ui5app2")
		Ω(err).Should(Succeed())
		Ω(module2.Parameters).ShouldNot(BeNil())
		Ω(module2.Parameters["memory"]).Should(Equal("256M"))
	})

	It("SetStrictParmeter", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc, false)
		Ω(moduleLoc.SetStrictParmeter(true)).Should(Equal(true))
	})

	It("ModuleLocation", func() {
		target := getPath("testdata", "result")
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ModuleLocation(projectLoc, false)).ShouldNot(BeNil())
	})
})
