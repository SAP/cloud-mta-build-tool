package artifacts

import (
	"archive/zip"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("ModuleArch", func() {

	var config []byte

	BeforeEach(func() {
		config = make([]byte, len(commands.ModuleTypeConfig))
		copy(config, commands.ModuleTypeConfig)
		// Simplified commands configuration (performance purposes). removed "npm prune --production"
		commands.ModuleTypeConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands: 
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  commands:
`)
	})

	AfterEach(func() {
		commands.ModuleTypeConfig = make([]byte, len(config))
		copy(commands.ModuleTypeConfig, config)
		os.RemoveAll(getResultPath())
		os.RemoveAll(getFullPathInTmpFolder("mta"))
	})

	m := mta.Module{
		Name: "node-js",
		Path: "node-js",
	}

	Describe("ExecuteBuild", func() {

		It("Sanity", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "node-js", "cf", os.Getwd)).Should(Succeed())
			Ω(getFullPathInTmpFolder("mta", "node-js", "data.zip")).Should(BeAnExistingFile())

		})

		It("Fails on empty module", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "", "cf", os.Getwd)).Should(HaveOccurred())

		})

		It("Fails on platform validation", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "node-js", "xx", os.Getwd)).Should(HaveOccurred())

		})

		It("Fails on location initialization", func() {
			Ω(ExecuteBuild("", "", nil, "ui5app", "cf", failingGetWd)).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "ui5app", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	Describe("ExecuteSoloBuild", func() {

		Describe("Sanity, with target path, with all dependencies", func() {
			AfterEach(func() {
				Ω(os.Remove(getTestPath("mtaModelsBuild", "ui5app", "test2.txt"))).Should(Succeed())
				Ω(os.Remove(getTestPath("mtaModelsBuild", "ui5app", "test2_copy.txt"))).Should(Succeed())
				Ω(os.Remove(getTestPath("mtaModelsBuild", "ui5app2", "test2_copy.txt"))).Should(Succeed())
			})
			It("required module m2 has ready artifact 'test2.txt' and creates a new one 'test2_copy.txt'", func() {
				Ω(ExecuteSoloBuild(getTestPath("mtaModelsBuild"), getResultPath(), nil, []string{"m1", "m3"}, true, os.Getwd)).Should(Succeed())
				Ω(getTestPath("result", "data.zip")).Should(BeAnExistingFile())
				Ω(getTestPath("result", "m3.zip")).Should(BeAnExistingFile())
				validateArchiveContents([]string{"test.txt", "test2.txt", "test2_copy.txt"}, getTestPath("result", "data.zip"))
			})
			It("modules m1 and m2 have conflicting build results", func() {
				err := ExecuteSoloBuild(getTestPath("mtaModelsBuild"), getResultPath(), nil, []string{"m1", "m2"}, true, os.Getwd)
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(multiBuildWithPathsConflictMsg, "m2", "m1", getResultPath(), "data.zip")))
			})
		})

		Describe("Sanity, with target path, without dependencies", func() {
			AfterEach(func() {
				Ω(os.Remove(getTestPath("mtaModelsBuild", "ui5app", "test2.txt"))).Should(Succeed())
			})
			It("required module m2 has ready artifact 'test2.txt', only this one will be copied to m1", func() {
				Ω(ExecuteSoloBuild(getTestPath("mtaModelsBuild"), getResultPath(), nil, []string{"m1", "m3"}, false, os.Getwd)).Should(Succeed())
				Ω(getTestPath("result", "data.zip")).Should(BeAnExistingFile())
				Ω(getTestPath("result", "m3.zip")).Should(BeAnExistingFile())
				validateArchiveContents([]string{"test.txt", "test2.txt"}, getTestPath("result", "data.zip"))
			})
		})

		It("Sanity, no target path", func() {
			Ω(ExecuteSoloBuild(getTestPath("mta"), "", nil, []string{"node-js"}, true,
				func() (string, error) {
					return getTestPath("result", "test_dir"), nil
				})).Should(Succeed())
			Ω(getTestPath("result", "test_dir", ".mta_mta_build_tmp", "node-js", "data.zip")).Should(BeAnExistingFile())
		})

		It("fails on empty list of modules", func() {
			Ω(ExecuteSoloBuild(getTestPath("mta"), getResultPath(), nil, []string{}, true, os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on source getter", func() {
			err := ExecuteSoloBuild("", "", nil, []string{"ui5app"}, true, failingGetWd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(buildFailedMsg, "ui5app")))
		})

		It("Fails on source getter with multiple modules", func() {
			err := ExecuteSoloBuild("", "", nil, []string{"ui5app", "ui5app2"}, true, failingGetWd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(multiBuildFailedMsg))
		})

		It("Fails on wrong build dependencies - on sortModules", func() {
			Ω(ExecuteSoloBuild(getTestPath("mtahtml5"), "", []string{"mtaExtWithCyclicDependencies.yaml"},
				[]string{"ui5app"}, true, os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on unknown builder", func() {
			Ω(ExecuteSoloBuild(getTestPath("mtahtml5"), "", []string{"mtaExtWithUnkownBuilder.yaml"},
				[]string{"ui5app"}, true, os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on location initialization", func() {
			counter := 0
			Ω(ExecuteSoloBuild("", "", nil, []string{"ui5app"}, true, func() (string, error) {
				if counter == 0 {
					counter++
					return "", nil
				}
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteSoloBuild(getTestPath("mta"), getResultPath(), nil, []string{"ui5app"}, true, os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on getting default source", func() {
			Ω(ExecuteSoloBuild(getTestPath("mta"), "", nil, []string{"ui5app"}, true,
				failingGetWd)).Should(HaveOccurred())
		})

		// the only purpose of the test to support coverage
		// it covers the path of the code that will never happen -
		// failure on dir.Location after 2 successful calls to getSoloModuleBuildAbsSource & getSoloModuleBuildAbsTarget
		It("Fails on creation of Location object", func() {
			counter := 1
			Ω(ExecuteSoloBuild("", "", nil, []string{"ui5app"}, true, func() (string, error) {
				if counter <= 2 {
					counter++
					return "", nil
				}
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("getSoloModuleBuildAbsTarget fails on current folder getter", func() {
			_, err := getSoloModuleBuildAbsTarget(getTestPath(), "", "m1", failingGetWd)
			Ω(err).Should(HaveOccurred())
		})
	})

	Describe("buildModules", func() {

		AfterEach(func() {
			Ω(os.RemoveAll(getTestPath("result"))).Should(Succeed())
		})

		It("sanity", func() {
			Ω(buildModules(getTestPath("mtahtml5"), getTestPath("result"), nil, []string{"ui5app"}, map[string]bool{"ui5app": true}, os.Getwd)).Should(Succeed())
			Ω(getTestPath("result", "data.zip")).Should(BeAnExistingFile())
		})

		It("fails on module location getter", func() {
			Ω(buildModules(getTestPath("mtahtml5"), "", nil, []string{"ui5app2"}, map[string]bool{}, failingGetWd)).Should(HaveOccurred())
		})

		It("fails on wrong selected module", func() {
			Ω(buildModules(getTestPath("mtahtml5"), "", nil, []string{"unknown"}, map[string]bool{"unknown": true}, os.Getwd)).Should(HaveOccurred())
		})

		It("fails on module location getter of dependency", func() {
			Ω(buildModules("", "", nil, []string{"ui5app"}, map[string]bool{"ui5app": true}, failingGetWd)).Should(HaveOccurred())
		})

		It("fails on buildModule because of the unknown builder", func() {
			Ω(buildModules(getTestPath("mtahtml5"), getTestPath("result"), nil, []string{"ui5app3"}, map[string]bool{"ui5app3": true}, os.Getwd)).Should(HaveOccurred())
		})
	})

	Describe("ExecutePack", func() {

		It("Sanity", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "node-js",
				"cf", os.Getwd)).Should(Succeed())
			Ω(getFullPathInTmpFolder("mta", "node-js", "data.zip")).Should(BeAnExistingFile())
		})

		It("no-source module", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "no_source",
				"cf", os.Getwd)).Should(Succeed())
			Ω(getFullPathInTmpFolder("mta", "node-js", "data.zip")).ShouldNot(BeAnExistingFile())
		})

		It("Fails on empty path", func() {
			Ω(ExecutePack(getTestPath("mta_no_path"), getResultPath(), nil, "no_path",
				"cf", os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on platform validation", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "node-js",
				"xx", os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on location initialization", func() {
			Ω(ExecutePack("", "", nil, "ui5app", "cf", failingGetWd)).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "ui5appx",
				"cf", os.Getwd)).Should(HaveOccurred())
		})

		It("Target folder exists as file", func() {
			createDirInTmpFolder("mta")
			createFileInTmpFolder("mta", "node-js")
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "node-js",
				"cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	Describe("Pack", func() {
		var _ = Describe("Sanity", func() {

			It("Default build-result - zip file, copy only", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &m, "node-js", "cf", "*.zip", true, map[string]string{})).Should(Succeed())
				Ω(getFullPathInTmpFolder("mta_with_zipped_module", "node-js", "abc.zip")).Should(BeAnExistingFile())
			})
			It("Build results - zip file not exists, fails", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				mod := mta.Module{
					Name: "node-js",
					Path: "notExists",
				}
				Ω(packModule(&ep, &mod, "node-js", "cf", "*.zip", true, map[string]string{})).Should(HaveOccurred())
			})

			It("zip file with ignored folder", func() {
				module := mta.Module{
					Name: "htmlapp2",
					Path: "htmlapp2",
					BuildParams: map[string]interface{}{
						// TODO this test doesn't check the ignore correctly. Even if there is no ignore it will pass.
						"ignore": []interface{}{"ignore/"},
					},
				}
				ep := dir.Loc{
					SourcePath: getTestPath("mta"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &module, "htmlapp2", "cf", "", true, map[string]string{})).Should(Succeed())
				Ω(getFullPathInTmpFolder("mta", "htmlapp2", "data.zip")).Should(BeAnExistingFile())
				validateArchiveContentsExcludes([]string{"ignore"}, getFullPathInTmpFolder("mta", "htmlapp2", "data.zip"))
			})

			It("Default build-result - zip file, copy only fails - no file matching wildcard", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &m, "node-js", "cf", "m*.zip", true, map[string]string{})).Should(HaveOccurred())
			})

			// ep.GetTargetModuleDir(moduleName)
			It("Wrong source", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_unknown"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(HaveOccurred())
			})
			It("Target directory exists as a file", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(dir.CreateDirIfNotExist(getFullPathInTmpFolder("mta_with_zipped_module"))).Should(Succeed())
				createFileInTmpFolder("mta_with_zipped_module", "node-js")
				Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(HaveOccurred())
			})
			When("build-artifact-name is defined for the module", func() {
				var ep dir.Loc
				BeforeEach(func() {
					ep = dir.Loc{
						SourcePath: getTestPath("mta_with_subfolder"),
						TargetPath: getResultPath(),
						Descriptor: dir.Dev,
					}
				})

				It("zips the build-result folder to build-artifact-name.zip when the build-result is defined and points to a folder", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-result":        "res",
							"build-artifact-name": "myresult",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "res", "myresult.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"file1"}, resultLocation)
				})
				It("zips the module folder to build-artifact-name.zip when there is no build-result", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "myresult",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "myresult.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/", "res/file1", "file2", "abc.war", "data.zip"}, resultLocation)
				})
				It("copies the build-result file to build-artifact-name when build-result is an archive file", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-result":        "abc.war",
							"build-artifact-name": "myresult",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "myresult.war")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"gulpfile.js", "server.js", "package.json"}, resultLocation)
				})
				It("fails when build-result doesn't exist", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-result":        "abc2.zip",
							"build-artifact-name": "myresult",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(HaveOccurred())
				})
				It("fails when build-artifact-name is not a string value", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": 1,
						},
					}
					err := packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})
					Ω(err).Should(HaveOccurred())
					Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(buildops.WrongBuildArtifactNameMsg, "1", "node-js")))
				})
				It("creates data.zip when build-artifact-name is data", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "data",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "data.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/", "res/file1", "file2", "abc.war", "data.zip"}, resultLocation)
				})
				It("creates build-artifact-name.zip when build-artifact-name is same as a file that exists in the project", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "file2",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "file2.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/", "res/file1", "file2", "abc.war", "data.zip"}, resultLocation)
				})
				It("creates build-artifact-name.zip when build-artifact-name is same as an archive file that exists in the project", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "abc",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "abc.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/", "res/file1", "file2", "abc.war", "data.zip"}, resultLocation)
				})
				It("creates build-artifact-name with the build-result extension when build-artifact-name is same as build-result, which is an archive file", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-result":        "abc.war",
							"build-artifact-name": "abc",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "abc.war")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"gulpfile.js", "server.js", "package.json"}, resultLocation)
				})
				It("creates build-artifact-name with the build-result extension when build-artifact-name is same as an archive file and different from build-result", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-result":        "abc.war",
							"build-artifact-name": "data",
						},
					}
					Ω(packModule(&ep, &m, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "data.war")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"gulpfile.js", "server.js", "package.json"}, resultLocation)
				})
			})
		})

		It("No platforms - no pack", func() {
			ep := dir.Loc{
				SourcePath: getTestPath("mta_with_zipped_module"),
				TargetPath: getResultPath(),
				Descriptor: dir.Dep,
			}
			mNoPlatforms := mta.Module{
				Name: "node-js",
				Path: "node-js",
				BuildParams: map[string]interface{}{
					buildops.SupportedPlatformsParam: []string{},
				},
			}
			Ω(packModule(&ep, &mNoPlatforms, "node-js", "cf", "", true, map[string]string{})).Should(Succeed())
			Ω(getFullPathInTmpFolder("mta_with_zipped_module", "node-js", "data.zip")).
				ShouldNot(BeAnExistingFile())
		})

	})

	Describe("Build", func() {

		var _ = Describe("build Module", func() {

			var config []byte

			BeforeEach(func() {
				config = make([]byte, len(commands.ModuleTypeConfig))
				copy(config, commands.ModuleTypeConfig)
				// Simplified commands configuration (performance purposes). removed "npm prune --production"
				commands.ModuleTypeConfig = []byte(`
builders:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands:
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  commands:
`)
			})

			It("Sanity", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				Ω(buildModule(&ep, &ep, "node-js", "cf", true, true, map[string]string{})).Should(Succeed())
				Ω(getFullPathInTmpFolder("mta", "node-js", "data.zip")).Should(BeAnExistingFile())
			})

			It("Sanity, not packed - platform not supported", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				Ω(buildModule(&ep, &ep, "node-js", "neo", true, true, map[string]string{})).Should(Succeed())
				Ω(getFullPathInTmpFolder("mta", "node-js", "data.zip")).ShouldNot(BeAnExistingFile())
			})

			It("Sanity, packed - platform not checked", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				Ω(buildModule(&ep, &ep, "node-js", "neo", false, true, map[string]string{})).Should(Succeed())
				Ω(getFullPathInTmpFolder("mta", "node-js", "data.zip")).Should(BeAnExistingFile())
			})

			It("empty path", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta_no_path"), TargetPath: getResultPath()}
				Ω(buildModule(&ep, &ep, "no_path", "cf", true, true, map[string]string{})).Should(HaveOccurred())
			})

			It("no source module", func() {
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				Ω(buildModule(&ep, &ep, "no_source", "cf", true, true, map[string]string{})).Should(Succeed())
				Ω(getTestPath("mta", "node-js", "data.zip")).ShouldNot(BeAnExistingFile())
			})

			It("Commands fail", func() {
				commands.ModuleTypeConfig = []byte(`
module-types:
- name: html5
  info: "installing module dependencies & execute grunt & remove dev dependencies"
  path: "path to config file which override the following default commands"
  commands:
    - command: go test exec_unknownTest.go
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  commands:
    - command: go test exec_unknownTest.go
`)

				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				Ω(buildModule(&ep, &ep, "node-js", "cf", true, true, map[string]string{})).Should(HaveOccurred())
			})

			It("fails when the command is invalid", func() {
				commands.ModuleTypeConfig = []byte(`
module-types:
- name: nodejs
  info: "build nodejs application"
  path: "path to config file which override the following default commands"
  commands:
    - command: bash -c "sleep 1
`)

				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				err := buildModule(&ep, &ep, "node-js", "cf", true, true, map[string]string{})
				checkError(err, commands.BadCommandMsg, `bash -c "sleep 1`)
			})

			It("Target folder exists as a file - dev", func() {
				createDirInTmpFolder("mta")
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				createFileInTmpFolder("mta", "node-js")
				Ω(buildModule(&ep, &ep, "node-js", "cf", true, true, map[string]string{})).Should(HaveOccurred())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getResultPath(), MtaFilename: mtaFilename}
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
				Ω(buildModule(&ep, &ep, moduleName, "cf", true, true, map[string]string{})).Should(HaveOccurred())
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
			},
				Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
				Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
				Entry("Invalid module name wrong build params", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
			)

			When("build parameters has timeout", func() {
				It("succeeds when timeout is not exceeded", func() {
					ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_with_timeout.yaml"}
					Ω(buildModule(&ep, &ep, "m2", "cf", true, true, map[string]string{})).Should(Succeed())
					Ω(getFullPathInTmpFolder("mta", "m2", "data.zip")).Should(BeAnExistingFile())
				})
				It("fails when timeout is exceeded", func() {
					ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_with_timeout.yaml"}
					err := buildModule(&ep, &ep, "m1", "cf", true, true, map[string]string{})
					checkError(err, exec.ExecTimeoutMsg, "2s")
				})
				It("fails when timeout is not a string", func() {
					ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_with_timeout.yaml"}
					err := buildModule(&ep, &ep, "m3", "cf", true, true, map[string]string{})
					checkError(err, exec.ExecInvalidTimeoutMsg, "1")
				})
			})
		})
	})

	Describe("CopyMtaContent", func() {
		var source string
		defaultDeploymentDescriptorName := "mtad.yaml"
		BeforeEach(func() {
			source, _ = ioutil.TempDir("", "testing-mta-content")
		})
		It("Without no deployment descriptor in the source directory", func() {
			err := CopyMtaContent(source, source, nil, true, os.Getwd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(dir.ReadFailedMsg, filepath.Join(source, "mtad.yaml"))))
		})
		It("Location initialization fails", func() {
			err := CopyMtaContent("", source, nil, false, failingGetWd)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(copyContentFailedOnLocMsg))
		})
		It("With a deployment descriptor in the source directory with only modules paths as zip archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 0, map[string]string{}, map[string]string{"test-module-0": "zip", "test-module-1": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, true, os.Getwd)).Should(Succeed())
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with one module path and one resource path as zip archive and a folder", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 1, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-module-0": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, true, os.Getwd)).Should(Succeed())
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only resources with zip and module archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 0, 2, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-resource-1": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, true, os.Getwd)).Should(Succeed())
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only resources with zip and module archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 2, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-resource-1": "zip", "test-module-0": "zip", "test-module-1": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, false, os.Getwd)).Should(Succeed())
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with only one module with zip and one requiredDependency with folder", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{"test-module-0": "test-required"}, map[string]string{"test-module-0": "folder", "test-required": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, false, os.Getwd)).Should(Succeed())
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only one module with zip and missing requiredDependency", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{"test-module-0": "test-required"}, map[string]string{"test-module-0": "folder", "test-required": "zip"})
			mta.Modules[0].Requires[0].Parameters["path"] = "zip1"
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, true, os.Getwd)).Should(HaveOccurred())
		})

		It("With a deployment descriptor in the source directory with only one module with non-existing content", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{}, map[string]string{"test-module-0": "not-existing-contet"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			checkError(err, pathNotExistsMsg, "not-existing-content")
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(false))
			Ω(dirContainsAllElements(filepath.Join(source, info.Name()+dir.TempFolderSuffix), map[string]bool{}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with a module with non-existing content and another which has content", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 0, map[string]string{}, map[string]string{"test-module-0": "not-existing-contet", "test-module-1": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			checkError(err, pathNotExistsMsg, "not-existing-content")
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(false))
			Ω(dirContainsAllElements(filepath.Join(source, info.Name()+dir.TempFolderSuffix), map[string]bool{}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with a lot of modules with zip contentt", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			modulesWithSameContent := make(map[string]string)
			for index := 0; index < 10; index++ {
				modulesWithSameContent["test-module-"+strconv.Itoa(index)] = "zip"
			}
			mta := generateTestMta(source, 10, 0, map[string]string{}, modulesWithSameContent)
			mtaBytes, _ := yaml.Marshal(mta)
			Ω(ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)).Should(Succeed())
			Ω(CopyMtaContent(source, source, nil, false, os.Getwd)).Should(Succeed())
			info, err := os.Stat(source)
			Ω(err).Should(Succeed())
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true}, true)).Should(Equal(true))
		})

		AfterEach(func() {
			os.RemoveAll(source)
		})
	})

	Describe("copyMtaContentFromPath", func() {
		It("content is file; fails because target folder exists and it's not a folder, but a file", func() {
			file, _ := os.Create(getTestPath("result"))
			defer file.Close()
			Ω(copyMtaContentFromPath(getTestPath("mta", "mta.yaml"), getTestPath("result", "mta.yaml"),
				getTestPath("result", "mta.yaml"), getTestPath("result"), true)).Should(HaveOccurred())
		})
		It("content is file; fails because target folder exists and it's not a folder, but a file; not parallel", func() {
			file, _ := os.Create(getTestPath("result"))
			defer file.Close()
			Ω(copyMtaContentFromPath(getTestPath("mta", "mta.yaml"), getTestPath("result", "mta.yaml"),
				getTestPath("result", "mta.yaml"), getTestPath("result"), false)).Should(HaveOccurred())
		})
	})

	Describe("cleanUpCopiedContent", func() {
		It("Sanity", func() {
			err := cleanUpCopiedContent(getTestPath(), []string{"result"})
			Ω(err).Should(Succeed())
		})
	})

	Describe("copyModuleArchiveToResultDir", func() {
		It("target folder is file", func() {
			err := copyModuleArchiveToResultDir(getTestPath("assembly", "data.jar"), getTestPath("assembly", "file", "file"), "m1")
			checkError(err, dir.FolderCreationFailedMsg, getTestPath("assembly", "file"))
		})

		It("target is existing folder", func() {
			err := copyModuleArchiveToResultDir(getTestPath("assembly", "file"), getTestPath("assembly", "folder3"), "m1")
			checkError(err, packFailedOnCopyMsg, "m1", getTestPath("assembly", "file"), getTestPath("assembly", "folder3"))
		})
	})

	Describe("getModuleLocation", func() {
		It("sanity, target provided", func() {
			loc, err := getModuleLocation(getTestPath(), getTestPath("result"), "m1", nil, os.Getwd)
			Ω(err).Should(Succeed())
			Ω(loc.GetTargetModuleDir("m1")).Should(Equal(getTestPath("result")))
		})

		It("sanity, target not provided - default one is used", func() {
			loc, err := getModuleLocation(getTestPath("mta"), "", "m1", nil, os.Getwd)
			Ω(err).Should(Succeed())
			currentFolder, err := os.Getwd()
			Ω(err).Should(Succeed())
			Ω(loc.GetTargetModuleDir("m1")).Should(Equal(filepath.Join(currentFolder, ".mta_mta_build_tmp", "m1")))
		})

		It("fails on default target detecting", func() {
			_, err := getModuleLocation("", "", "m1", nil, failingGetWd)
			Ω(err).Should(HaveOccurred())
		})

		// this test was added for coverage purpose only,
		// in case of problems with getWd function failure will happen before location initialization
		It("fails on location initialization", func() {
			count := 0
			_, err := getModuleLocation("", "", "m1", nil, func() (string, error) {
				if count == 0 {
					count++
					return "", nil
				}
				return "", errors.New("error")
			})

			Ω(err).Should(HaveOccurred())
		})
	})

	Describe("collectSelectedModulesAndDependencies", func() {

		var mtaObj *mta.MTA

		BeforeEach(func() {
			mtaObj = getMtaObj("mtahtml5", "mtaWithWrongBuildRequirements.yaml")
		})

		It("sanity", func() {
			collection := make(map[string]bool)
			err := collectSelectedModulesAndDependencies(mtaObj, collection, "m1")
			Ω(err).Should(Succeed())
			validateMapKeys(collection, []string{"m1", "m4", "m3", "m2"})
		})
		DescribeTable("failures", func(targetModule string) {
			collection := make(map[string]bool)
			err := collectSelectedModulesAndDependencies(mtaObj, collection, targetModule)
			Ω(err).Should(HaveOccurred())
		},
			Entry("fails on none existing target module", "m5"),
			Entry("fails on none existing module required by the target module", "n1"),
			Entry("fails on none existing module required by the module that is required by the target module", "n2"))
	})

	Describe("sortModules", func() {
		It("sanity", func() {
			mtaObj := getMtaObj("mtahtml5", "mtaWithBuildRequirements.yaml")
			allModulesSorted, err := buildops.GetModulesNames(mtaObj)
			Ω(err).Should(Succeed())
			selectedModulesMap := map[string]bool{"n1": true, "m1": true}
			selectedModulesSorted := sortModules(allModulesSorted, selectedModulesMap)
			Ω(selectedModulesSorted).Should(Equal([]string{"m1", "n1"}))
		})
	})
})

