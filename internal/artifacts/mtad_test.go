package artifacts

import (
	"os"

	"github.com/SAP/cloud-mta/mta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
)

var _ = Describe("Mtad", func() {

	BeforeEach(func() {
		os.MkdirAll(getTestPath("resultMtad"), os.ModePerm)
	})

	AfterEach(func() {
		os.RemoveAll(getTestPath("resultMtad"))
	})

	var _ = Describe("ExecuteGenMtad", func() {
		It("Sanity", func() {
			os.MkdirAll(getTestPath("resultMtad", ".mta_mta_build_tmp", "node-js"), os.ModePerm)
			Ω(ExecuteGenMtad(getTestPath("mta"), getTestPath("resultMtad"), "cf", os.Getwd)).Should(Succeed())
			Ω(getTestPath("resultMtad", "mtad.yaml")).Should(BeAnExistingFile())
		})
		It("Fails on location initialization", func() {
			Ω(ExecuteGenMtad("", getTestPath("resultMtad"), "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Fails on platform validation", func() {
			Ω(ExecuteGenMtad(getTestPath("mta"), getTestPath("resultMtad"), "ab", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Fails on wrong source path - parse fails", func() {
			Ω(ExecuteGenMtad(getTestPath("mtax"), getTestPath("resultMtad"), "cf", os.Getwd)).Should(HaveOccurred())
		})
		It("Fails on broken extension file - parse ext fails", func() {
			Ω(ExecuteGenMtad(getTestPath("mtaWithBrokenExt"), getTestPath("resultMtad"), "cf", os.Getwd)).Should(HaveOccurred())
		})
		It("Fails on broken platforms configuration", func() {
			cfg := platform.PlatformConfig
			platform.PlatformConfig = []byte("abc abc")
			Ω(ExecuteGenMtad(getTestPath("mta"), getTestPath("resultMtad"), "cf", os.Getwd)).Should(HaveOccurred())
			platform.PlatformConfig = cfg
		})
	})

	var _ = Describe("genMtad", func() {

		It("Fails on META folder creation", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("resultMtad")}
			metaPath := ep.GetMetaPath()
			tmpDir := ep.GetTargetTmpDir()
			os.MkdirAll(tmpDir, os.ModePerm)
			file, err := os.Create(metaPath)
			Ω(err).Should(Succeed())
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(genMtad(mtaStr, &ep, ep.IsDeploymentDescriptor(), "cf", yaml.Marshal)).Should(HaveOccurred())
			file.Close()
		})
		It("Fails on mtad marshalling", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("resultMtad")}
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(genMtad(mtaStr, &ep, ep.IsDeploymentDescriptor(), "cf", func(i interface{}) (out []byte, err error) {
				return nil, errors.New("err")
			})).Should(HaveOccurred())
		})
	})

})

var _ = Describe("adaptMtadForDeployment", func() {
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
		err := adaptMtadForDeployment(&testMtadLoc{}, &mta, "neo")
		Ω(err).Should(Succeed())
		Ω(len(mta.Modules)).Should(Equal(1))
		Ω(mta.Modules[0].Name).Should(Equal("htmlapp2"))
		Ω(mta.Modules[0].Path).Should(Equal("htmlapp2"))
		Ω(mta.Parameters["hcp-deployer-version"]).ShouldNot(BeNil())
	})
})

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
})
