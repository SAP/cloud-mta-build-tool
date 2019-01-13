package dir

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/logs"
)

func TestFsys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fsys Suite")
}

var _ = BeforeSuite(func() {
	logs.NewLogger()
})
