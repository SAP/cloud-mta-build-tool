package dir

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	It("getRelativePath", func() {
		Ω(getRelativePath(getFullPath("abc", "xyz", "fff"),
			filepath.Join(getFullPath()))).Should(Equal(string(filepath.Separator) + filepath.Join("abc", "xyz", "fff")))
	})
	It("GetSource - Explicit", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("abc")}
		Ω(location.GetSource()).Should(Equal(getFullPath("abc")))
	})
	It("GetSource - Implicit", func() {
		location := MtaLocationParameters{}
		Ω(location.GetSource()).Should(Equal(getFullPath()))
	})
	It("GetTarget - Explicit", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetTarget()).Should(Equal(getFullPath("abc")))
	})
	It("GetTarget - Implicit", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz")}
		Ω(location.GetTarget()).Should(Equal(getFullPath("xyz")))
	})
	It("GetTargetTmpDir", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetTargetTmpDir()).Should(Equal(getFullPath("abc", "xyz")))
	})
	It("GetTargetModuleDir", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetTargetModuleDir("mmm")).Should(Equal(getFullPath("abc", "xyz", "mmm")))
	})
	It("GetTargetModuleZipPath", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetTargetModuleZipPath("mmm")).Should(Equal(getFullPath("abc", "xyz", "mmm", "data.zip")))
	})
	It("GetSourceModuleDir", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetSourceModuleDir("mpath")).Should(Equal(getFullPath("xyz", "mpath")))
	})
	It("getMtaYamlFilename - Explicit", func() {
		location := MtaLocationParameters{MtaFilename: "mymta.yaml"}
		Ω(location.getMtaYamlFilename()).Should(Equal("mymta.yaml"))
	})
	It("getMtaYamlFilename - Implicit", func() {
		location := MtaLocationParameters{}
		Ω(location.getMtaYamlFilename()).Should(Equal("mta.yaml"))
	})
	It("GetMtaYamlPath", func() {
		location := MtaLocationParameters{}
		Ω(location.GetMtaYamlPath()).Should(Equal(getFullPath("mta.yaml")))
	})
	It("GetMetaPath", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetMetaPath()).Should(Equal(getFullPath("abc", "xyz", "META-INF")))
	})
	It("GetMtadPath", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetMtadPath()).Should(Equal(getFullPath("abc", "xyz", "META-INF", "mtad.yaml")))
	})
	It("GetManifestPath", func() {
		location := MtaLocationParameters{SourcePath: getFullPath("xyz"), TargetPath: getFullPath("abc")}
		Ω(location.GetManifestPath()).Should(Equal(getFullPath("abc", "xyz", "META-INF", "MANIFEST.MF")))
	})
})
