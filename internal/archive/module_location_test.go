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

	It("ModuleLocation", func() {
		target := getPath("testdata", "result")
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ModuleLocation(projectLoc)).ShouldNot(BeNil())
	})
})
