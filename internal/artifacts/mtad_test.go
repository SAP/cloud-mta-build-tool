package artifacts

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Mtad", func() {

	BeforeEach(func() {
		Ω(dir.CreateDirIfNotExist(getTestPath("result"))).Should(Succeed())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})

	var _ = Describe("ExecuteMtadGen", func() {
		It("Sanity", func() {
			Ω(ExecuteMtadGen(getTestPath("mta"), getTestPath("result"), nil, "cf", os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", "mtad.yaml")).Should(BeAnExistingFile())
		})
		It("Fails on creating META-INF folder", func() {
			createDirInTmpFolder("mta")
			file, err := os.Create(getFullPathInTmpFolder("mta", "META-INF"))
			Ω(err).Should(Succeed())
			Ω(file.Close()).Should(Succeed())
			Ω(ExecuteMtadGen(getTestPath("mta", "META-INF"), getTestPath("result"), nil, "cf", os.Getwd)).Should(HaveOccurred())
		})
		It("Fails on location initialization", func() {
			Ω(ExecuteMtadGen("", getTestPath("result"), nil, "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Fails on platform validation", func() {
			Ω(ExecuteMtadGen(getTestPath("mta"), getTestPath("result"), nil, "ab", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Fails on broken extension file - parse ext fails", func() {
			Ω(ExecuteMtadGen(getTestPath("mtaWithBrokenExt"), getTestPath("result"), []string{"cf-mtaext.yaml"}, "cf", os.Getwd)).Should(HaveOccurred())
		})
		It("Fails on wrong source path - parse fails", func() {
			Ω(ExecuteMtadGen(getTestPath("mtax"), getTestPath("result"), nil, "cf", os.Getwd)).Should(HaveOccurred())
		})
		It("Fails on broken platforms configuration", func() {
			cfg := platform.PlatformConfig
			platform.PlatformConfig = []byte("abc abc")
			Ω(ExecuteMtadGen(getTestPath("mta"), getTestPath("result"), nil, "cf", os.Getwd)).Should(HaveOccurred())
			platform.PlatformConfig = cfg
		})
	})

	var _ = Describe("genMtad", func() {

		It("Fails on META folder creation", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			metaPath := ep.GetMetaPath()
			tmpDir := ep.GetTargetTmpDir()
			Ω(dir.CreateDirIfNotExist(tmpDir)).Should(Succeed())
			file, err := os.Create(metaPath)
			Ω(err).Should(Succeed())
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(genMtad(mtaStr, &ep, ep.IsDeploymentDescriptor(), "cf", yaml.Marshal)).Should(HaveOccurred())
			Ω(file.Close()).Should(Succeed())
		})
		It("Fails on mtad marshalling", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(genMtad(mtaStr, &ep, ep.IsDeploymentDescriptor(), "cf", func(i interface{}) (out []byte, err error) {
				return nil, errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Fails on mtad schema version change", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result"), MtaFilename: "mtaBadSchemaVersion.yaml"}
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(genMtad(mtaStr, &ep, ep.IsDeploymentDescriptor(), "cf", yaml.Marshal)).Should(HaveOccurred())
		})
	})

})

var _ = Describe("adaptModulePath", func() {
	It("path by module name", func() {
		mod := mta.Module{Name: "htmlapp2", Path: "xyz"}
		Ω(adaptModulePath(&testMtadLoc{}, &mod, true)).Should(Succeed())
		Ω(mod.Path).Should(Equal("htmlapp2"))
	})
	It("Fails on location initialization - validatePaths is set to true", func() {
		Ω(adaptModulePath(&testMtadLoc{}, &mta.Module{Name: "htmlapp", Path: "xxx"}, true)).Should(HaveOccurred())
	})
})

var _ = Describe("removeUndeployedModules", func() {
	It("Sanity", func() {
		mta := mta.MTA{
			ID:      "mta_proj",
			Version: "1.0.0",
			Modules: []*mta.Module{
				{
					Name: "htmlapp",
					Type: "javascript.nodejs",
					Path: "app",
					BuildParams: map[string]interface{}{
						buildops.SupportedPlatformsParam: []string{},
					},
				},
				{
					Name: "htmlapp2",
					Type: "javascript.nodejs",
					Path: "node-js1",
					BuildParams: map[string]interface{}{
						buildops.SupportedPlatformsParam: nil,
					},
				},
				{
					Name: "java",
					Type: "java.tomcat",
					Path: "app3",
					BuildParams: map[string]interface{}{
						buildops.SupportedPlatformsParam: []string{},
					},
				},
			},
		}
		removeUndeployedModules(&mta, "neo")
		Ω(len(mta.Modules)).Should(Equal(1))
		Ω(mta.Modules[0].Name).Should(Equal("htmlapp2"))
	})
})

var _ = Describe("setPlatformSpecificParameters", func() {
	It("Sanity", func() {
		mta := mta.MTA{
			ID:      "mta_proj",
			Version: "1.0.0",
			Modules: []*mta.Module{
				{
					Name: "htmlapp",
					Type: "javascript.nodejs",
					Path: "app",
					BuildParams: map[string]interface{}{
						buildops.SupportedPlatformsParam: []string{},
					},
				},
				{
					Name: "my-html-app",
					Type: "javascript.nodejs",
					Path: "node-js1",
					Parameters: map[string]interface{}{
						"a": "b",
					},
				},
				{
					Name: "java",
					Type: "java.tomcat",
					Path: "app3",
					Parameters: map[string]interface{}{
						"name": "someName",
					},
				},
			},
		}

		setPlatformSpecificParameters(&mta, "neo")

		Ω(mta.Parameters["hcp-deployer-version"]).ShouldNot(BeNil())

		Ω(len(mta.Modules)).Should(Equal(3))
		Ω(mta.Modules[0].Parameters).ShouldNot(BeNil())
		// Check it sets the name when parameters don't exist
		Ω(mta.Modules[0].Parameters["name"]).Should(Equal("htmlapp"))
		Ω(mta.Modules[1].Parameters).ShouldNot(BeNil())
		// Check it sets a supported name based on the module name
		Ω(mta.Modules[1].Parameters["name"]).Should(Equal("myhtmlapp"))
		// Check it doesn't change pre-set names
		Ω(mta.Modules[2].Parameters).ShouldNot(BeNil())
		Ω(mta.Modules[2].Parameters["name"]).Should(Equal("someName"))
	})
})

var _ = DescribeTable("adjustNeoAppName", func(name string, expected string) {
	Ω(adjustNeoAppName(name)).Should(Equal(expected))
},
	Entry("supported name", "abc123cde", "abc123cde"),
	Entry("name with uppercase letters", "Abc123P", "abc123p"),
	Entry("name starts with numbers", "87xaa8s", "xaa8s"),
	Entry("name starts with numbers and uppercase", "87Xaa8s", "xaa8s"),
	Entry("name with unsupported characters", "my-module", "mymodule"),
	Entry("name starts with unsupported characters and numbers", "!~11-2ave_1gne", "ave1gne"),
	Entry("name longer than 30 characters before removing unsupported characters", "a123456789-12345678901234567890", "a12345678912345678901234567890"),
	Entry("name longer than 30 characters after removing unsupported characters", "a1234567890-12345678901234567890", "a12345678901234567890123456789"),
	Entry("mixed conditions", "1234567890abcdeABCDE--==?12345mymodule", "abcdeabcde12345mymodule"),
)

type testMtadLoc struct {
}

func (loc *testMtadLoc) GetTarget() string {
	return getTestPath("mta")
}
func (loc *testMtadLoc) GetTargetTmpDir() string {
	return getTestPath("mta")
}

var _ = Describe("mtadLoc", func() {
	It("GetManifestPath", func() {
		loc := mtadLoc{"anyPath"}
		Ω(loc.GetManifestPath()).Should(Equal(""))
	})
	It("GetMtarDir", func() {
		loc := mtadLoc{"anyPath"}
		Ω(loc.GetMtarDir(true)).Should(Equal(""))
	})
	It("GetManifestPath", func() {
		loc := mtadLoc{"anyPath"}
		Ω(loc.GetMetaPath()).Should(Equal("anyPath"))
	})
})
