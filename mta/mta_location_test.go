package mta

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func getPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, filepath.Join(relPath...))
}

var _ = Describe("Path", func() {

	It("GetSource - Explicit", func() {
		location := Loc{SourcePath: getPath("abc")}
		Ω(location.GetSource()).Should(Equal(getPath("abc")))
	})
	It("GetSource - Implicit", func() {
		location := Loc{}
		Ω(location.GetSource()).Should(Equal(getPath()))
	})
	It("GetTarget - Explicit", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTarget()).Should(Equal(getPath("abc")))
	})
	It("GetTarget - Implicit", func() {
		location := Loc{SourcePath: getPath("xyz")}
		Ω(location.GetTarget()).Should(Equal(getPath("xyz")))
	})
	It("GetTargetTmpDir", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTargetTmpDir()).Should(Equal(getPath("abc", "xyz")))
	})
	It("GetTargetModuleDir", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTargetModuleDir("mmm")).Should(Equal(getPath("abc", "xyz", "mmm")))
	})
	It("GetTargetModuleZipPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTargetModuleZipPath("mmm")).Should(Equal(getPath("abc", "xyz", "mmm", "data.zip")))
	})
	It("GetSourceModuleDir", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetSourceModuleDir("mpath")).Should(Equal(getPath("xyz", "mpath")))
	})
	It("getMtaYamlFilename - Explicit", func() {
		location := Loc{MtaFilename: "mymta.yaml"}
		Ω(location.getMtaYamlFilename()).Should(Equal("mymta.yaml"))
	})
	It("getMtaYamlFilename - Implicit", func() {
		location := Loc{}
		Ω(location.getMtaYamlFilename()).Should(Equal("mta.yaml"))
	})
	It("getMtaYamlFilename - Implicit- MTAD", func() {
		location := Loc{Descriptor: "dep"}
		Ω(location.getMtaYamlFilename()).Should(Equal("mtad.yaml"))
	})
	It("GetMtaYamlPath", func() {
		location := Loc{}
		Ω(location.GetMtaYamlPath()).Should(Equal(getPath("mta.yaml")))
	})
	It("GetMetaPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMetaPath()).Should(Equal(getPath("abc", "xyz", "META-INF")))
	})
	It("GetMtadPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMtadPath()).Should(Equal(getPath("abc", "xyz", "META-INF", "mtad.yaml")))
	})
	It("GetManifestPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetManifestPath()).Should(Equal(getPath("abc", "xyz", "META-INF", "MANIFEST.MF")))
	})
	It("ValidateDeploymentDescriptor - Valid", func() {
		Ω(ValidateDeploymentDescriptor("")).Should(Succeed())
	})
	It("ValidateDeploymentDescriptor - Invalid", func() {
		Ω(ValidateDeploymentDescriptor("xxx")).Should(HaveOccurred())
	})
	It("IsDeploymentDescriptor", func() {
		location := Loc{}
		Ω(location.IsDeploymentDescriptor()).Should(Equal(false))
	})
})

var _ = Describe("Path Failures", func() {

	var storedWorkingDirectory func() (string, error)
	lp := Loc{}

	BeforeEach(func() {
		storedWorkingDirectory = GetWorkingDirectory
		GetWorkingDirectory = func() (string, error) {
			return "", errors.New("Dummy error")
		}
	})

	AfterEach(func() {
		GetWorkingDirectory = storedWorkingDirectory
	})

	It("GetSource - Implicit", func() {
		_, err := lp.GetSource()
		Ω(err).Should(HaveOccurred())
	})
	It("GetTarget - Implicit", func() {
		_, err := lp.GetTarget()
		Ω(err).Should(HaveOccurred())
	})
	It("GetTargetTmpDir", func() {
		_, err := lp.GetTargetTmpDir()
		Ω(err).Should(HaveOccurred())
	})
	It("GetTargetModuleDir", func() {
		_, err := lp.GetTargetModuleDir("mmm")
		Ω(err).Should(HaveOccurred())
	})
	It("GetTargetModuleZipPath", func() {
		_, err := lp.GetTargetModuleZipPath("mmm")
		Ω(err).Should(HaveOccurred())
	})
	It("GetSourceModuleDir", func() {
		_, err := lp.GetSourceModuleDir("mpath")
		Ω(err).Should(HaveOccurred())
	})
	It("GetMtaYamlPath", func() {
		_, err := lp.GetMtaYamlPath()
		Ω(err).Should(HaveOccurred())
	})
	It("GetMetaPath", func() {
		_, err := lp.GetMetaPath()
		Ω(err).Should(HaveOccurred())
	})
	It("GetMtadPath", func() {
		_, err := lp.GetMtadPath()
		Ω(err).Should(HaveOccurred())
	})
	It("GetManifestPath", func() {
		_, err := lp.GetManifestPath()
		Ω(err).Should(HaveOccurred())
	})
})
