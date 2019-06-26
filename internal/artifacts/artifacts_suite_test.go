package artifacts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

func TestArtifacts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artifacts Suite")
}

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}

func getResultPath() string {
	return getTestPath("result")
}

func removeSpecialSymbols(b []byte) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9._{}]+")
	s := string(b)
	fmt.Println(s)
	s = strings.Replace(s, "0xd, ", "", -1)
	s = reg.ReplaceAllString(s, "")
	return s
}

func getFileContent(filePath string) string {
	expected, _ := ioutil.ReadFile(filePath)
	return removeSpecialSymbols(expected)
}
