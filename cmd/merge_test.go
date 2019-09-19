package commands

import (
	"github.com/SAP/cloud-mta/mta"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"path/filepath"
)

var _ = Describe("Merge commands call", func() {
	AfterEach(func() {
		Ω(os.RemoveAll(getTestPath("result", "result.yaml"))).Should(Succeed())
	})
	It("merges with one extension", func() {
		path := getTestPath("mtaext")
		mergeCmdSrc = path
		mergeCmdTrg = getTestPath("result")
		mergeCmdExtensions = []string{"ext.mtaext"}
		mergeCmdName = "result.yaml"

		Ω(mergeCmd.RunE(nil, []string{})).Should(Succeed())

		mtadPath := filepath.Join(mergeCmdTrg, "result.yaml")
		Ω(mtadPath).Should(BeAnExistingFile())
		content, e := ioutil.ReadFile(mtadPath)
		Ω(e).Should(Succeed())
		mtaObj, e := mta.Unmarshal(content)
		Ω(e).Should(Succeed())
		expected, e := mta.Unmarshal([]byte(`
ID: mta_demo
_schema-version: '2.1'
version: 0.0.1

modules:
- name: node
  type: nodejs
  path: node
  provides:
  - name: node_api
    properties:
      url: ${default-url}
  build-parameters:
    supported-platforms: [cf]
- name: node-js
  type: nodejs
  path: node-js
  provides:
  - name: node-js_api
    properties:
      url: ${default-url}
  build-parameters:
    builder: zip
    supported-platforms: [neo]
`))
		Ω(e).Should(Succeed())
		Ω(mtaObj).Should(Equal(expected))
	})
})
