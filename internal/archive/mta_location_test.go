package dir

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta/mta"
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

var _ = Describe("ParseMtaFile", func() {
	wd, _ := os.Getwd()

	It("Valid filename", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata")}
		mta, err := ep.ParseMtaFile()
		Ω(mta).ShouldNot(BeNil())
		Ω(err).Should(Succeed())
	})
	It("Invalid filename", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mtax.yaml"}
		_, err := ep.ParseMtaFile()
		Ω(err).Should(HaveOccurred())
	})
})

var _ = Describe("ParseExtFile", func() {

	wd, _ := os.Getwd()

	It("Valid filename", func() {
		ep := Loc{SourcePath: filepath.Join(wd, "testdata", "testext")}
		mta, err := ep.ParseExtFile("cf-mtaext.yaml")
		Ω(err).Should(Succeed())
		Ω(mta).ShouldNot(BeNil())
	})
	It("Invalid filename", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata", "testext")}
		_, e := ep.ParseExtFile("neo-mtaext.yaml")
		Ω(e).Should(HaveOccurred())
	})
	It("Invalid file content", func() {
		ep := &Loc{SourcePath: filepath.Join(wd, "testdata", "testext")}
		_, e := ep.ParseExtFile("invalid.mtaext")
		Ω(e).Should(HaveOccurred())
		Ω(e.Error()).Should(ContainSubstring(fmt.Sprintf(parseExtFileFailed, ep.GetMtaExtYamlPath("invalid.mtaext"))))
	})
})

var _ = Describe("ParseFile", func() {
	wd, _ := os.Getwd()

	It("Parse the mta.yaml file and returns it when there are no extension files", func() {
		ep := Loc{SourcePath: filepath.Join(wd, "testdata", "testext")}
		mta, err := ep.ParseFile()
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
		mta, err := ep.ParseFile()
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
		mta, err := ep.ParseFile()
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
	})

	It("fails when an extension version mismatches the MTA version", func() {
		ep := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"bad_version.mtaext", "other.mtaext"},
		}
		_, err := ep.ParseFile()
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(versionMismatchMsg, "3.1", "mtahtml5ext", "2.1")))
	})

	It("fails when an the extension cannot be merged", func() {
		ep := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"bad_module.mtaext"},
		}
		_, err := ep.ParseFile()
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("ui5app3"))
	})
})

