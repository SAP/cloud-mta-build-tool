package artifacts

import (
	"os"
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/conttype"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("manifest", func() {

	BeforeEach(func() {
		createDirInTmpFolder("mta", "META-INF")
	})

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})

	var _ = Describe("setManifestDesc", func() {
		It("Sanity", func() {
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "data.zip")
			createFileInTmpFolder("mta", "config-site-host.json")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
			actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			Ω(actual).Should(Equal(golden))
		})
		It("Unknown content type, assembly scenario", func() {
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "server.js")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), Descriptor: dir.Dep}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, []*mta.Resource{}, "cf")
			checkError(err, conttype.ContentTypeUndefinedMsg, ".js")
		})
		It("Sanity - with configuration provided", func() {
			// no_source module (with no-source build parameter) is not referenced in the manifest
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "data.zip")
			createFileInTmpFolder("mta", "config-site-host.json")
			createFileInTmpFolder("mta", "config-site-host1.json")
			createFileInTmpFolder("mta", "xs-security.json")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_cfg.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
			actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest_cfg.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			Ω(actual).Should(Equal(golden))
		})
		It("wrong Commands configuration", func() {
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "data.zip")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			moduleConf := commands.ModuleTypeConfig
			commands.ModuleTypeConfig = []byte("bad module conf")
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(HaveOccurred())
			commands.ModuleTypeConfig = moduleConf
		})
		It("module with defined build-result fails when build-result file does not exist in source directory", func() {
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "some.war")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(HaveOccurred())
		})
		It("module with defined build-result fails when build-result file does not exist in target temp directory", func() {
			createDirInTmpFolder("mta", "node-js")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult2.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(HaveOccurred())
		})
		It("entry for module with defined build-result has the build-result file", func() {
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "data1.zip")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildResult.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
			actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifestBuildResult.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			Ω(actual).Should(Equal(golden))
		})
		It("wrong content types configuration", func() {
			createDirInTmpFolder("mta", "node-js")
			createFileInTmpFolder("mta", "node-js", "data.zip")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			contentTypesOrig := conttype.ContentTypeConfig
			conttype.ContentTypeConfig = []byte(`wrong configuraion`)
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(HaveOccurred())
			conttype.ContentTypeConfig = contentTypesOrig

		})
		It("Sanity - module with no path is not added to the manifest", func() {
			createDirInTmpFolder("mta", "node-js")
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_no_paths.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
			actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_assembly_manifest_no_paths.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			Ω(actual).Should(Equal(golden))
		})
		It("With resources", func() {
			createDirInTmpFolder("assembly-sample", "META-INF")
			createDirInTmpFolder("assembly-sample", "web")
			createFileInTmpFolder("assembly-sample", "config-site-host.json")
			createFileInTmpFolder("assembly-sample", "xs-security.json")
			createDirInTmpFolder("assembly-sample", "inner")
			createFileInTmpFolder("assembly-sample", "inner", "uaa.json")
			createFileInTmpFolder("assembly-sample", "inner", "xs-security.json")
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources, "cf")).Should(Succeed())
			actual := getFileContent(getFullPathInTmpFolder("assembly-sample", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_assembly_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			Ω(actual).Should(Equal(golden))
		})
		It("With missing module path", func() {
			createDirInTmpFolder("assembly-sample", "META-INF")
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources, "cf")
			checkError(err, wrongArtifactPathMsg, "java-hello-world")
		})
		It("With missing resource", func() {
			createDirInTmpFolder("assembly-sample", "META-INF")
			createDirInTmpFolder("assembly-sample", "web")
			createFileInTmpFolder("assembly-sample", "config-site-host.json")
			createDirInTmpFolder("assembly-sample", "inner")
			createFileInTmpFolder("assembly-sample", "inner", "uaa.json")
			createFileInTmpFolder("assembly-sample", "inner", "xs-security.json")
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources, "cf")
			checkError(err, unknownResourceContentTypeMsg, "java-uaa")

		})
		It("required resource with path fails when the path doesn't exist", func() {
			createDirInTmpFolder("assembly-sample", "META-INF")
			createDirInTmpFolder("assembly-sample", "web")
			createFileInTmpFolder("assembly-sample", "xs-security.json")
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources, "cf")
			// This fails because the config-site-host.json file (from the path of the required java-site-host) doesn't exist
			checkError(err, requiredEntriesProblemMsg, "java-hello-world-backend")
		})
		When("build-artifact-name is defined in the build parameters", func() {
			It("should take the defined build artifact name when the build artifact exists", func() {
				createDirInTmpFolder("mta", "node-js")
				createFileInTmpFolder("mta", "node-js", "data2.zip")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
				actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the archive.zip with the build artifact name when the build result is a folder", func() {
				createDirInTmpFolder("mta", "node-js")
				createFileInTmpFolder("mta", "node-js", "data2.zip")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
				actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should skip the module when it has no path", func() {
				createDirInTmpFolder("mta", "node-js", "data2")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifactNoPath.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
				actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_assembly_manifest_no_paths.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the build artifact name when the build result is also defined", func() {
				createDirInTmpFolder("mta", "node-js")
				createFileInTmpFolder("mta", "node-js", "ROOT.war")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildResultAndArtifact.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")).Should(Succeed())
				actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildResultAndArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should fail when build-artifact-name is not a string value", func() {
				createDirInTmpFolder("mta", "node-js")
				createFileInTmpFolder("mta", "node-js", "data.zip")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifactBad.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				err = setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")
				checkError(err, buildops.WrongBuildArtifactNameMsg, "1", "node-js")
			})
			It("should fail when data.zip exists instead of the build artifact name", func() {
				createDirInTmpFolder("mta", "node-js")
				createFileInTmpFolder("mta", "node-js", "data.zip")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				err = setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")
				checkError(err, wrongArtifactPathMsg, "node-js")
			})
			It("should fail when the build artifact doesn't exist in the module folder", func() {
				createDirInTmpFolder("mta", "node-js")
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				err = setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")
				checkError(err, wrongArtifactPathMsg, "node-js")
			})
			It("should fail when the module folder doesn't exist", func() {
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				err = setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{}, "cf")
				checkError(err, wrongArtifactPathMsg, "node-js")
			})
		})
	})

	var _ = Describe("genManifest", func() {
		It("Sanity", func() {
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			entries := []entry{
				{
					EntryName:   "node-js",
					EntryPath:   "node-js/data.zip",
					EntryType:   moduleEntry,
					ContentType: "application/zip",
				},
			}
			Ω(genManifest(loc.GetManifestPath(), entries)).Should(Succeed())
			actual := getFileContent(getFullPathInTmpFolder("mta", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			Ω(actual).Should(Equal(golden))
		})
		It("Fails on wrong location", func() {
			loc := dir.Loc{}
			Ω(genManifest(loc.GetManifestPath(), []entry{})).Should(HaveOccurred())
		})
		It("Fails on wrong version configuration", func() {
			versionCfg := version.VersionConfig
			version.VersionConfig = []byte(`
bad config
`)
			loc := dir.Loc{}
			Ω(genManifest(loc.GetManifestPath(), []entry{})).Should(HaveOccurred())
			version.VersionConfig = versionCfg
		})
	})
	var _ = Describe("moduleDefined", func() {
		It("not defined", func() {
			Ω(moduleDefined("x", []string{"a"})).Should(BeFalse())
		})
		It("empty list", func() {
			Ω(moduleDefined("x", []string{})).Should(BeTrue())
		})
		It("defined", func() {
			Ω(moduleDefined("x", []string{"y", "x"})).Should(BeTrue())
		})

	})

	var _ = Describe("populateManifest", func() {
		It("Nil file failure", func() {

			Ω(populateManifest(&testWriter{}, template.FuncMap{})).Should(HaveOccurred())
		})
	})

	var _ = Describe("buildEntries", func() {
		It("Sanity", func() {
			createDirInTmpFolder("result", "node-js")
			mod := mta.Module{Name: "module1"}
			requires := []mta.Requires{
				{
					Name: "req",
					Parameters: map[string]interface{}{
						"path": "node-js",
					},
				},
			}
			entries, err := buildEntries(&dir.Loc{SourcePath: getTestPath("result"), TargetPath: getTestPath("result")}, &mod, requires, &conttype.ContentTypes{})
			Ω(len(entries)).Should(Equal(1))
			Ω(err).Should(Succeed())
			e := entries[0]
			Ω(e.EntryType).Should(Equal(requiredEntry))
			Ω(e.EntryName).Should(Equal("module1/req"))
			Ω(e.EntryPath).Should(Equal("node-js"))
			Ω(e.ContentType).Should(Equal(dirContentType))
		})

	})

	DescribeTable("mergeDuplicateEntries", func(entries []entry, expected []entry) {
		actual := mergeDuplicateEntries(entries)
		Ω(actual).Should(Equal(expected))
	},
		Entry("returns non-module entries unchanged", []entry{
			{EntryName: "e1", EntryPath: "a", EntryType: resourceEntry, ContentType: "t"},
			{EntryName: "e2", EntryPath: "a", EntryType: requiredEntry, ContentType: "t"},
			{EntryName: "e3", EntryPath: "b", EntryType: resourceEntry, ContentType: "t"},
			{EntryName: "e4", EntryPath: "b", EntryType: resourceEntry, ContentType: "t"},
		}, []entry{
			{EntryName: "e2", EntryPath: "a", EntryType: requiredEntry, ContentType: "t"},
			{EntryName: "e1", EntryPath: "a", EntryType: resourceEntry, ContentType: "t"},
			{EntryName: "e3", EntryPath: "b", EntryType: resourceEntry, ContentType: "t"},
			{EntryName: "e4", EntryPath: "b", EntryType: resourceEntry, ContentType: "t"},
		},
		),
		Entry("merges module entries with the same path and keeps the modules order", []entry{
			{EntryName: "e1", EntryPath: "a", EntryType: moduleEntry, ContentType: "t1"},
			{EntryName: "e2", EntryPath: "a", EntryType: moduleEntry, ContentType: "t2"},
			{EntryName: "e3", EntryPath: "a", EntryType: moduleEntry, ContentType: "t3"},
			{EntryName: "e4", EntryPath: "b", EntryType: moduleEntry, ContentType: "t4"},
			{EntryName: "e5", EntryPath: "b", EntryType: moduleEntry, ContentType: "t5"},
			{EntryName: "e6", EntryPath: "c", EntryType: moduleEntry, ContentType: "t6"},
		}, []entry{
			{EntryName: "e1, e2, e3", EntryPath: "a", EntryType: moduleEntry, ContentType: "t1"},
			{EntryName: "e4, e5", EntryPath: "b", EntryType: moduleEntry, ContentType: "t4"},
			{EntryName: "e6", EntryPath: "c", EntryType: moduleEntry, ContentType: "t6"},
		},
		),
		Entry("merges required entries with the same path and keeps the entries order", []entry{
			{EntryName: "e1", EntryPath: "a", EntryType: requiredEntry, ContentType: "t1"},
			{EntryName: "e2", EntryPath: "a", EntryType: requiredEntry, ContentType: "t2"},
			{EntryName: "e3", EntryPath: "a", EntryType: requiredEntry, ContentType: "t3"},
			{EntryName: "e4", EntryPath: "b", EntryType: requiredEntry, ContentType: "t4"},
			{EntryName: "e5", EntryPath: "b", EntryType: requiredEntry, ContentType: "t5"},
			{EntryName: "e6", EntryPath: "c", EntryType: requiredEntry, ContentType: "t6"},
		}, []entry{
			{EntryName: "e1, e2, e3", EntryPath: "a", EntryType: requiredEntry, ContentType: "t1"},
			{EntryName: "e4, e5", EntryPath: "b", EntryType: requiredEntry, ContentType: "t4"},
			{EntryName: "e6", EntryPath: "c", EntryType: requiredEntry, ContentType: "t6"},
		},
		),
		Entry("merges module entries and keeps non-module entries unchanged at the end", []entry{
			{EntryName: "e1", EntryPath: "a", EntryType: resourceEntry, ContentType: "t1"},
			{EntryName: "e2", EntryPath: "a", EntryType: moduleEntry, ContentType: "t2"},
			{EntryName: "e3", EntryPath: "a", EntryType: moduleEntry, ContentType: "t3"},
			{EntryName: "e4", EntryPath: "b", EntryType: moduleEntry, ContentType: "t4"},
			{EntryName: "e5", EntryPath: "b", EntryType: requiredEntry, ContentType: "t5"},
		}, []entry{
			{EntryName: "e2, e3", EntryPath: "a", EntryType: moduleEntry, ContentType: "t2"},
			{EntryName: "e4", EntryPath: "b", EntryType: moduleEntry, ContentType: "t4"},
			{EntryName: "e5", EntryPath: "b", EntryType: requiredEntry, ContentType: "t5"},
			{EntryName: "e1", EntryPath: "a", EntryType: resourceEntry, ContentType: "t1"},
		},
		),
		Entry("merges module&required entries and keeps the other entries unchanged at the end", []entry{
			{EntryName: "e1", EntryPath: "a", EntryType: resourceEntry, ContentType: "t1"},
			{EntryName: "e2", EntryPath: "a", EntryType: moduleEntry, ContentType: "t2"},
			{EntryName: "e3", EntryPath: "a", EntryType: moduleEntry, ContentType: "t3"},
			{EntryName: "e4", EntryPath: "b", EntryType: moduleEntry, ContentType: "t4"},
			{EntryName: "e5", EntryPath: "b", EntryType: requiredEntry, ContentType: "t5"},
			{EntryName: "e6", EntryPath: "b", EntryType: requiredEntry, ContentType: "t6"},
			{EntryName: "e7", EntryPath: "c", EntryType: requiredEntry, ContentType: "t7"},
		}, []entry{
			{EntryName: "e2, e3", EntryPath: "a", EntryType: moduleEntry, ContentType: "t2"},
			{EntryName: "e4", EntryPath: "b", EntryType: moduleEntry, ContentType: "t4"},
			{EntryName: "e5, e6", EntryPath: "b", EntryType: requiredEntry, ContentType: "t5"},
			{EntryName: "e7", EntryPath: "c", EntryType: requiredEntry, ContentType: "t7"},
			{EntryName: "e1", EntryPath: "a", EntryType: resourceEntry, ContentType: "t1"},
		},
		),
	)

	var _ = Describe("getContentType", func() {
		It("fails on empty path", func() {
			_, err := getContentType("", &conttype.ContentTypes{})
			Ω(err).Should(HaveOccurred())
		})
	})
})

type testWriter struct {
}

func (t *testWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func getFileContentWithCliVersion(path string) string {
	content := getFileContent(path)
	v, _ := version.GetVersion()
	content = strings.Replace(content, "{{cli_version}}", v.CliVersion, -1)
	return content
}
