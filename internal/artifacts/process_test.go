package artifacts

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fsys"
)

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}

var _ = Describe("Process", func() {
	It("Sanity", func() {
		ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module")}
		Ω(processMta("empty", &ep, []string{}, func(file []byte, args []string) error { return nil })).Should(Succeed())
	})
	It("Sanity", func() {
		ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module"), MtaFilename: "unknown.yaml"}
		Ω(processMta("empty", &ep, []string{}, func(file []byte, args []string) error { return nil })).Should(HaveOccurred())
	})
})
