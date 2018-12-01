package artifacts

import (
	"cloud-mta-build-tool/mta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mtad", func() {

	It("CleanMtaForDeployment", func() {
		mta := mta.MTA{
			ID:      "mta_proj",
			Version: "1.0.0",
			Modules: []*mta.Module{
				{
					Name: "htmlapp",
					Type: "javascript.nodejs",
					Path: "app",
					BuildParams: mta.BuildParameters{
						SupportedPlatforms: []string{},
					},
				},
				{
					Name: "htmlapp2",
					Type: "javascript.nodejs",
					Path: "app2",
					BuildParams: mta.BuildParameters{
						SupportedPlatforms: nil,
					},
				},
				{
					Name: "java",
					Type: "java.tomcat",
					Path: "app3",
					BuildParams: mta.BuildParameters{
						SupportedPlatforms: []string{},
					},
				},
			},
		}
		CleanMtaForDeployment(&mta)
		Ω(len(mta.Modules)).Should(Equal(1))
		Ω(mta.Modules[0].Name).Should(Equal("htmlapp2"))
	})
})
