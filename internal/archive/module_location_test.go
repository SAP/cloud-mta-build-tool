package dir

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("ModuleLocation", func() {
	It("GetTarget", func() {
		loc := &ModuleLoc{targetPath: getPath("test")}
		Ω(loc.GetTarget()).Should(Equal(getPath("test")))
	})

	It("GetTargetTmpDir - default target", func() {
		projectLoc, err := Location(getPath("testdata"), "", Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := &ModuleLoc{loc: projectLoc}
		Ω(moduleLoc.GetTargetTmpDir()).Should(Equal(getPath("testdata", ".testdata_mta_build_tmp")))
	})

	It("GetTargetTmpDir - target provided", func() {
		target := getPath("testdata", "result")
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := &ModuleLoc{loc: projectLoc, targetPath: target}
		Ω(moduleLoc.GetTargetTmpDir()).Should(Equal(target))
	})

	It("GetSourceModuleDir", func() {
		target := ""
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := &ModuleLoc{loc: projectLoc, targetPath: target}
		Ω(moduleLoc.GetSourceModuleDir("path1")).Should(Equal((getPath("testdata", "path1"))))
	})

	It("GetSourceModuleArtifactRelPath", func() {
		target := ""
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := &ModuleLoc{loc: projectLoc, targetPath: target}
		Ω(moduleLoc.GetSourceModuleArtifactRelPath("path1", getPath("testdata", "path1", "data.zip"), false)).Should(Equal(""))
	})

	It("GetTargetModuleDir - default target", func() {
		target := ""
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := &ModuleLoc{loc: projectLoc, targetPath: target}
		Ω(moduleLoc.GetTargetModuleDir("module1")).Should(Equal(getPath("testdata", ".testdata_mta_build_tmp", "module1")))
	})

	It("GetTargetModuleDir - target provided", func() {
		target := getPath("testdata", "result")
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		moduleLoc := &ModuleLoc{loc: projectLoc, targetPath: target}
		Ω(moduleLoc.GetTargetModuleDir("module1")).Should(Equal(target))
	})

	It("ModuleLocation", func() {
		target := getPath("testdata", "result")
		projectLoc, err := Location(getPath("testdata"), target, Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ModuleLocation(projectLoc, target)).ShouldNot(BeNil())
	})
})
