package artifacts

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
)

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}

var _ = Describe("Process", func() {

})
