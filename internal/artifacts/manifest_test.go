package artifacts

import (
	"fmt"
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
		Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "META-INF"), os.ModePerm)).Should(Succeed())
	})

	AfterEach(func() {
		err1 := os.RemoveAll(getTestPath("result"))
		err2 := os.RemoveAll(getTestPath("result1"))
		err3 := os.RemoveAll(getTestPath("result2"))
		Ω(err1).Should(Succeed())
		Ω(err2).Should(Succeed())
		Ω(err3).Should(Succeed())
	})

	var _ = Describe("setManifestDesc", func() {
		It("Sanity", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			create(getTestPath("result", ".mta_mta_build_tmp", "config-site-host.json"))
			dirC, _ := ioutil.ReadDir(getTestPath("result", ".mta_mta_build_tmp"))
			for _, c := range dirC {
				fmt.Println(c.Name())
			}
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("Sanity - with configuration provided", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			create(getTestPath("result", ".mta_mta_build_tmp", "config-site-host.json"))
			create(getTestPath("result", ".mta_mta_build_tmp", "xs-security.json"))
			dirC, _ := ioutil.ReadDir(getTestPath("result", ".mta_mta_build_tmp"))
			for _, c := range dirC {
				fmt.Println(c.Name())
			}
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_cfg.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
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
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			moduleConf := commands.ModuleTypeConfig
			commands.ModuleTypeConfig = []byte("bad module conf")
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(HaveOccurred())
			commands.ModuleTypeConfig = moduleConf
		})
		It("module with defined build-result fails when build-result file does not exist in source directory", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "some.war"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(HaveOccurred())
		})
		It("module with defined build-result fails when build-result file does not exist in target temp directory", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult2.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(HaveOccurred())
		})
		It("entry for module with defined build-result has the build-result file", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data1.zip"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildResult.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifestBuildResult.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("Sanity - with list of modules provided; second module ignored", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			dirC, _ := ioutil.ReadDir(getTestPath("result", ".mta_mta_build_tmp"))
			for _, c := range dirC {
				fmt.Println(c.Name())
			}
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_2modules.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{"node-js"})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("wrong content types configuration", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			dirC, _ := ioutil.ReadDir(getTestPath("result", ".mta_mta_build_tmp"))
			for _, c := range dirC {
				fmt.Println(c.Name())
			}
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			contentTypesOrig := conttype.ContentTypeConfig
			conttype.ContentTypeConfig = []byte(`wrong configuraion`)
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(HaveOccurred())
			conttype.ContentTypeConfig = contentTypesOrig

		})
		It("Sanity - module with no path is not added to the manifest", func() {
			Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_no_paths.yaml"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_assembly_manifest_no_paths.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("With resources", func() {
			Ω(os.MkdirAll(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)).Should(Succeed())
			Ω(os.MkdirAll(getTestPath("result", ".assembly-sample_mta_build_tmp", "web"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result", ".assembly-sample_mta_build_tmp", "config-site-host.json"))
			create(getTestPath("result", ".assembly-sample_mta_build_tmp", "xs-security.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})).Should(Succeed())
			actual := getFileContent(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF", "MANIFEST.MF"))
			golden := getFileContent(getTestPath("golden_assembly_manifest.mf"))
			v, _ := version.GetVersion()
			golden = strings.Replace(golden, "{{cli_version}}", v.CliVersion, -1)
			fmt.Println(actual)
			fmt.Println(golden)
			Ω(actual).Should(Equal(golden))
		})
		It("With missing module path", func() {
			Ω(os.MkdirAll(getTestPath("result", "assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "java-hello-world" module content type`))
		})
		It("With missing resource", func() {
			Ω(os.MkdirAll(getTestPath("result1", ".assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)).Should(Succeed())
			Ω(os.MkdirAll(getTestPath("result1", ".assembly-sample_mta_build_tmp", "web"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result1", ".assembly-sample_mta_build_tmp", "config-site-host.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result1"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "java-uaa" resource content type`))

		})
		It("required resource with path fails when the path doesn't exist", func() {
			Ω(os.MkdirAll(getTestPath("result2", ".assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)).Should(Succeed())
			Ω(os.MkdirAll(getTestPath("result2", ".assembly-sample_mta_build_tmp", "web"), os.ModePerm)).Should(Succeed())
			create(getTestPath("result2", ".assembly-sample_mta_build_tmp", "xs-security.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result2"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})
			Ω(err).Should(HaveOccurred())
			// This fails because the config-site-host.json file (from the path of the required java-site-host) doesn't exist
			Ω(err.Error()).Should(
				ContainSubstring(`failed to generate the manifest file when building the required entries of the "java-hello-world-backend" module`))
		})
		When("build-artifact-name is defined in the build parameters", func() {
			It("should take the defined build artifact name when the build artifact exists", func() {
				Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
				create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data2.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the folder with the build artifact name when the build result is a folder", func() {
				Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data2"), os.ModePerm)).Should(Succeed())
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifactFolder.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the folder with the build artifact name when the build result is the mta root folder", func() {
				Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data2"), os.ModePerm)).Should(Succeed())
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifactNoPath.yaml"}
				mtaObj, err := loc.ParseFile()
				Ω(err).Should(Succeed())
				Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildArtifactNoPath.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should take the build artifact name when the build result is also defined", func() {
				Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
				create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "ROOT.war"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildResultAndArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(Succeed())
				actual := getFileContent(getTestPath("result", ".mta_mta_build_tmp", "META-INF", "MANIFEST.MF"))
				golden := getFileContentWithCliVersion(getTestPath("golden_manifestBuildResultAndArtifact.mf"))
				Ω(actual).Should(Equal(golden))
			})
			It("should fail when build-artifact-name is not a string value", func() {
				Ω(os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
				create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifactBad.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(`the node-js module has a non-string build-artifact-name in its build parameters`))
			})
			It("should fail when data.zip exists instead of the build artifact name", func() {
				Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
				create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "node-js" module content type`))
			})
			It("should fail when the build artifact doesn't exist in the module folder", func() {
				Ω(os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "node-js" module content type`))
			})
			It("should fail when the module folder doesn't exist", func() {
				loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaBuildArtifact.yaml"}
				mtaObj, _ := loc.ParseFile()
				err := setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "node-js" module content type`))
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
			loc := testTargetPathGetter{}
			mod := mta.Module{Name: "module1"}
			requires := []mta.Requires{
				{
					Name: "req",
					Parameters: map[string]interface{}{
						"path": "node-js",
					},
				},
			}
			entries, err := buildEntries(loc, &mod, requires, &conttype.ContentTypes{})
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

type testTargetPathGetter struct {
}

func (testTargetPathGetter) GetTarget() string {
	return getTestPath("mta")
}
func (testTargetPathGetter) GetTargetTmpDir() string {
	return getTestPath("mta")
}

type testWriter struct {
}

func (t *testWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func create(path string) {
	file, err := os.Create(path)
	Ω(err).Should(Succeed())
	Ω(file.Close()).Should(Succeed())
}

func getFileContentWithCliVersion(path string) string {
	content := getFileContent(path)
	v, _ := version.GetVersion()
	content = strings.Replace(content, "{{cli_version}}", v.CliVersion, -1)
	return content
}