func validateMapKeys(actualMap map[string]bool, expectedKeys []string) {
	var actualKeys []string
	for actualKey := range actualMap {
		actualKeys = append(actualKeys, actualKey)
	}
	Ω(expectedKeys).Should(ConsistOf(actualKeys))
}

func getMtaObj(projectName string, mtaFilename string) *mta.MTA {
	loc, err := dir.Location(getTestPath(projectName), getTestPath("result"), dir.Dev, nil, os.Getwd)
	Ω(err).Should(Succeed())
	loc.MtaFilename = mtaFilename
	mtaObj, err := loc.ParseFile()
	Ω(err).Should(Succeed())
	return mtaObj
}

func dirContainsAllElements(source string, elements map[string]bool, validateEntitiesCount bool) bool {
	sourceElements, _ := ioutil.ReadDir(source)
	if validateEntitiesCount {
		Ω(len(sourceElements)).Should(Equal(len(elements)))
	}
	for _, el := range sourceElements {
		if elements[el.Name()] {
			delete(elements, el.Name())
		}
	}

	return len(elements) == 0
}

func generateTestMta(source string, numberOfModules, numberOfResources int, moduleWithReqDependencies, moduleAndResourcesAndRequiredDependenciesContentTypes map[string]string) mta.MTA {
	mtaObj := mta.MTA{SchemaVersion: &[]string{"3.0.0"}[0], ID: "test-mta-id"}
	// populate modules
	for index := 0; index < numberOfModules; index++ {
		moduleName := "test-module-" + strconv.Itoa(index)
		mtaObj.Modules = append(mtaObj.Modules, generateTestModule(moduleName, moduleAndResourcesAndRequiredDependenciesContentTypes[moduleName], source))
	}

	for index := 0; index < numberOfResources; index++ {
		resourceName := "test-resource-" + strconv.Itoa(index)
		mtaObj.Resources = append(mtaObj.Resources, generateTestResource(resourceName, moduleAndResourcesAndRequiredDependenciesContentTypes[resourceName], source))
	}

	for moduleName, requiredDependencyName := range moduleWithReqDependencies {
		for _, module := range mtaObj.Modules {
			if module.Name == moduleName {
				module.Requires = append(module.Requires, generateRequiredDependency(requiredDependencyName, moduleAndResourcesAndRequiredDependenciesContentTypes[requiredDependencyName], source))
			}
		}
	}
	return mtaObj
}

