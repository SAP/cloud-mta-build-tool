package dir

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getPath(relPath ...string) string {
	wd, err := os.Getwd()
	Ω(err).Should(Succeed())
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
	It("GetSourceModuleArtifactRelPath - artifact is file", func() {
		location := Loc{SourcePath: getPath("testdata", "mtahtml5")}
		relPath, err := location.GetSourceModuleArtifactRelPath("ui5app", getPath("testdata", "mtahtml5", "ui5app", "webapp", "Component.js"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(Equal("webapp"))
	})
	It("GetSourceModuleArtifactRelPath - artifact is an archive", func() {
		location := Loc{SourcePath: getPath("testdata", "mtahtml5")}
		relPath, err := location.GetSourceModuleArtifactRelPath("ui5app", getPath("testdata", "mtahtml5", "ui5app", "webapp", "abc.jar"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(Equal("webapp"))
	})
	It("GetSourceModuleArtifactRelPath - artifact does not exist", func() {
		location := Loc{SourcePath: getPath("testdata", "mtahtml5")}
		_, err := location.GetSourceModuleArtifactRelPath("ui5app", getPath("testdata", "mtahtml5", "ui5app", "webapp", "ComponentA.js"))
		Ω(err).Should(HaveOccurred())
	})
	It("GetSourceModuleArtifactRelPath - artifact is folder", func() {
		location := Loc{SourcePath: getPath("testdata", "mtahtml5")}
		relPath, err := location.GetSourceModuleArtifactRelPath("ui5app", getPath("testdata", "mtahtml5", "ui5app", "webapp", "view"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(Equal(filepath.Join("webapp", "view")))
	})
	It("GetSourceModuleArtifactRelPath - artifact is module itself", func() {
		location := Loc{SourcePath: getPath("testdata", "mtahtml5")}
		relPath, err := location.GetSourceModuleArtifactRelPath("ui5app", getPath("testdata", "mtahtml5", "ui5app"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(Equal("."))
	})
	It("GetSourceModuleArtifactRelPath - artifact is module which is file", func() {
		location := Loc{SourcePath: getPath("testdata", "mtahtml5")}
		relPath, err := location.GetSourceModuleArtifactRelPath(filepath.Join("ui5app", "Gruntfile.js"), getPath("testdata", "mtahtml5", "ui5app", "Gruntfile.js"))
		Ω(err).Should(Succeed())
		Ω(relPath).Should(BeEmpty())
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
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMtarDir(false)).Should(Equal(getPath("abc", "mta_archives")))
	})
	It("GetMtarDir - target folder", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetMtarDir(true)).Should(Equal(getPath("abc")))
	})
	It("GetManifestPath", func() {
		location := Loc{SourcePath: getPath("xyz"), TargetPath: getPath("abc")}
		Ω(location.GetManifestPath()).Should(Equal(getPath("abc", ".xyz_mta_build_tmp", "META-INF", "MANIFEST.MF")))
	})
	It("ValidateDeploymentDescriptor - Valid", func() {
		Ω(ValidateDeploymentDescriptor("")).Should(Succeed())
	})
	It("ValidateDeploymentDescriptor - Invalid", func() {
		err := ValidateDeploymentDescriptor("xxx")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(Equal(fmt.Sprintf(InvalidDescMsg, "xxx")))
	})
	It("IsDeploymentDescriptor", func() {
		location := Loc{}
		Ω(location.IsDeploymentDescriptor()).Should(Equal(false))
	})

	It("GetMtaExtYamlPath when path is relative returns the path relative to the source folder", func() {
		location := Loc{SourcePath: getPath("xyz")}
		Ω(location.GetMtaExtYamlPath("somefile.mtaext")).Should(Equal(getPath("xyz", "somefile.mtaext")))
		Ω(location.GetMtaExtYamlPath(filepath.Join("innerfolder", "somefile"))).Should(Equal(getPath("xyz", "innerfolder", "somefile")))
	})

	It("GetMtaExtYamlPath when path is absolute returns the same path", func() {
		location := Loc{SourcePath: getPath("xyz")}
		pathToMtaExt := filepath.Join("a_folder", "somefile.mtaext")
		absPath, err := filepath.Abs(pathToMtaExt)
		Ω(err).Should(Succeed())
		Ω(location.GetMtaExtYamlPath(absPath)).Should(Equal(absPath))
		Ω(location.GetMtaExtYamlPath(pathToMtaExt)).ShouldNot(Equal(absPath))
	})

	It("GetTargetTmpRoot, target path calculated", func() {
		projectLoc, err := Location(getPath("test"), "", Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(projectLoc.GetTargetTmpRoot()).Should(Equal(getPath("test", ".test_mta_build_tmp")))
	})

	It("GetTargetTmpRoot, target path defined", func() {
		projectLoc, err := Location(getPath("test"), getPath("test1"), Dev, []string{}, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(projectLoc.GetTargetTmpRoot()).Should(Equal(getPath("test1", ".test_mta_build_tmp")))
	})

	Describe("GetExtensionFilePaths", func() {
		It("returns empty list when ExtensionFileNames is nil", func() {
			loc := Loc{}
			Ω(loc.GetExtensionFilePaths()).Should(BeEmpty())
		})
		It("returns empty list when ExtensionFileNames is an empty list", func() {
			loc := Loc{ExtensionFileNames: []string{}}
			Ω(loc.GetExtensionFilePaths()).Should(BeEmpty())
		})
		It("returns one item when ExtensionFileNames has one item", func() {
			loc := Loc{SourcePath: getPath("xyz"), ExtensionFileNames: []string{"a.mtaext"}}
			Ω(loc.GetExtensionFilePaths()).Should(Equal([]string{getPath("xyz", "a.mtaext")}))
		})
		It("returns an item for each ExtensionFileNames entry", func() {
			bPath, err := filepath.Abs("b.mtaext")
			Ω(err).Should(Succeed())

			loc := Loc{SourcePath: getPath("xyz"), ExtensionFileNames: []string{"a.mtaext", bPath, "somefile"}}
			Ω(loc.GetExtensionFilePaths()).Should(Equal([]string{
				getPath("xyz", "a.mtaext"),
				bPath,
				getPath("xyz", "somefile"),
			}))
		})
	})
})

var _ = Describe("ParseFile", func() {
	wd, _ := os.Getwd()

	It("Parse the mta.yaml file and returns it when there are no extension files", func() {
		ep := Loc{SourcePath: filepath.Join(wd, "testdata", "testext")}
		mta, err := ep.ParseFile(true)
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

	It("Parses the mta.yaml and merges the extension file when there is one extension file", func() {
		ep := Loc{SourcePath: filepath.Join(wd, "testdata", "testext"), ExtensionFileNames: []string{"cf-mtaext.yaml"}}
		mta, err := ep.ParseFile(true)
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(Succeed())

		module1, err := mta.GetModuleByName("ui5app")
		Ω(err).Should(Succeed())
		Ω(module1.Properties).ShouldNot(BeNil())
		Ω(module1.Properties["my_prop"]).Should(Equal(1))

		module2, err := mta.GetModuleByName("ui5app2")
		Ω(err).Should(Succeed())
		Ω(module2.Parameters).ShouldNot(BeNil())
		Ω(module2.Parameters["memory"]).Should(Equal("512M"))
	})

	It("Parses the mta.yaml file and merges all extension files when there are several extension files", func() {
		ep := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"other.mtaext", "cf-mtaext.yaml"},
		}
		mta, err := ep.ParseFile(true)
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(Succeed())

		module1, err := mta.GetModuleByName("ui5app")
		Ω(err).Should(Succeed())
		Ω(module1.Properties).ShouldNot(BeNil())
		Ω(module1.Properties["my_prop"]).Should(Equal(1))
		Ω(module1.Properties["other_prop"]).Should(Equal("abc"))

		module2, err := mta.GetModuleByName("ui5app2")
		Ω(err).Should(Succeed())
		Ω(module2.Parameters).ShouldNot(BeNil())
		Ω(module2.Parameters["memory"]).Should(Equal("1024M"))

		resource := mta.GetResourceByName("uaa_mtahtml5")
		Ω(resource).ShouldNot(BeNil())
		Ω(resource.Active).ShouldNot(BeNil())
		Ω(*resource.Active).Should(BeFalse())
	})

	It("fails on not existing file", func() {
		ep := Loc{
			SourcePath:  filepath.Join(wd, "testdata", "testext"),
			MtaFilename: "some.yaml",
		}
		_, err := ep.ParseFile(true)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(filepath.Join(ep.SourcePath, ep.MtaFilename)))
	})

	It("fails when an extension version mismatches the MTA version", func() {
		ep := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"bad_version.mtaext", "other.mtaext"},
		}
		_, err := ep.ParseFile(true)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("3.1"))
	})

	It("fails when the extension cannot be merged", func() {
		ep := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"bad_module.mtaext"},
		}
		_, err := ep.ParseFile(true)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("ui5app3"))
	})

	It("fails when the sent extensions don't extend this mta.yaml", func() {
		ep := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"other.mtaext"},
		}
		_, err := ep.ParseFile(true)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("mtahtml5ext"))
	})
})

var _ = Describe("Location", func() {
	It("Dev Descritor", func() {
		ep, err := Location("", "", "", nil, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ep.GetMtaYamlFilename()).Should(Equal("mta.yaml"))
	})
	It("Dep Descriptor", func() {
		ep, err := Location("", "", Dep, nil, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ep.GetMtaYamlFilename()).Should(Equal("mtad.yaml"))
	})
	It("Dev Descriptor - Explicit", func() {
		ep, err := Location("", "", Dep, nil, os.Getwd)
		Ω(err).Should(Succeed())
		Ω(ep.GetDescriptor()).Should(Equal(Dep))
	})
	It("Dev Descriptor - Implicit", func() {
		ep := &Loc{Descriptor: ""}
		Ω(ep.GetDescriptor()).Should(Equal(Dev))
	})
	It("Fails on descriptor validation", func() {
		_, err := Location("", "", "xx", nil, os.Getwd)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(InvalidDescMsg, "xx")))
	})
	It("Fails when it can't get the current working directory", func() {
		_, err := Location("", "", Dev, nil, func() (string, error) {
			return "", errors.New("err")
		})
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(InitLocFailedOnWorkDirMsg))
	})
})
