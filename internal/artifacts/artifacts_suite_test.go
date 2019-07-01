package artifacts

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
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
	s = strings.Replace(s, "0xd, ", "", -1)
	s = reg.ReplaceAllString(s, "")
	return s
}

func getFileContent(filePath string) string {
	expected, _ := ioutil.ReadFile(filePath)
	return removeSpecialSymbols(expected)
}

func createFileInTmpFolder(projectName string, path ...string) {
	file, err := os.Create(getFullPathInTmpFolder(projectName, path...))
	Ω(err).Should(Succeed())
	err = file.Close()
	Ω(err).Should(Succeed())
}

func createFileInGivenPath(path string) {
	file, err := os.Create(path)
	Ω(err).Should(Succeed())
	err = file.Close()
	Ω(err).Should(Succeed())
}

func createDirInTempFolder(projectName string, path ...string) {
	err := dir.CreateDirIfNotExist(getFullPathInTmpFolder(projectName, path...))
	Ω(err).Should(Succeed())
}

func getFullPathInTmpFolder(projectName string, path ...string) string {
	pathWithResultFolder := []string{"result", "." + projectName + "_mta_build_tmp"}
	pathWithResultFolder = append(pathWithResultFolder, path...)
	return getTestPath(pathWithResultFolder...)
}
