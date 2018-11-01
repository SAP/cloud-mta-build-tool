package builders_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuilders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builders Suite")
}
