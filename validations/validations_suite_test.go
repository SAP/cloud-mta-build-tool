package validate_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validations Suite")
}
