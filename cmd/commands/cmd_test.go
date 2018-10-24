package commands

import (
	"cloud-mta-build-tool/internal/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cmd", func() {
	var _ = Describe("Version", func() {
		It("Sanity", func() {
			version.VersionConfig = []byte(`
cli_version: 0.0.1
makefile_version: 10.5.3
`)
			out := executeAndProvideOutput(func() {
				versionCmd.Run(nil, []string{})
			})
			Ω(out).Should(Equal("0.0.1\n"))
		})
	})

})
