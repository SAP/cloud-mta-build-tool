package artifacts

import (
	"fmt"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/SAP/cloud-mta/mta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/conttype"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"strings"
)

var _ = Describe("manifest", func() {

	BeforeEach(func() {
		Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "META-INF"))).Should(Succeed())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})

	var _ = Describe("setManifestDesc", func() {
		It("Sanity", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "config-site-host.json"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("Unknown content type, assembly scenario", func() {
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "server.js"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), Descriptor: dir.Dep}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, []*mta.Resource{})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`content type for the ".js" extension is not defined`))
		})
		It("Sanity - with configuration provided", func() {
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "config-site-host.json"))
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "xs-security.json"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_cfg.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest_cfg.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("wrong Commands configuration", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			moduleConf := commands.ModuleTypeConfig
			commands.ModuleTypeConfig = []byte("bad module conf")
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(HaveOccurred())
			commands.ModuleTypeConfig = moduleConf
		})
		It("module with defined build-result fails when build-result file does not exist in source directory", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "some.war"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(HaveOccurred())
		})
		It("module with defined build-result fails when build-result file does not exist in target temp directory", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult2.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(HaveOccurred())
		})
		It("entry for module with defined build-result has the build-result file", func() {
			Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data1.zip"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildResult.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifestBuildResult.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("wrong content types configuration", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			dirC, _ := ioutil.ReadDir(getTestPath("result", ".mta_mta_build_tmp"))
			for _, c := range dirC {
				fmt.Println(c.Name())
			}
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			contentTypesOrig := conttype.ContentTypeConfig
			conttype.ContentTypeConfig = []byte(`wrong configuraion`)
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(HaveOccurred())
			conttype.ContentTypeConfig = contentTypesOrig

		})
		It("Sanity - module with no path is not added to the manifest", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_no_paths.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_assembly_manifest_no_paths.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("With resources", func() {
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF"))).Should(Succeed())
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".assembly-sample_mta_build_tmp", "web"))).Should(Succeed())
			createTmpFile(getTestPath("result", ".assembly-sample_mta_build_tmp", "config-site-host.json"))
			createTmpFile(getTestPath("result", ".assembly-sample_mta_build_tmp", "xs-security.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources)).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_assembly_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("With missing module path", func() {
			Ω(dir.CreateDirIfNotExist(getTestPath("result", "assembly-sample_mta_build_tmp", "META-INF"))).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(wrongArtifactPathMsg, "java-hello-world")))
		})
		It("With missing resource", func() {
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF"))).Should(Succeed())
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".assembly-sample_mta_build_tmp", "web"))).Should(Succeed())
			createTmpFile(getTestPath("result", ".assembly-sample_mta_build_tmp", "config-site-host.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(unknownResourceContentTypeMsg, "java-uaa")))

		})
		It("required resource with path fails when the path doesn't exist", func() {
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF"))).Should(Succeed())
			Ω(dir.CreateDirIfNotExist(getTestPath("result", ".assembly-sample_mta_build_tmp", "web"))).Should(Succeed())
			createTmpFile(getTestPath("result", ".assembly-sample_mta_build_tmp", "xs-security.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, &loc, true, mtaObj.Modules, mtaObj.Resources)
			Ω(err).Should(HaveOccurred())
			// This fails because the config-site-host.json file (from the path of the required java-site-host) doesn't exist
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(requiredEntriesProblemMsg, "java-hello-world-backend")))
		})
		When("build-artifact-name is defined in the build parameters", func() {
			It("should take the defined build artifact name when the build artifact exists", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
				createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data2.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the archive.zip with the build artifact name when the build result is a folder", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
				createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data2.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should skip module with no path the build artifact name when the build result is the mta root folder", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data2"))).Should(Succeed())
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifactNoPath.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_assembly_manifest_no_paths.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the build artifact name when the build result is also defined", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
				createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "ROOT.war"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildResultAndArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				Ω(setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildResultAndArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should fail when build-artifact-name is not a string value", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
				createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifactBad.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(buildops.WrongBuildArtifactNameMsg, "1", "node-js")))
			})
			It("should fail when data.zip exists instead of the build artifact name", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
				createTmpFile(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(wrongArtifactPathMsg, "node-js")))
			})
			It("should fail when the build artifact doesn't exist in the module folder", func() {
				Ω(dir.CreateDirIfNotExist(getTestPath("result", ".mta_mta_build_tmp", "node-js"))).Should(Succeed())
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(wrongArtifactPathMsg, "node-js")))
			})
			It("should fail when the module folder doesn't exist", func() {
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, &loc, false, mtaObj.Modules, []*mta.Resource{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(wrongArtifactPathMsg, "node-js")))
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
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
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
			err := dir.CreateDirIfNotExist(getTestPath("result", ".result_mta_build_tmp", "node-js"))
			Ω(err).Should(Succeed())
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
