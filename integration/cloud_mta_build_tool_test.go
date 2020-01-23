// +build integration

package integration_test

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta/mta"
)

const (
	demoArchiveName = "mta_demo_0.0.1.mtar"
	//javaArchiveName = "com.fetcher.project_0.0.1.mtar"
	binPath = "mbt"
)

var _ = Describe("Integration - CloudMtaBuildTool", func() {

	var mbtName = ""

	BeforeSuite(func() {
		By("Building MBT")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			mbtName = "mbt"
		} else {
			mbtName = "mbt.exe"
		}
		// This runs locally for testing purpose only
		/* #nosec */
		cmd := exec.Command("go", "build", "-o", filepath.Join(os.Getenv("GOPATH"), "/bin/"+mbtName), ".")
		cmd.Dir = filepath.FromSlash("../")
		err := cmd.Run()
		if err != nil {
			fmt.Println("binary creation failed: ", err)
		}
	})

	AfterSuite(func() {
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_demo/" + mbtName))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_demo/Makefile.mta"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_demo/mtad.yaml"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_demo/abc.mtar"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_demo/mta_archives"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_java/myModule/target"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_java/Makefile.mta"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_java/mtad.yaml"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_java/mta_archives"))).Should(Succeed())
		Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_demo/node/package-lock.json"))).Should(Succeed())
	})

	var _ = Describe("Command to provide the list of modules", func() {

		It("Getting module", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash(binPath)
			cmdOut, cmdErr, err := execute(bin, "provide modules", path)
			Ω(err).Should(Succeed(), cmdErr)
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(ContainSubstring("[node node-js]" + "\n"))
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/")
			bin := filepath.FromSlash(binPath)
			_, _, err := execute(bin, "provide modules 2", path)
			Ω(err).Should(HaveOccurred())
		})
	})
	var _ = Describe("Generate the Makefile according to the mta.yaml file", func() {

		It("Generate Makefile for mta_demo", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_demo")
			bin := filepath.FromSlash(binPath)
			_, errOut, err := execute(bin, "init", path)
			Ω(err).Should(Succeed(), errOut)

			// Check the MakeFile was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta")).Should(BeAnExistingFile())
		})

		It("Generate Makefile for mta_java", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_java")
			bin := filepath.FromSlash(binPath)
			_, errOut, err := execute(bin, "init", path)
			Ω(err).Should(Succeed(), errOut)

			// Check the MakeFile was generated
			Ω(filepath.Join(dir, "testdata", "mta_java", "Makefile.mta")).Should(BeAnExistingFile())
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash(binPath)
			_, _, err := execute(bin, "init 2", path)
			Ω(err).Should(HaveOccurred())

		})
	})

	var _ = Describe("Generate MTAR", func() {
		It("Generate MTAR with provided target and mtar name", func() {
			dir, _ := os.Getwd()
			Ω(os.RemoveAll(filepath.Join(dir, "testdata", "mta_demo", demoArchiveName))).Should(Succeed())
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("make")
			cmdOut, errOut, err := execute(bin, "-f Makefile.mta p=cf mtar=abc t="+path, path)
			Ω(err).Should(Succeed(), errOut)
			Ω(cmdOut).ShouldNot(BeEmpty())
			// Check the archive was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "abc.mtar")).Should(BeAnExistingFile())
		})

		It("Generate MTAR - wrong platform", func() {

			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("make")
			out, errOut, err := execute(bin, "-f Makefile.mta p=xxx mtar=xyz1", path)
			Ω(err).Should(HaveOccurred())
			Ω(out + errOut).Should(ContainSubstring(`ERROR invalid target platform "xxx"; supported platforms are: "cf", "neo", "xsa"`))
			Ω(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", "xyz1.mtar")).ShouldNot(BeAnExistingFile())
		})

		var _ = Describe("MBT build - generates Makefile and executes it", func() {

			It("MBT build for mta_demo", func() {
				dir, _ := os.Getwd()
				path := filepath.Join(dir, "testdata", "mta_demo")
				bin := filepath.FromSlash(binPath)
				_, errOut, err := execute(bin, "build -p=cf", path)
				Ω(err).Should(Succeed(), errOut)

				// Check the MTAR was generated
				validateMtaArchiveContents([]string{"node/", "node/data.zip", "node-js/", "node-js/data.zip"}, filepath.Join(path, "mta_archives", "mta_demo_0.0.1.mtar"))
			})

			It("MBT build - wrong platform", func() {
				dir, _ := os.Getwd()
				path := filepath.Join(dir, "testdata", "mta_demo")
				bin := filepath.FromSlash(binPath)
				_, _, err := execute(bin, "build -p=xxx", path)
				Ω(err).Should(HaveOccurred())
			})

			It("MBT build with timeout", func() {
				dir, _ := os.Getwd()
				path := filepath.Join(dir, "testdata", "moduletimeout")
				bin := filepath.FromSlash(binPath)

				start := time.Now()
				out, errOut, err := execute(bin, "build", path)
				elapsed := time.Since(start)
				Ω(err).Should(HaveOccurred())
				Ω(out + errOut).Should(ContainSubstring("timed out"))

				// Check elapsed time
				Ω(elapsed).Should(BeNumerically(">=", time.Duration(5)*time.Second))
				Ω(elapsed).Should(BeNumerically("<=", time.Duration(10)*time.Second))
			})
		})

		It("Generate MTAR - unsupported platform, module removed from mtad", func() {

			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("make")
			_, errOut, err := execute(bin, "-f Makefile.mta p=neo mtar=xyz", path)
			Ω(err).Should(Succeed(), errOut)
			mtarFilename := filepath.Join(dir, "testdata", "mta_demo", "mta_archives", "xyz.mtar")
			Ω(mtarFilename).Should(BeAnExistingFile())
			// check that module with unsupported platform 'neo' is not presented in mtad.yaml
			mtadContent, e := getFileContentFromZip(mtarFilename, "mtad.yaml")
			Ω(e).Should(Succeed())
			actual, e := mta.Unmarshal(mtadContent)
			Ω(e).Should(Succeed())
			expected, e := mta.Unmarshal([]byte(`
_schema-version: "3.1"
ID: mta_demo
version: 0.0.1
modules:
- name: node-js
  type: nodejs
  path: node-js
  provides:
  - name: node-js_api
    properties:
      url: ${default-url}
  parameters:
    name: nodejs
parameters:
  hcp-deployer-version: 1.1.0
`))
			Ω(e).Should(Succeed())
			Ω(actual).Should(Equal(expected))
		})

		It("Generate MTAR for mta_demo", func() {

			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("make")
			_, errOut, err := execute(bin, "-f Makefile.mta p=cf", path)
			Ω(err).Should(Succeed(), errOut)
			// Check the archive was generated
			mtarFilename := filepath.Join(dir, "testdata", "mta_demo", "mta_archives", demoArchiveName)
			Ω(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", demoArchiveName)).Should(BeAnExistingFile())
			// check that module with unsupported platform 'cf' is presented in mtad.yaml
			mtadContent, e := getFileContentFromZip(mtarFilename, "mtad.yaml")
			Ω(e).Should(Succeed())
			actual, e := mta.Unmarshal(mtadContent)
			Ω(e).Should(Succeed())
			expected, e := mta.Unmarshal([]byte(`
_schema-version: "3.1"
ID: mta_demo
version: 0.0.1
modules:
- name: node
  type: javascript.nodejs
  path: node
  provides:
  - name: node_api
    properties:
      url: ${default-url}
- name: node-js
  type: javascript.nodejs
  path: node-js
  provides:
  - name: node-js_api
    properties:
      url: ${default-url}
`))
			Ω(e).Should(Succeed())
			Ω(actual).Should(Equal(expected))
			validateMtaArchiveContents([]string{"node/", "node/data.zip", "node-js/", "node-js/data.zip"}, filepath.Join(path, "mta_archives", "mta_demo_0.0.1.mtar"))
		})
		//It("Generate MTAR for mta_java", func() {
		//
		//	dir, _ := os.Getwd()
		//	path := dir + filepath.FromSlash("/testdata/mta_java")
		//	bin := filepath.FromSlash("make")
		//	_, errOut, err, _ := execute(bin, "-f Makefile.mta p=cf", path)
		//				Ω(err).Should(Succeed(), errOut)
		//				// Check the archive was generated
		//				mtarFilename := filepath.Join(dir, "testdata", "mta_java", "mta_archives", javaArchiveName)
		//				Ω(filepath.Join(dir, "testdata", "mta_java", "mta_archives", javaArchiveName)).Should(BeAnExistingFile())
		//				// check that module with unsupported platform 'cf' is presented in mtad.yaml
		//				mtadContent, e := getFileContentFromZip(mtarFilename, "mtad.yaml")
		//				Ω(e).Should(Succeed())
		//				actual, e := mta.Unmarshal(mtadContent)
		//				Ω(e).Should(Succeed())
		//				expected, e := mta.Unmarshal([]byte(`
		//	_schema-version: 2.0.0
		//	ID: com.fetcher.project
		//	version: 0.0.1
		//	modules:
		//	- name: myModule
		//	 type: java.tomcat
		//	 path: myModule
		//	 requires:
		//	 - name: otracker-uaa
		//	 - name: otracker-managed-hdi
		//	 parameters:
		//	   buildpack: sap_java_buildpack
		//	   stack: cflinuxfs3
		//	resources:
		//	- name: otracker-uaa
		//	 type: com.sap.xs.uaa-space
		//	 parameters:
		//	   config-path: xs-security.json
		//	- name: otracker-managed-hdi
		//	 type: com.sap.xs.managed-hdi-container
		//	`))
		//				Ω(e).Should(Succeed())
		//				Ω(actual).Should(Equal(expected))
		//				validateMtaArchiveContents([]string{"myModule/", "myModule/java-xsahaa-1.1.2.war"}, filepath.Join(path, "mta_archives", "com.fetcher.project_0.0.1.mtar"))
		//})

		When("Running MBT commands with MTA extension descriptors", func() {
			var path string
			var mtarFilename string
			var makefileName string
			BeforeEach(func() {
				dir, err := os.Getwd()
				Ω(err).Should(Succeed())
				path = filepath.Join(dir, "testdata", "mta_demo")
				mtarFilename = filepath.Join(path, "mta_archives", demoArchiveName)
				makefileName = filepath.Join(path, "Makefile.mta")
			})
			AfterEach(func() {
				Ω(os.RemoveAll(makefileName)).Should(Succeed())
				Ω(os.RemoveAll(mtarFilename)).Should(Succeed())
			})

			var validateMtar = func() {
				// Check the MTAR was generated without the node-js module (since the extension file overrides its supported-platforms)
				Ω(mtarFilename).Should(BeAnExistingFile())
				validateMtaArchiveContents([]string{"node/", "node/data.zip"}, mtarFilename)

				// Check the mtad.yaml has the parts from the extension file
				// check that module with unsupported platform 'neo' is not present in the mtad.yaml
				mtadContent, e := getFileContentFromZip(mtarFilename, "mtad.yaml")
				Ω(e).Should(Succeed())
				actual, e := mta.Unmarshal(mtadContent)
				Ω(e).Should(Succeed())
				expected, e := mta.Unmarshal([]byte(`
_schema-version: "3.1"
ID: mta_demo
version: 0.0.1
modules:
- name: node
  type: javascript.nodejs
  path: node
  provides:
  - name: node_api
    properties:
      url: ${default-url}
`))
				Ω(e).Should(Succeed())
				Ω(actual).Should(Equal(expected))
			}

			It("MBT build for mta_demo with extension", func() {
				bin := filepath.FromSlash(binPath)
				_, errOut, err := execute(bin, "build -e=ext.mtaext -p=cf", path)
				Ω(err).Should(Succeed(), errOut)
				validateMtar()
			})

			It("MBT init and run make for mta_demo with extension - non-verbose", func() {
				bin := filepath.FromSlash(binPath)
				cmdOut, errOut, err := execute(bin, "init -e=ext.mtaext", path)
				Ω(err).Should(Succeed(), errOut)
				Ω(cmdOut).ShouldNot(BeNil())
				// Read the MakeFile was generated
				Ω(makefileName).Should(BeAnExistingFile())
				// generate mtar
				_, errOut, err = execute("make", "-f Makefile.mta p=cf", path)
				Ω(err).Should(Succeed(), errOut)
				validateMtar()
			})

			It("MBT init and run make for mta_demo with extension - verbose", func() {
				bin := filepath.FromSlash(binPath)
				cmdOut, errOut, err := execute(bin, "init -m=verbose -e=ext.mtaext", path)
				Ω(err).Should(Succeed(), errOut)
				Ω(cmdOut).ShouldNot(BeNil())
				// Read the MakeFile was generated
				Ω(makefileName).Should(BeAnExistingFile())
				// generate mtar
				_, errOut, err = execute("make", "-f Makefile.mta p=cf", path)
				Ω(err).Should(Succeed(), errOut)
				validateMtar()
			})
		})
	})

	var _ = Describe("Generate the Verbose Makefile and use it for mtar generation", func() {

		It("Generate Verbose Makefile", func() {
			dir, _ := os.Getwd()
			Ω(os.RemoveAll(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta"))).Should(Succeed())
			Ω(os.RemoveAll(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", demoArchiveName))).Should(Succeed())
			path := filepath.Join(dir, "testdata", "mta_demo")
			bin := filepath.FromSlash(binPath)
			cmdOut, errOut, err := execute(bin, "init -m=verbose", path)
			Ω(err).Should(Succeed(), errOut)
			Ω(cmdOut).ShouldNot(BeNil())
			// Read the MakeFile was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta")).Should(BeAnExistingFile())
			// generate mtar
			bin = filepath.FromSlash("make")
			_, errOut, err = execute(bin, "-f Makefile.mta p=cf", path)
			Ω(err).Should(Succeed(), errOut)
			// Check the archive was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", demoArchiveName)).Should(BeAnExistingFile())
		})

		Describe("module with dependencies", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "moduledep")
			archivePath := filepath.Join(path, "mta_archives", "f1_0.0.1.mtar")
			tempZipPath := filepath.Join(path, "mta_archives", "data.zip")

			AfterEach(func() {
				Ω(os.RemoveAll(filepath.Join(path, "Makefile.mta"))).Should(Succeed())
				Ω(os.RemoveAll(filepath.Join(path, "mta_archives"))).Should(Succeed())
				Ω(os.RemoveAll(filepath.Join(path, "public", "client"))).Should(Succeed())
				Ω(os.RemoveAll(filepath.Join(path, "public", "client2"))).Should(Succeed())
			})

			It("Generate Verbose Makefile with module dependencies", func() {
				bin := filepath.FromSlash(binPath)
				cmdOut, errOut, err := execute(bin, "init -m=verbose", path)
				Ω(err).Should(Succeed(), errOut)
				Ω(cmdOut).ShouldNot(BeNil())
				// Check the MakeFile was generated
				Ω(filepath.Join(path, "Makefile.mta")).Should(BeAnExistingFile())

				// Generate mtar
				bin = filepath.FromSlash("make")
				_, errOut, err = execute(bin, "-f Makefile.mta p=cf", path)
				Ω(err).Should(Succeed(), errOut)
				// Check the archive was generated
				Ω(archivePath).Should(BeAnExistingFile())
				validateMtaArchiveContents([]string{"module_with_dep/", "module_with_dep/data.zip"}, archivePath)

				// Extract data.zip and check its content
				err = extractFileFromZip(archivePath, "module_with_dep/data.zip", tempZipPath)
				Ω(err).Should(Succeed())
				validateArchiveContents([]string{"package.json", "client/", "client/client_package.json", "client2/", "client2/client_package.json"}, tempZipPath)
			})
		})
	})

	Describe("module with dependencies", func() {
		dir, _ := os.Getwd()
		path := filepath.Join(dir, "testdata", "moduledep")
		archivePath := filepath.Join(path, "mta_archives", "f1_0.0.1.mtar")
		tempZipPath := filepath.Join(path, "mta_archives", "data.zip")

		AfterEach(func() {
			Ω(os.RemoveAll(filepath.Join(path, "mta_archives"))).Should(Succeed())
			Ω(os.RemoveAll(filepath.Join(path, "public", "client"))).Should(Succeed())
			Ω(os.RemoveAll(filepath.Join(path, "public", "client2"))).Should(Succeed())
		})

		DescribeTable("Build MTA with module dependencies", func(additionalBuildOpts []string) {
			bin := filepath.FromSlash(binPath)
			cmdOut, _, err := executeWithArgs(bin, path, append([]string{"build"}, additionalBuildOpts...)...)
			Ω(err).Should(Succeed())
			Ω(cmdOut).ShouldNot(BeNil())

			// Check the archive was generated
			Ω(archivePath).Should(BeAnExistingFile())
			validateMtaArchiveContents([]string{"module_with_dep/", "module_with_dep/data.zip"}, archivePath)

			// Extract data.zip and check its content
			err = extractFileFromZip(archivePath, "module_with_dep/data.zip", tempZipPath)
			Ω(err).Should(Succeed())
			validateArchiveContents([]string{"package.json", "client/", "client/client_package.json", "client2/", "client2/client_package.json"}, tempZipPath)
		},
			Entry("Non-verbose build", []string{}),
			Entry("Parallel verbose build", []string{"--mode=verbose"}),
			Entry("Serial verbose build", []string{"--mode=verbose", "--jobs=1"}),
		)
	})

	var _ = Describe("MBT gen commands", func() {
		It("Generate mtad", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_demo")
			Ω(os.MkdirAll(filepath.Join(path, ".mta_demo_mta_build_tmp", "node"), os.ModePerm)).Should(Succeed())
			Ω(os.MkdirAll(filepath.Join(path, ".mta_demo_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			bin := filepath.FromSlash(binPath)
			_, errOut, err := execute(bin, "mtad gen", path)
			Ω(err).Should(Succeed(), errOut)
			mtadPath := filepath.Join(path, "mtad.yaml")
			Ω(mtadPath).Should(BeAnExistingFile())
			content, _ := ioutil.ReadFile(mtadPath)
			mtadObj, _ := mta.Unmarshal(content)
			Ω(len(mtadObj.Modules)).Should(Equal(2))
			Ω(mtadObj.Modules[0].Type).Should(Equal("javascript.nodejs"))
			Ω(mtadObj.Modules[1].Type).Should(Equal("javascript.nodejs"))
		})

		It("Generate mtad with mta extension", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_demo")
			Ω(os.MkdirAll(filepath.Join(path, ".mta_demo_mta_build_tmp", "node"), os.ModePerm)).Should(Succeed())
			Ω(os.MkdirAll(filepath.Join(path, ".mta_demo_mta_build_tmp", "node-js"), os.ModePerm)).Should(Succeed())
			bin := filepath.FromSlash(binPath)
			_, errOut, err := execute(bin, `mtad gen -e="ext.mtaext"`, path)
			Ω(err).Should(Succeed(), errOut)
			mtadPath := filepath.Join(path, "mtad.yaml")
			Ω(mtadPath).Should(BeAnExistingFile())
			content, _ := ioutil.ReadFile(mtadPath)
			mtadObj, _ := mta.Unmarshal(content)
			Ω(len(mtadObj.Modules)).Should(Equal(1))
			Ω(mtadObj.Modules[0].Name).Should(Equal("node"))
			Ω(mtadObj.Modules[0].Type).Should(Equal("javascript.nodejs"))
		})
	})

	var _ = Describe("Deploy basic mta archive", func() {
		AfterEach(func() {
			resourceCleanup("node")
			resourceCleanup("node-js")
		})
		It("Deploy MTAR", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo/mta_archives")
			bin := filepath.FromSlash("cf")
			// Execute deployment process with output to make the deployment success/failure more clear
			err := executeWithOutput(bin, "deploy "+demoArchiveName+" -f", path)
			Ω(err).Should(Succeed())
			// Check if the deploy succeeded by using curl command response.
			// Receiving the output status code 200 represents successful deployment
			args := "-s -o /dev/null -w '%{http_code}' " + os.Getenv("NODE_APP_ROUTE")
			path = dir + filepath.FromSlash("/testdata/mta_demo")
			bin = filepath.FromSlash("curl")
			cmdOut, errOut, err := executeEverySecond(bin, args, path)
			if len(errOut) > 0 {
				log.Println(errOut)
			}
			Ω(err).Should(Succeed())
			Ω(cmdOut).Should(Equal("'200'"))
		})
	})

	var _ = Describe("Generate merged mta.yaml", func() {
		AfterEach(func() {
			dir, _ := os.Getwd()
			Ω(os.RemoveAll(filepath.Join(dir, "testdata", "mta_demo", "result.yaml"))).Should(Succeed())
		})
		It("merges with one extension", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_demo")
			bin := filepath.FromSlash(binPath)
			_, errOut, err := execute(bin, `merge -e=ext.mtaext -n=result.yaml`, path)
			Ω(err).Should(Succeed(), errOut)
			mtadPath := filepath.Join(path, "result.yaml")
			Ω(mtadPath).Should(BeAnExistingFile())
			content, _ := ioutil.ReadFile(mtadPath)
			mtaObj, _ := mta.Unmarshal(content)
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

	var _ = Describe("Assemble MTAR", func() {
		var currentWorkingDirectory string
		var mtaAssemblePath string
		var resultMtarPath string

		BeforeEach(func() {
			currentWorkingDirectory, _ = os.Getwd()
			mtaAssemblePath = currentWorkingDirectory + filepath.FromSlash("/testdata/mta_assemble")
			resultMtarPath = filepath.Join(mtaAssemblePath, "mta_archives", "mta.assembly.example_1.3.3.mtar")
		})

		AfterEach(func() {
			Ω(os.RemoveAll(filepath.Join(mtaAssemblePath, "mta.assembly.example.mtar"))).Should(Succeed())
			Ω(os.Chdir(currentWorkingDirectory)).Should(Succeed())
			Ω(os.RemoveAll(filepath.FromSlash("./testdata/mta_assemble/mta_archives"))).Should(Succeed())
		})

		It("Assemble MTAR", func() {
			bin := filepath.FromSlash(binPath)
			cmdOut, errOut, err := execute(bin, "assemble", mtaAssemblePath)
			Ω(err).Should(Succeed(), errOut)
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(ContainSubstring("assembling the MTA project..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("copying the MTA content..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("generating the metadata..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("generating the MTA archive..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("the MTA archive generated at: " + resultMtarPath + "\n"))
			Ω(cmdOut).Should(ContainSubstring("cleaning temporary files..." + "\n"))
			Ω(resultMtarPath).Should(BeAnExistingFile())
			validateMtaArchiveContents([]string{
				"node.zip", "xs-security.json",
				"node/", "node/.eslintrc", "node/.eslintrc.ext", "node/.gitignore", "node/.npmrc", "node/jest.json", "node/package.json", "node/runTest.js", "node/server.js",
				"node/.che/", "node/.che/project.json",
				"node/tests/", "node/tests/sample-spec.js",
			}, resultMtarPath)
		})

		It("Assemble MTAR with MTA extension", func() {
			bin := filepath.FromSlash(binPath)
			cmdOut, errOut, err := execute(bin, fmt.Sprintf(`assemble -e="my.mtaext"`), mtaAssemblePath)
			Ω(err).Should(Succeed(), errOut)
			Ω(cmdOut).ShouldNot(Equal(""))
			Ω(resultMtarPath).Should(BeAnExistingFile())
			// TODO the assemble command copies the contents of excluded modules to the archive (unrelated to the extension files) even though
			// the modules are removed from the mtad.yaml and manifest.mf
			validateMtaArchiveContents([]string{
				"node.zip", "xs-security.json",
				"node/", "node/.eslintrc", "node/.eslintrc.ext", "node/.gitignore", "node/.npmrc", "node/jest.json", "node/package.json", "node/runTest.js", "node/server.js",
				"node/.che/", "node/.che/project.json",
				"node/tests/", "node/tests/sample-spec.js",
			}, resultMtarPath)
			mtadContent, e := getFileContentFromZip(resultMtarPath, "mtad.yaml")
			Ω(e).Should(Succeed())
			actual, e := mta.Unmarshal(mtadContent)
			Ω(e).Should(Succeed())
			expected, e := mta.Unmarshal([]byte(`
_schema-version: "3"
ID: mta.assembly.example
version: 1.3.3
modules:
  - name: example2
    type: javascript.nodejs
    path: node.zip
    provides:
      - name: backend
        properties:
          url: "${default-url}"
    requires:
      - name: assembly-uaa
resources:
  - name: mta-assembly-uaa
    type: org.cloudfoundry.managed-service
    parameters:
      service: xsuaa
      service-plan: space
      path: xs-security.json

`))
			Ω(e).Should(Succeed())
			Ω(actual).Should(Equal(expected))
		})
	})
})

func getFileContentFromZip(path string, filename string) ([]byte, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()
	for _, file := range zipFile.File {
		if strings.Contains(file.Name, filename) {
			fc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fc.Close() // If we got here there won't be another loop iteration
			return ioutil.ReadAll(fc)
		}
	}
	return nil, fmt.Errorf(`file "%s" not found`, filename)
}

func extractFileFromZip(archivePath string, filename string, dst string) error {
	zipFile, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
	}()
	var fileToExtract *zip.File = nil
	for _, file := range zipFile.File {
		if strings.Contains(file.Name, filename) {
			fileToExtract = file
		}
	}
	if fileToExtract == nil {
		return fmt.Errorf(`file "%s" not found`, filename)
	}

	in, err := fileToExtract.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()
	err = dir.WriteFile(in, dst)
	return err
}

func validateMtaArchiveContents(expectedAdditionalFilesInArchive []string, archiveLocation string) {
	expectedFilesInArchive := append(expectedAdditionalFilesInArchive, "META-INF/", "META-INF/MANIFEST.MF", "META-INF/mtad.yaml")
	validateArchiveContents(expectedFilesInArchive, archiveLocation)
}

func validateArchiveContents(expectedFilesInArchive []string, archiveLocation string) {
	archive, err := zip.OpenReader(archiveLocation)
	Ω(err).Should(Succeed())
	defer func() {
		_ = archive.Close()
	}()
	var filesInArchive []string
	for _, file := range archive.File {
		filesInArchive = append(filesInArchive, file.Name)
	}
	for _, expectedFile := range expectedFilesInArchive {
		Ω(contains(expectedFile, filesInArchive)).Should(BeTrue(), fmt.Sprintf("expected %s to be in the archive; archive contains %v", expectedFile, filesInArchive))
	}
	for _, existingFile := range filesInArchive {
		Ω(contains(existingFile, expectedFilesInArchive)).Should(BeTrue(), fmt.Sprintf("did not expect %s to be in the archive; archive contains %v", existingFile, filesInArchive))
	}
}

func contains(element string, elements []string) bool {
	for _, el := range elements {
		if el == element {
			return true
		}
	}
	return false
}

// execute with live output
func executeWithOutput(bin string, args string, path string) error {
	cmd := exec.Command(bin, strings.Split(args, " ")...)
	cmd.Dir = path
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "Error creating StdoutPipe for Cmd")
	}
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, "Error creating StderrPipe for Cmd")
	}

	scanner := bufio.NewScanner(cmdReader)
	scannerErr := bufio.NewScanner(cmdErrReader)

	go func() {
		for scanner.Scan() {
			fmt.Printf("process output | %s\n", scanner.Text())
		}
	}()
	go func() {
		for scannerErr.Scan() {
			fmt.Printf("process error output | %s\n", scannerErr.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "Error starting Cmd")
	}
	err = cmd.Wait()
	if err != nil {
		return errors.Wrapf(err, "Error waiting for Cmd")
	}
	return nil
}

// Delete deployed app
func resourceCleanup(appName string) {
	wd, _ := os.Getwd()
	path := wd + filepath.FromSlash("/testdata/mta_demo")
	bin := filepath.FromSlash("cf")
	cmdOut, cmdErr, err := execute(bin, "delete "+appName+" -r -f", path)
	Ω(err).Should(Succeed(), cmdErr)
	Ω(cmdOut).ShouldNot(BeEmpty())
}

// Execute command every second for 40 times
func executeEverySecond(bin string, args string, path string) (string, errorOut string, err error) {
	n := 0
	cmdOut, errOut, err := execute(bin, args, path)
	for range time.Tick(time.Second) {
		if n == 40 || strings.Compare(cmdOut, "'200'") == 0 {
			break
		}
		cmdOut, errOut, err = execute(bin, args, path)
		n++
	}
	return cmdOut, errOut, err
}

// Execute commands and get outputs
func execute(bin string, args string, path string) (output string, errorOutput string, err error) {
	return executeWithArgs(bin, path, strings.Split(args, " ")...)
}

func executeWithArgs(bin string, path string, args ...string) (output string, errorOutput string, err error) {
	// Provide list of commands
	cmd := exec.Command(bin, args...)
	// bin path
	cmd.Dir = path
	// std out
	stdoutBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	// std error
	stdErrBuf := &bytes.Buffer{}
	cmd.Stderr = stdErrBuf
	// Start command
	if err = cmd.Start(); err != nil {
		return "", "", err
	}
	// wait to the command to finish
	err = cmd.Wait()
	return stdoutBuf.String(), stdErrBuf.String(), err
}
