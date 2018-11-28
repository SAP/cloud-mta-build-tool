package validate_test

import (
	"testing"

	"cloud-mta-build-tool/internal/logs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validations Suite")
}

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})