func generateRequiredDependency(name, contentType, source string) mta.Requires {
	requiredDep := mta.Requires{Name: name}
	requiredDep.Parameters = make(map[string]interface{})
	requiredDep.Parameters["path"] = getContentPath(contentType, source)
	return requiredDep
}

func generateTestResource(resourceName, contentType, source string) *mta.Resource {
	resource := mta.Resource{Name: resourceName, Type: "test-resource-type"}
	resource.Parameters = make(map[string]interface{})
	resource.Parameters["path"] = getContentPath(contentType, source)
	return &resource
}

func generateTestModule(moduleName, contentType, source string) *mta.Module {
	module := mta.Module{Name: moduleName, Type: "test-module-type"}
	module.Path = getContentPath(contentType, source)
	return &module
}

func getContentPath(contentType, source string) string {
	if contentType == "zip" {
		Ω(dir.CopyFile(getTestPath("mta_content_copy_test", "test.zip"), filepath.Join(source, "test.zip"))).Should(Succeed())
		return "test.zip"
	}
	if contentType == "folder" {
		Ω(dir.CopyDir(
			getTestPath("mta_content_copy_test", "test-content"), filepath.Join(source, "test-content"),
			true, dir.CopyEntries,
		)).Should(Succeed())
		return "test-content"
	}

	return "not-existing-content"
}

