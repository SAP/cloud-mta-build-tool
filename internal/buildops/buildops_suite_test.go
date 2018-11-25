package buildops_test

import (
	"testing"

	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buildops Suite")
}

var _ = BeforeSuite(func() {
	logs.NewLogger()
})
