package commands

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/logs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pack", func() {
	BeforeEach(func() {
		logs.Logger = logs.NewLogger()
	})

	DescribeTable("Standard cases", func(args []string, validator func(projectPath string)) {
		packCmd.Run(nil, args)
		validator(args[0])
	},
		Entry("SanityTest",
			[]string{
				getFullPath("testdata", "result"),
				filepath.Join("testdata", "mtahtml5", "testapp"),
				"ui5app",
			},
			func(projectPath string) {
				resultPath := filepath.Join(projectPath, "ui5app", "data.zip")
				fileInfo, _ := os.Stat(resultPath)
				Ω(fileInfo).ShouldNot(BeNil())
				Ω(fileInfo.IsDir()).Should(BeFalse())
				os.RemoveAll(projectPath)
			}),
		Entry("Wrong relative path to module",
			[]string{
				getFullPath("testdata", "result"),
				filepath.Join("testdata", "mtahtml5", "ui5app"),
				"ui5app",
			},
			func(projectPath string) {
				fileInfo, _ := os.Stat(filepath.Join(projectPath, "ui5app", "data.zip"))
				Ω(fileInfo).Should(BeNil())
			}),
		Entry("Missing arguments",
			[]string{
				getFullPath("testdata", "result"),
				"ui5app",
			},
			func(projectPath string) {
				fileInfo, _ := os.Stat(filepath.Join(projectPath, "ui5app", "data.zip"))
				Ω(fileInfo).Should(BeNil())
			}),
	)
})
