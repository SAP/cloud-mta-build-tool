package artifacts

import (
	"testing"

	"cloud-mta-build-tool/internal/logs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArtifacts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artifacts Suite")
}

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})
