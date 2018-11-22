package exec_test

import (
	"testing"

	"cloud-mta-build-tool/internal/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Exec Suite")
}

var _ = BeforeSuite(func() {
	logs.NewLogger()
})
