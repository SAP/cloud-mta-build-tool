package dir

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	It("GetRelativePath", func() {
		Ω(GetRelativePath(getFullPath("abc", "xyz", "fff"),
			filepath.Join(getFullPath()))).Should(Equal(string(filepath.Separator) + filepath.Join("abc", "xyz", "fff")))
	})
	It("GetSource - Explicit", func() {
		Ω(EndPoints{SourcePath: getFullPath("abc")}.GetSource()).Should(Equal(getFullPath("abc")))
	})
	It("GetSource - Implicit", func() {
		Ω(EndPoints{}.GetSource()).Should(Equal(getFullPath()))
	})
	It("GetTarget - Explicit", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetTarget()).Should(Equal(getFullPath("abc")))
	})
	It("GetTarget - Implicit", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz")}.GetTarget()).Should(Equal(getFullPath("xyz")))
	})
	It("GetTargetTmpDir", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetTargetTmpDir()).Should(Equal(getFullPath("abc", "xyz")))
	})
	It("GetTargetModuleDir", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetTargetModuleDir("mmm")).Should(Equal(getFullPath("abc", "xyz", "mmm")))
	})
	It("GetTargetModuleZipPath", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetTargetModuleZipPath("mmm")).Should(Equal(getFullPath("abc", "xyz", "mmm", "data.zip")))
	})
	It("GetSourceModuleDir", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetSourceModuleDir("mpath")).Should(Equal(getFullPath("xyz", "mpath")))
	})
	It("GetMtaYamlFilename - Explicit", func() {
		Ω(EndPoints{MtaFilename: "mymta.yaml"}.GetMtaYamlFilename()).Should(Equal("mymta.yaml"))
	})
	It("GetMtaYamlFilename - Implicit", func() {
		Ω(EndPoints{}.GetMtaYamlFilename()).Should(Equal("mta.yaml"))
	})
	It("GetMtaYamlPath", func() {
		Ω(EndPoints{}.GetMtaYamlPath()).Should(Equal(getFullPath("mta.yaml")))
	})
	It("GetMetaPath", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetMetaPath()).Should(Equal(getFullPath("abc", "xyz", "META-INF")))
	})
	It("GetMtadPath", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetMtadPath()).Should(Equal(getFullPath("abc", "xyz", "META-INF", "mtad.yaml")))
	})
	It("GetManifestPath", func() {
		Ω(EndPoints{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}.GetManifestPath()).Should(Equal(getFullPath("abc", "xyz", "META-INF", "MANIFEST.MF")))
	})
})