func validateArchiveContents(expectedFilesInArchive []string, archiveLocation string) {
	archiveReader, err := zip.OpenReader(archiveLocation)
	Ω(err).Should(Succeed())
	defer archiveReader.Close()
	var filesInArchive []string
	for _, file := range archiveReader.File {
		filesInArchive = append(filesInArchive, file.Name)
	}
	for _, expectedFile := range expectedFilesInArchive {
		Ω(contains(expectedFile, filesInArchive)).Should(BeTrue(), fmt.Sprintf("expected %s to be in the archive; archive contains %v", expectedFile, filesInArchive))
	}
	for _, existingFile := range filesInArchive {
		Ω(contains(existingFile, expectedFilesInArchive)).Should(BeTrue(), fmt.Sprintf("did not expect %s to be in the archive; archive contains %v", existingFile, filesInArchive))
	}
}

func validateArchiveContentsExcludes(unexpectedFilesInArchive []string, archiveLocation string) {
	archiveReader, err := zip.OpenReader(archiveLocation)
	Ω(err).Should(Succeed())
	defer archiveReader.Close()
	var filesInArchive []string
	for _, file := range archiveReader.File {
		filesInArchive = append(filesInArchive, file.Name)
	}
	for _, unexpectedFile := range unexpectedFilesInArchive {
		Ω(contains(unexpectedFile, filesInArchive)).Should(BeFalse(), fmt.Sprintf("did not expect %s to be in the archive; archive contains %v", unexpectedFile, filesInArchive))
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

func failingGetWd() (string, error) {
	return "", errors.New("error")
}
