package commands

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
)

var _ = Describe("cp command", func() {
	BeforeEach(func() {
		copyCmdSrc = getTestPath("mtahtml5")
		copyCmdTrg = getTestPath("result")
		copyCmdPatterns = []string{}
		err := dir.CreateDirIfNotExist(copyCmdTrg)
		Ω(err).Should(Succeed())
	})

	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
	})

	It("copy - Sanity", func() {
		copyCmdPatterns = []string{"mta.*", "ui5app2"}
		Ω(copyCmd.RunE(nil, []string{})).Should(Succeed())
		validateFilesInDir(getTestPath("result"), []string{"mta.sh", "mta.yaml", "ui5app2/", "ui5app2/test.txt"})
	})
	It("copy should not fail when the pattern is valid but doesn't match any files", func() {
		copyCmdPatterns = []string{"xxx.xxx"}
		Ω(copyCmd.RunE(nil, []string{})).Should(Succeed())
		validateFilesInDir(getTestPath("result"), []string{})
	})
	It("copy should not fail when there are no patterns", func() {
		copyCmdPatterns = []string{}
		Ω(copyCmd.RunE(nil, []string{})).Should(Succeed())
		validateFilesInDir(getTestPath("result"), []string{})
	})
	It("copy should return error when the source folder doesn't exist", func() {
		copyCmdSrc = getTestPath("mtahtml6")
		copyCmdPatterns = []string{"*"}
		Ω(copyCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})
	It("copy should return error when a pattern is invalid", func() {
		copyCmdPatterns = []string{"["}
		Ω(copyCmd.RunE(nil, []string{})).Should(HaveOccurred())
	})
})

// Check the folder exists and includes exactly the expected files
func validateFilesInDir(src string, expectedFilesInDir []string) {
	// List all files in the folder recursively
	filesInDir := make([]string, 0)
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Don't include the folder itself
		if filepath.Clean(path) == filepath.Clean(src) {
			return nil
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		if info.IsDir() {
			relPath += "/"
		}
		filesInDir = append(filesInDir, relPath)
		return nil
	})
	Ω(err).Should(Succeed())

	for _, expectedFile := range expectedFilesInDir {
		Ω(contains(expectedFile, filesInDir)).Should(BeTrue(), fmt.Sprintf("expected %s to be in the directory; directory contains %v", expectedFile, filesInDir))
	}
	for _, existingFile := range filesInDir {
		Ω(contains(existingFile, expectedFilesInDir)).Should(BeTrue(), fmt.Sprintf("did not expect %s to be in the directory; directory contains %v", existingFile, filesInDir))
	}
}

func contains(element string, elements []string) bool {
	for _, el := range elements {
		if el == element {
			return true
		}
	}
	return false
}