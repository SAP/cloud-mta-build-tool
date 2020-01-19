package conttype

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestContenttype(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contenttype Suite")
}
