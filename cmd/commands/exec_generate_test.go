package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
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

})

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