var _ = Describe("getSortedExtensions", func() {
	wd, _ := os.Getwd()

	It("fails when one of the files cannot be read", func() {
		loc := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"cf-mtaext.yaml", "unknownfile.mtaext"},
		}
		_, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(ReadFailedMsg, loc.GetMtaExtYamlPath("unknownfile.mtaext"))))
	})
	It("fails when there are several extensions with the same ID", func() {
		// This takes care of the cyclic extends case too when the extension's ID is not the mta.yaml ID
		// (because then there will be 2 extensions with the same ID)
		loc := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"cf-mtaext.yaml", "other.mtaext", "third.mtaext", "third_copy_diff_extends.mtaext"},
		}
		_, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(duplicateExtensionIDMsg,
			loc.GetMtaExtYamlPath("third.mtaext"),
			loc.GetMtaExtYamlPath("third_copy_diff_extends.mtaext"),
			"mtahtml5ext3"),
		))
	})
	It("fails when there are several extensions that extend the same ID", func() {
		loc := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"cf-mtaext.yaml", "other.mtaext", "third.mtaext", "third_copy_diff_id.mtaext"},
		}
		_, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(duplicateExtendsMsg,
			loc.GetMtaExtYamlPath("third.mtaext"),
			loc.GetMtaExtYamlPath("third_copy_diff_id.mtaext"),
			"mtahtml5ext2"),
		))
	})
	It("fails when there are extensions that extend unknown files", func() {
		loc := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"cf-mtaext.yaml", "third.mtaext", "unknown_extends.mtaext"},
		}
		_, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(unknownExtendsMsg, "")))
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(extendsMsg, loc.GetMtaExtYamlPath("third.mtaext"), "mtahtml5ext2")))
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(extendsMsg, loc.GetMtaExtYamlPath("unknown_extends.mtaext"), "mtahtml5unknown")))
	})
	It("fails when there is an extension with the MTA ID", func() {
		// This covers the cyclic case too (cf-mtaext.yaml <-> mtaid.mtaext)
		loc := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"cf-mtaext.yaml", "mtaid.mtaext"},
		}
		_, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(extensionIDSameAsMtaIDMsg,
			loc.GetMtaExtYamlPath("mtaid.mtaext"), "mtahtml5", loc.GetMtaYamlFilename()),
		))
	})
	It("fails when none of the extensions extends the MTA", func() {
		loc := Loc{
			SourcePath:         filepath.Join(wd, "testdata", "testext"),
			ExtensionFileNames: []string{"other.mtaext", "third.mtaext"},
		}
		_, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(unknownExtendsMsg, "")))
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(extendsMsg, loc.GetMtaExtYamlPath("other.mtaext"), "mtahtml5ext")))
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(extendsMsg, loc.GetMtaExtYamlPath("third.mtaext"), "mtahtml5ext2")))
	})
	DescribeTable("returns the extensions sorted by extends chain order", func(files []string, expectedIDsOrder []string) {
		loc := Loc{SourcePath: filepath.Join(wd, "testdata", "testext"), ExtensionFileNames: files}
		extensions, err := loc.getSortedExtensions("mtahtml5")
		Ω(err).Should(Succeed())
		extIDs := make([]string, 0)
		for _, ext := range extensions {
			extIDs = append(extIDs, ext.ID)
		}
		Ω(extIDs).Should(Equal(expectedIDsOrder))
	},
		Entry("nil table", nil, []string{}),
		Entry("empty table", []string{}, []string{}),
		Entry("there is only one entry", []string{"cf-mtaext.yaml"}, []string{"mtahtml5ext"}),
		Entry("extensions are in order of chain", []string{"cf-mtaext.yaml", "other.mtaext", "third.mtaext"}, []string{"mtahtml5ext", "mtahtml5ext2", "mtahtml5ext3"}),
		Entry("extensions are not in the order of the chain", []string{"third.mtaext", "cf-mtaext.yaml", "other.mtaext"}, []string{"mtahtml5ext", "mtahtml5ext2", "mtahtml5ext3"}),
	)
})

var _ = DescribeTable("checkSchemaVersionMatches", func(mtaVersion *string, extVersion *string, expectedError bool) {
	err := checkSchemaVersionMatches(&mta.MTA{SchemaVersion: mtaVersion}, &mta.EXT{SchemaVersion: extVersion})
	if expectedError {
		Ω(err).Should(HaveOccurred())
	} else {
		Ω(err).Should(Succeed())
	}
},
	Entry("nil versions", nil, nil, false),
	Entry("empty versions", ptr(""), ptr(""), false),
	Entry("mta version is empty, ext version isn't empty", ptr(""), ptr("3.1"), true),
	Entry("mta version isn't empty, ext version is empty", ptr("2.1"), ptr(""), true),
	Entry("different major versions when minor version is specified", ptr("2.1"), ptr("3.1"), true),
	Entry("different major versions when minor version isn't specified", ptr("2"), ptr("3"), true),
	Entry("different minor versions", ptr("3.3"), ptr("3.1"), false),
	Entry("mta version is major.minor, ext only has major part", ptr("3.1"), ptr("3"), false),
	Entry("mta version only has major part, ext is major.minor", ptr("3"), ptr("3.2"), false),
	Entry("different patch version", ptr("3.2.1"), ptr("3.2.2"), false),
	Entry("only mta has patch version", ptr("3.2.1"), ptr("3.2"), false),
	Entry("same version - major", ptr("3"), ptr("3"), false),
	Entry("same version - major.minor", ptr("3.3"), ptr("3.3"), false),
	Entry("same version - major.minor.patch", ptr("3.4.5"), ptr("3.4.5"), false),
)

func ptr(str string) *string {
	return &str
}

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
