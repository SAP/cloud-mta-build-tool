package dir

import (
	"testing"

	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFsys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fsys Suite")
}

var _ = BeforeSuite(func() {
	logs.NewLogger()
})
