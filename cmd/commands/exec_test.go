package commands

import (
	"bytes"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Commands", func() {

	var _ = Describe("Pack and cleanup commands", func() {
		It("Target file in opened status", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			f, _ := os.Create(filepath.Join("testdata", "temp"))

			args := []string{getFullPath("testdata", "temp"), filepath.Join("testdata", "mtahtml5", "testapp"), "ui5app"}

			packCmd.Run(nil, args)
			Ω(str.String()).Should(ContainSubstring("ERROR mkdir"))

			f.Close()
			// cleanup command used for test temp file removal
			cleanupCmd.Run(nil, []string{filepath.Join("testdata", "temp")})
		})
	})

	var _ = DescribeTable("Generate commands call with no effect", func(projectRelPath string, module string,
		expectedFileRelPath string, genGommand *cobra.Command) {
		genGommand.Run(nil, []string{projectRelPath, module})
		Ω(getFullPath(projectRelPath, expectedFileRelPath)).ShouldNot(BeAnExistingFile())
	},
		Entry("Generate META", filepath.Join("testdata", "result"), "testapp", filepath.Join("META-INF", "mtad.yaml"), genMetaCmd),
		Entry("Generate MTAR", filepath.Join("testdata", "mtahtml5"), "testdata", filepath.Join("testdata-INF", "mtahtml5.mtar"), genMtarCmd),
		Entry("Generate MTAD", filepath.Join("testdata", "mtahtml5"), "testdata", filepath.Join("testdata-INF", "mtahtml5.mtar"), genMtadCmd),
	)

	var _ = Describe("Validate", func() {
		It("Invalid yaml path", func() {
			var str bytes.Buffer
			// navigate log output to local string buffer. It will be used for error analysis
			logs.Logger.SetOutput(&str)
			validateCmd.Run(nil, []string{})

			Ω(str.String()).Should(ContainSubstring("ERROR MTA validation failed. Error reading the MTA file:"))
		})
	})

})
