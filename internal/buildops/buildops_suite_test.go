package buildops_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

func TestBuildops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buildops Suite")
}

var _ = BeforeSuite(func() {
	logs.NewLogger()
})
