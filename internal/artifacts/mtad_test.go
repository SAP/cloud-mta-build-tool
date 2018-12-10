package artifacts

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("Mtad", func() {

	AfterEach(func() {
		os.RemoveAll(getTestPath("result"))
	})

	var _ = Describe("ExecuteGenMtad", func() {
		It("Sanity", func() {
			Ω(ExecuteGenMtad(getTestPath("mta"), getTestPath("result"), "dev", "cf", os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", "mta", "META-INF", "mtad.yaml")).Should(BeAnExistingFile())
		})
		It("Fails on location initialization", func() {
			Ω(ExecuteGenMtad("", getTestPath("result"), "dev", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})
		It("Fails on wrong source path - parse fails", func() {
			Ω(ExecuteGenMtad(getTestPath("mtax"), getTestPath("result"), "dev", "cf", os.Getwd)).Should(HaveOccurred())
		})
		It("Fails on broken extension file - parse ext fails", func() {
			Ω(ExecuteGenMtad(getTestPath("mtaWithBrokenExt"), getTestPath("result"), "dev", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("genMtad", func() {

		It("Fails on META folder creation", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			metaPath := ep.GetMetaPath()
			tmpDir := ep.GetTargetTmpDir()
			os.MkdirAll(tmpDir, os.ModePerm)
			_, err := os.Create(metaPath)
			Ω(err).Should(Succeed())
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(genMtad(mtaStr, &ep, ep.IsDeploymentDescriptor(), "cf")).Should(HaveOccurred())
		})
	})

	It("adaptMtadForDeployment", func() {
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
					Path: "app2",
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
		adaptMtadForDeployment(&mta, "neo")
		Ω(len(mta.Modules)).Should(Equal(1))
		Ω(mta.Modules[0].Name).Should(Equal("htmlapp2"))
		Ω(mta.Parameters["hcp-deployer-version"]).ShouldNot(BeNil())
	})

})
