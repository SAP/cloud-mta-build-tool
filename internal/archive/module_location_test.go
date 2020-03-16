package dir

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("ModuleLocation", func() {
	It("GetTarget", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("test"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		loc := ModuleLocation(projectLoc)
		Ω(loc.GetTarget()).Should(Equal(getPath("test")))
	})

	It("GetTargetTmpDir", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc)
		Ω(moduleLoc.GetTargetTmpDir()).Should(Equal(getPath("testdata")))
	})

	It("GetSourceModuleDir", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc)
		Ω(moduleLoc.GetSourceModuleDir("path1")).Should(Equal((getPath("testdata", "path1"))))
	})

	It("GetSourceModuleArtifactRelPath", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc)
		relPath, err := moduleLoc.GetSourceModuleArtifactRelPath("path1", getPath("testdata", "path1", "data.zip"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(Equal(""))
	})

	It("GetTargetModuleDir", func() {
		projectLoc, err := Location(getPath("testdata"), getPath("testdata"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc)
		Ω(moduleLoc.GetTargetModuleDir("module1")).Should(Equal(getPath("testdata")))
	})

	It("ParseExtFile", func() {
		projectLoc, err := Location(getPath("testdata", "testext"), "", Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := ModuleLocation(projectLoc)
		mta, err := moduleLoc.ParseExtFile("cf-mtaext.yaml")
		Ω(err).Should(Succeed())
		Ω(mta).ShouldNot(BeNil())
	})

	It("ParseFile", func() {
		ep := Loc{SourcePath: getPath("testdata", "testext")}
		moduleLoc := ModuleLocation(&ep)
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

	It("ModuleLocation", func() {
		target := getPath("testdata", "result")
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ModuleLocation(projectLoc)).ShouldNot(BeNil())
	})
})
