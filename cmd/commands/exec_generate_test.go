package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
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
