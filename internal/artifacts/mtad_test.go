package artifacts

import (
	"os"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mtad", func() {

	var _ = Describe("GenMtad", func() {
		BeforeEach(func() {
			os.RemoveAll(getTestPath("result"))
		})
		It("Fails on META folder creation", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getTestPath("result")}
			metaPath, err := ep.GetMetaPath()
			Ω(err).Should(Succeed())
			tmpDir, err := ep.GetTargetTmpDir()
			Ω(err).Should(Succeed())
			os.MkdirAll(tmpDir, os.ModePerm)
			_, err = os.Create(metaPath)
			Ω(err).Should(Succeed())
			mtaBytes, err := dir.Read(&ep)
			Ω(err).Should(Succeed())
			mtaStr, err := mta.Unmarshal(mtaBytes)
			Ω(err).Should(Succeed())
			Ω(GenMtad(mtaStr, &ep, "cf", func(mtaStr *mta.MTA, platform string) {

			})).Should(HaveOccurred())
		})
	})

	It("AdaptMtadForDeployment", func() {
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
		AdaptMtadForDeployment(&mta, "cf")
		Ω(len(mta.Modules)).Should(Equal(1))
		Ω(mta.Modules[0].Name).Should(Equal("htmlapp2"))
	})

})
