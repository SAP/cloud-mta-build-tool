package artifacts

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v2"

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
	})

	m := mta.Module{
		Name: "node-js",
		Path: "node-js",
	}

	var _ = Describe("ExecuteBuild", func() {

		It("Sanity", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "node-js", "cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())

		})

		It("Fails on platform validation", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "node-js", "xx", os.Getwd)).Should(HaveOccurred())

		})

		It("Fails on location initialization", func() {
			Ω(ExecuteBuild("", "", nil, "ui5app", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), nil, "ui5app", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("ExecutePack", func() {

		It("Sanity", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "node-js",
				"cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		It("Fails on platform validation", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), nil, "node-js",
				"xx", os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on location initialization", func() {
			Ω(ExecutePack("", "", nil, "ui5app", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
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

	var _ = Describe("Pack", func() {
		var _ = Describe("Sanity", func() {

			It("Default build-result - zip file, copy only", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &ep, &m, "node-js", "cf", "*.zip")).Should(Succeed())
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
				Ω(packModule(&ep, &ep, &mod, "node-js", "cf", "*.zip")).Should(HaveOccurred())
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
				Ω(packModule(&ep, &ep, &module, "htmlapp2", "cf", "")).Should(Succeed())
				Ω(getFullPathInTmpFolder("mta", "htmlapp2", "data.zip")).Should(BeAnExistingFile())
				validateArchiveContents([]string{"ignore"}, ep.GetTargetModuleZipPath("htmlapp2"), false)
			})

			It("Default build-result - zip file, copy only fails - no file matching wildcard", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &ep, &m, "node-js", "cf", "m*.zip")).Should(HaveOccurred())
			})

			// ep.GetTargetModuleDir(moduleName)
			It("Wrong source", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_unknown"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(HaveOccurred())
			})
			It("Target directory exists as a file", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: dir.Dev,
				}
				Ω(dir.CreateDirIfNotExist(filepath.Join(ep.GetTarget(), ".mta_with_zipped_module_mta_build_tmp"))).Should(Succeed())
				createFileInTmpFolder("mta_with_zipped_module", "node-js")
				Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(HaveOccurred())
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
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "res", "myresult.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"file1"}, resultLocation, true)
				})
				It("zips the module folder to build-artifact-name.zip when there is no build-result", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "myresult",
						},
					}
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "myresult.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/file1", "file2", "abc.war", "data.zip"}, resultLocation, true)
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
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "myresult.war")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"gulpfile.js", "server.js", "package.json"}, resultLocation, true)
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
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(HaveOccurred())
				})
				It("fails when build-artifact-name is not a string value", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": 1,
						},
					}
					err := packModule(&ep, &ep, &m, "node-js", "cf", "")
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
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getFullPathInTmpFolder("mta_with_subfolder", "node-js", "data.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/file1", "file2", "abc.war", "data.zip"}, resultLocation, true)
				})
				It("creates build-artifact-name.zip when build-artifact-name is same as a file that exists in the project", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "file2",
						},
					}
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getTestPath("result", ".mta_with_subfolder_mta_build_tmp", "node-js", "file2.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/file1", "file2", "abc.war", "data.zip"}, resultLocation, true)
				})
				It("creates build-artifact-name.zip when build-artifact-name is same as an archive file that exists in the project", func() {
					m := mta.Module{
						Name: "node-js",
						Path: "node-js",
						BuildParams: map[string]interface{}{
							"build-artifact-name": "abc",
						},
					}
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getTestPath("result", ".mta_with_subfolder_mta_build_tmp", "node-js", "abc.zip")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"res/file1", "file2", "abc.war", "data.zip"}, resultLocation, true)
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
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getTestPath("result", ".mta_with_subfolder_mta_build_tmp", "node-js", "abc.war")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"gulpfile.js", "server.js", "package.json"}, resultLocation, true)
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
					Ω(packModule(&ep, &ep, &m, "node-js", "cf", "")).Should(Succeed())
					resultLocation := getTestPath("result", ".mta_with_subfolder_mta_build_tmp", "node-js", "data.war")
					Ω(resultLocation).Should(BeAnExistingFile())
					validateArchiveContents([]string{"gulpfile.js", "server.js", "package.json"}, resultLocation, true)
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
			Ω(packModule(&ep, &ep, &mNoPlatforms, "node-js", "cf", "")).Should(Succeed())
			Ω(getTestPath("result", "mta_with_zipped_module_mta_build_tmp", "node-js", "data.zip")).
				ShouldNot(BeAnExistingFile())
		})

	})

	var _ = Describe("Build", func() {

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
				Ω(buildModule(&ep, &ep, &ep, "node-js", "cf")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
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
				Ω(buildModule(&ep, &ep, &ep, "node-js", "cf")).Should(HaveOccurred())
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
				err := buildModule(&ep, &ep, &ep, "node-js", "cf")
				checkError(err, commands.BadCommandMsg, `bash -c "sleep 1`)
			})

			It("Target folder exists as a file - dev", func() {
				createDirInTmpFolder("mta")
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				createFileInTmpFolder("mta", "node-js")
				Ω(buildModule(&ep, &ep, &ep, "node-js", "cf")).Should(HaveOccurred())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getResultPath(), MtaFilename: mtaFilename}
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
				Ω(buildModule(&ep, &ep, &ep, moduleName, "cf")).Should(HaveOccurred())
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
			},
				Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
				Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
				Entry("Invalid module name wrong build params", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
			)

			When("build parameters has timeout", func() {
				It("succeeds when timeout is not exceeded", func() {
					ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_with_timeout.yaml"}
					Ω(buildModule(&ep, &ep, &ep, "m2", "cf")).Should(Succeed())
					Ω(ep.GetTargetModuleZipPath("m2")).Should(BeAnExistingFile())
				})
				It("fails when timeout is exceeded", func() {
					ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_with_timeout.yaml"}
					err := buildModule(&ep, &ep, &ep, "m1", "cf")
					checkError(err, exec.ExecTimeoutMsg, "2s")
				})
				It("fails when timeout is not a string", func() {
					ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath(), MtaFilename: "mta_with_timeout.yaml"}
					err := buildModule(&ep, &ep, &ep, "m3", "cf")
					checkError(err, exec.ExecInvalidTimeoutMsg, "1")
				})
			})
		})
	})

	var _ = Describe("CopyMtaContent", func() {
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
			err := CopyMtaContent("", source, nil, false, func() (string, error) {
				return "", errors.New("error")
			})
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(ContainSubstring(copyContentFailedOnLocMsg))
		})
		It("With a deployment descriptor in the source directory with only modules paths as zip archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 0, map[string]string{}, map[string]string{"test-module-0": "zip", "test-module-1": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, true, os.Getwd)
			Ω(err).Should(Succeed())
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with one module path and one resource path as zip archive and a folder", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 1, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-module-0": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, true, os.Getwd)
			Ω(err).Should(Succeed())
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only resources with zip and module archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 0, 2, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-resource-1": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, true, os.Getwd)
			Ω(err).Should(Succeed())
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only resources with zip and module archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 2, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-resource-1": "zip", "test-module-0": "zip", "test-module-1": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			Ω(err).Should(Succeed())
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with only one module with zip and one requiredDependency with folder", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{"test-module-0": "test-required"}, map[string]string{"test-module-0": "folder", "test-required": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			Ω(err).Should(Succeed())
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only one module with zip and missing requiredDependency", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{"test-module-0": "test-required"}, map[string]string{"test-module-0": "folder", "test-required": "zip"})
			mta.Modules[0].Requires[0].Parameters["path"] = "zip1"
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, true, os.Getwd)
			Ω(err).Should(HaveOccurred())
		})

		It("With a deployment descriptor in the source directory with only one module with non-existing content", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{}, map[string]string{"test-module-0": "not-existing-contet"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			checkError(err, pathNotExistsMsg, "not-existing-content")
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(false))
			Ω(dirContainsAllElements(filepath.Join(source, info.Name()+dir.TempFolderSuffix), map[string]bool{}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with a module with non-existing content and another which has content", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 0, map[string]string{}, map[string]string{"test-module-0": "not-existing-contet", "test-module-1": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			checkError(err, pathNotExistsMsg, "not-existing-content")
			info, _ := os.Stat(source)
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
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, nil, false, os.Getwd)
			Ω(err).Should(Succeed())
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true}, true)).Should(Equal(true))
		})

		AfterEach(func() {
			os.RemoveAll(source)
		})
	})

	var _ = Describe("copyMtaContentFromPath", func() {
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

	var _ = Describe("cleanUpCopiedContent", func() {
		It("Sanity", func() {
			err := cleanUpCopiedContent(getTestPath(), []string{"result"})
			Ω(err).Should(Succeed())
		})
	})

	var _ = Describe("copyModuleArchiveToResultDir", func() {
		It("target folder is file", func() {
			err := copyModuleArchiveToResultDir(getTestPath("assembly", "data.jar"), getTestPath("assembly", "file", "file"), "m1")
			checkError(err, dir.FolderCreationFailedMsg, getTestPath("assembly", "file"))
		})
	})

})

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
	mta := mta.MTA{SchemaVersion: &[]string{"3.0.0"}[0], ID: "test-mta-id"}
	// populate modules
	for index := 0; index < numberOfModules; index++ {
		moduleName := "test-module-" + strconv.Itoa(index)
		mta.Modules = append(mta.Modules, generateTestModule(moduleName, moduleAndResourcesAndRequiredDependenciesContentTypes[moduleName], source))
	}

	for index := 0; index < numberOfResources; index++ {
		resourceName := "test-resource-" + strconv.Itoa(index)
		mta.Resources = append(mta.Resources, generateTestResource(resourceName, moduleAndResourcesAndRequiredDependenciesContentTypes[resourceName], source))
	}

	for moduleName, requiredDependencyName := range moduleWithReqDependencies {
		for _, module := range mta.Modules {
			if module.Name == moduleName {
				module.Requires = append(module.Requires, generateRequiredDependency(requiredDependencyName, moduleAndResourcesAndRequiredDependenciesContentTypes[requiredDependencyName], source))
			}
		}
	}
	return mta
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
		dir.CopyFile(getTestPath("mta_content_copy_test", "test.zip"), filepath.Join(source, "test.zip"))
		return "test.zip"
	}
	if contentType == "folder" {
		dir.CopyDir(getTestPath("mta_content_copy_test", "test-content"),
			filepath.Join(source, "test-content"), true, dir.CopyEntries)
		return "test-content"
	}

	return "not-existing-content"
}

func validateArchiveContents(expectedFilesInArchive []string, archiveLocation string, isExists bool) {
	archiveReader, err := zip.OpenReader(archiveLocation)
	Ω(err).Should(BeNil())
	defer archiveReader.Close()
	var filesInArchive []string
	for _, file := range archiveReader.File {
		filesInArchive = append(filesInArchive, file.Name)
	}
	for _, expectedFile := range expectedFilesInArchive {
		Ω(contains(expectedFile, filesInArchive)).Should(Equal(isExists), "Did not find "+expectedFile+" in archive")
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
