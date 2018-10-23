package commands

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func getFullPath(relPath ...string) string {
	path, _ := dir.GetFullPath(relPath...)
	return path
}

var _ = Describe("Sub Commands", func() {
	BeforeEach(func() {
		logs.Logger = logs.NewLogger()
	})

	readFileContent := func(args ...string) string {
		content, _ := ioutil.ReadFile(getFullPath(args...))
		contentString := string(content[:])
		contentString = strings.Replace(contentString, "\n", "", -1)
		contentString = strings.Replace(contentString, "\r", "", -1)
		return contentString
	}

	It("Generate Meta", func() {
		generateMeta(filepath.Join("testdata", "mtahtml5"), []string{filepath.Join("testdata", "result"), "testapp"})
		Ω(readFileContent("testdata", "result", "META-INF", "mtad.yaml")).Should(Equal(readFileContent("testdata", "golden", "mtad.yaml")))
		os.RemoveAll(getFullPath("testdata", "result"))
	})

	It("Generate Mtar", func() {
		generateMtar(filepath.Join("testdata", "mtahtml5"), []string{getFullPath("testdata", "mtahtml5"), getFullPath("testdata")})
		mtarPath := getFullPath("testdata", "mtahtml5.mtar")
		Ω(mtarPath).Should(BeAnExistingFile())
		os.RemoveAll(mtarPath)
	})

	Describe("Pack", func() {
		DescribeTable("Standard cases", func(args []string, validator func(projectPath string)) {
			pack.Run(nil, args)
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

	It("Target file in opened status", func() {
		var str bytes.Buffer
		// navigate log output to local string buffer. It will be used for error analysis
		logs.Logger.SetOutput(&str)
		f, _ := os.Create(filepath.Join("testdata", "temp"))

		args := []string{getFullPath("testdata", "temp"), filepath.Join("testdata", "mtahtml5", "testapp"), "ui5app"}

		pack.Run(nil, args)
		Ω(str.String()).Should(ContainSubstring("ERROR mkdir"))

		f.Close()
		// cleanup command used for test temp file removal
		cleanup.Run(nil, []string{filepath.Join("testdata", "temp")})
	})

	var _ = DescribeTable("Generate commands call with no effect", func(projectRelPath string, module string,
		expectedFileRelPath string, genGommand *cobra.Command) {
		genGommand.Run(nil, []string{projectRelPath, module})
		Ω(getFullPath(projectRelPath, expectedFileRelPath)).ShouldNot(BeAnExistingFile())
	},
		Entry("Generate META", filepath.Join("testdata", "result"), "testapp", filepath.Join("META-INF", "mtad.yaml"), genMeta),
		Entry("Generate MTAR", filepath.Join("testdata", "mtahtml5"), "testdata", filepath.Join("testdata-INF", "mtahtml5.mtar"), genMtar),
	)

})
