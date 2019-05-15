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
		os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp", "META-INF"), os.ModePerm)
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
		os.RemoveAll(getTestPath("result1"))
		os.RemoveAll(getTestPath("result2"))
	})

	var _ = Describe("setManifestDesc", func() {
		It("Sanity", func() {
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "config-site-host.json"))
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
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "config-site-host.json"))
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "xs-security.json"))
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
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			moduleConf := commands.ModuleTypeConfig
			commands.ModuleTypeConfig = []byte("bad module conf")
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(HaveOccurred())
			commands.ModuleTypeConfig = moduleConf
		})
		It("wrong build results", func() {
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mtaWrongBuildResult.yaml"}
			mtaObj, _ := loc.ParseFile()
			Ω(setManifestDesc(&loc, &loc, mtaObj.Modules, []*mta.Resource{}, []string{})).Should(HaveOccurred())
		})
		It("correct build results", func() {
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			file, _ := os.Create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data1.zip"))
			file.Close()
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
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			file, _ := os.Create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
			file.Close()
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
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			os.Create(getTestPath("result", ".mta_mta_build_tmp", "node-js", "data.zip"))
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
		It("Sanity - no paths", func() {
			os.Mkdir(getTestPath("result", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
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
			os.MkdirAll(getTestPath("result", ".assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)
			os.MkdirAll(getTestPath("result", ".assembly-sample_mta_build_tmp", "web"), os.ModePerm)
			os.Create(getTestPath("result", ".assembly-sample_mta_build_tmp", "config-site-host.json"))
			os.Create(getTestPath("result", ".assembly-sample_mta_build_tmp", "xs-security.json"))
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
			os.MkdirAll(getTestPath("result", "assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getResultPath(), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "java-hello-world" module content type`))
		})
		It("With missing resource", func() {
			os.MkdirAll(getTestPath("result1", ".assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)
			os.MkdirAll(getTestPath("result1", ".assembly-sample_mta_build_tmp", "web"), os.ModePerm)
			os.Create(getTestPath("result1", ".assembly-sample_mta_build_tmp", "config-site-host.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result1"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(`failed to generate the manifest file when getting the "java-uaa" resource content type`))

		})
		It("With missing requirement", func() {
			os.MkdirAll(getTestPath("result2", ".assembly-sample_mta_build_tmp", "META-INF"), os.ModePerm)
			os.MkdirAll(getTestPath("result2", ".assembly-sample_mta_build_tmp", "web"), os.ModePerm)
			os.Create(getTestPath("result2", ".assembly-sample", "xs-security.json"))
			loc := dir.Loc{SourcePath: getTestPath("assembly-sample"), TargetPath: getTestPath("result2"), Descriptor: "dep"}
			mtaObj, err := loc.ParseFile()
			Ω(err).Should(Succeed())
			err = setManifestDesc(&loc, &loc, mtaObj.Modules, mtaObj.Resources, []string{})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(
				ContainSubstring(`failed to generate the manifest file when building the required entries of the "java-hello-world-backend" module`))

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
