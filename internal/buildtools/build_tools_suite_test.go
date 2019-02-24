package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildTools(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BuildTools Suite")
}
