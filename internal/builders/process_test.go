package builders

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Process", func() {
	var buildCfg = []byte(`
version: 1
builders:
  - name: html5
    info: "build UI5 application"
    type:
    - command: npm install
    - command: grunt
    - command: npm prune --production
  - name: java
    info: "build java application"
    type:
    - command: mvn clean install
  - name: nodejs
    info: "build nodejs application"
    type:
    - command: npm install
  - name: golang
    info: "build golang application"
    type:
    - command: go build *.go
`)
	var builders = Builders{
		Version: "1",
		Builders: []builder{
			{
				Name: "html5",
				Info: "build UI5 application",
				Type: []Commands{
					{Command: "npm install"},
					{Command: "grunt"},
					{Command: "npm prune --production"},
				},
			},
			{
				Name: "java",
				Info: "build java application",
				Type: []Commands{
					{Command: "mvn clean install"},
				},
			},
			{
				Name: "nodejs",
				Info: "build nodejs application",
				Type: []Commands{
					{Command: "npm install"},
				},
			},
			{
				Name: "golang",
				Info: "build golang application",
				Type: []Commands{
					{Command: "go build *.go"},
				},
			},
		}}
	var malformedBuildCfg = []byte(`bad:  "YAML" syntax`)

	var _ = DescribeTable("Unmarshal", func(input []byte, expected Builders, match types.GomegaMatcher) {
		actual, err := parse(input)
		Ω(actual).Should(Equal(expected))
		Ω(err).Should(match)
	},
		Entry("Sanity", buildCfg, builders, Succeed()),
		Entry("MalformedCfg", malformedBuildCfg, Builders{}, HaveOccurred()))

})
