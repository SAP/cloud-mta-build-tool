package dir

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta/mta"
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
	It("GetTarget - Explicit", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTarget()).Should(Equal(getPath("abc")))
	})
	It("GetTargetTmpDir", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTargetTmpDir()).Should(Equal(getPath("abc", ".xyz_mta_build_tmp")))
	})
	It("GetTargetModuleDir", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTargetModuleDir("mmm")).Should(
			Equal(getPath("abc", ".xyz_mta_build_tmp", "mmm")))
	})
	It("GetTargetModuleZipPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetTargetModuleZipPath("mmm")).Should(
			Equal(getPath("abc", ".xyz_mta_build_tmp", "mmm", "data.zip")))
	})
	It("GetSourceModuleDir", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetSourceModuleDir("mpath")).Should(Equal(getPath("xyz", "mpath")))
	})
	It("getMtaYamlFilename - Explicit", func() {
		location := Loc{MtaFilename: "mymta.yaml"}
		Ω(location.GetMtaYamlFilename()).Should(Equal("mymta.yaml"))
	})
	It("getMtaYamlFilename - Implicit", func() {
		location := Loc{}
		Ω(location.GetMtaYamlFilename()).Should(Equal("mta.yaml"))
	})
	It("getMtaYamlFilename - Implicit- MTAD", func() {
		location := Loc{Descriptor: Dep}
		Ω(location.GetMtaYamlFilename()).Should(Equal("mtad.yaml"))
	})
	It("GetMtaYamlPath", func() {
		location := Loc{SourcePath: getPath()}
		Ω(location.GetMtaYamlPath()).Should(Equal(getPath("mta.yaml")))
	})
	It("GetMetaPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMetaPath()).Should(Equal(getPath("abc", ".xyz_mta_build_tmp", "META-INF")))
	})
	It("GetMtadPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMtadPath()).Should(Equal(getPath("abc", ".xyz_mta_build_tmp", "META-INF", "mtad.yaml")))
	})
	It("GetMtarDir - mta_archives subfolder", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("xyz")}
		Ω(location.GetMtarDir()).Should(Equal(getPath("xyz", "mta_archives")))
	})
	It("GetMtarDir - target folder", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMtarDir()).Should(Equal(getPath("abc")))
	})
	It("GetManifestPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetManifestPath()).Should(Equal(getPath("abc", ".xyz_mta_build_tmp", "META-INF", "MANIFEST.MF")))
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

var _ = Describe("ParseFile MTA", func() {

	wd, _ := os.Getwd()

	It("Valid filename", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata")}
		mta, err := ep.ParseFile()
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})
	It("Invalid filename", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtax.yaml"}
		_, err := ep.ParseFile()
		Ω(err).ShouldNot(BeNil())
	})
})

var _ = Describe("ParseExtFile MTA", func() {

	wd, _ := os.Getwd()

	It("Valid filename", func() {
		ep := Loc{SourcePath: filepath.Join(wd, "testdata", "testproject")}
		mta, err := ep.ParseExtFile("cf")
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})
	It("Invalid filename", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata", "testproject"), MtaFilename: "mtax.yaml"}
		Ω(ep.ParseExtFile("neo")).Should(Equal(&mta.EXT{}))
	})
})

var _ = Describe("Location", func() {
	It("Dev Descritor", func() {
		ep, err := Location("", "", "", os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ep.GetMtaYamlFilename()).Should(Equal("mta.yaml"))
	})
	It("Dep Descriptor", func() {
		ep, err := Location("", "", Dep, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ep.GetMtaYamlFilename()).Should(Equal("mtad.yaml"))
	})
	It("Dev Descriptor - Explicit", func() {
		ep, err := Location("", "", Dep, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ep.GetDescriptor()).Should(Equal(Dep))
	})
	It("Dev Descriptor - Implicit", func() {
		ep := &Loc{Descriptor: ""}
		Ω(ep.GetDescriptor()).Should(Equal(Dev))
	})
	It("Fails on descriptor validation", func() {
		_, err := Location("", "", "xx", os.Getwd)
		Ω(err).Should(HaveOccurred())
	})
	It("Fails on implicit source", func() {
		_, err := Location("", "", Dev, func() (string, error) {
			return "", errors.New("err")
		})
		Ω(err).Should(HaveOccurred())
	})
})
