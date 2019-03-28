package artifacts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SAP/cloud-mta/mta"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
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
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), "dev", "node-js", "cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())

		})

		It("Fails on platform validation", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), "dev", "node-js", "xx", os.Getwd)).Should(HaveOccurred())

		})

		It("Fails on location initialization", func() {
			Ω(ExecuteBuild("", "", "dev", "ui5app", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteBuild(getTestPath("mta"), getResultPath(), "dev", "ui5app", "cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("ExecutePack", func() {

		It("Sanity", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), "dev", "node-js",
				"cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		It("Fails on platform validation", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), "dev", "node-js",
				"xx", os.Getwd)).Should(HaveOccurred())
		})

		It("Fails on location initialization", func() {
			Ω(ExecutePack("", "", "dev", "ui5app", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), "dev", "ui5appx",
				"cf", os.Getwd)).Should(HaveOccurred())
		})

		It("Target folder exists as file", func() {
			os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp"), os.ModePerm)
			createFile("result", ".mta_mta_build_tmp", "node-js")
			Ω(ExecutePack(getTestPath("mta"), getResultPath(), "dev", "node-js",
				"cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("ExecuteZip", func() {

		It("Sanity", func() {
			Ω(ExecuteZip(getTestPath("mta_zip_builder"), getResultPath(), "dev", "node-js",
				"cf", os.Getwd)).Should(Succeed())
			loc := dir.Loc{SourcePath: getTestPath("mta_zip_builder"), TargetPath: getResultPath()}
			Ω(loc.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})

		It("Fails on location initialization", func() {
			Ω(ExecuteZip("", "", "dev", "mta_zip_builder", "cf", func() (string, error) {
				return "", errors.New("err")
			})).Should(HaveOccurred())
		})

		It("Fails on wrong module", func() {
			Ω(ExecuteZip(getTestPath("mta_zip_builder"), getResultPath(), "dev", "ui5appx",
				"cf", os.Getwd)).Should(HaveOccurred())
		})

		It("Target folder exists as file", func() {
			os.MkdirAll(getTestPath("result", ".mta_zip_builder_mta_build_tmp"), os.ModePerm)
			createFile("result", ".mta_zip_builder_mta_build_tmp", "node-js")
			Ω(ExecuteZip(getTestPath("mta_zip_builder"), getResultPath(), "dev", "node-js",
				"cf", os.Getwd)).Should(HaveOccurred())
		})
	})

	var _ = Describe("Pack", func() {
		var _ = Describe("Sanity", func() {

			It("Deployment descriptor - Copy only", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: "dep",
				}
				Ω(packModule(&ep, true, &m, "node-js", "cf")).Should(Succeed())
				Ω(getTestPath("result", ".mta_with_zipped_module_mta_build_tmp", "node-js", "data.zip")).Should(BeAnExistingFile())
			})

			// ep.GetTargetModuleDir(moduleName)
			It("Wrong source", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_unknown"),
					TargetPath: getResultPath(),
					Descriptor: "dev",
				}
				Ω(packModule(&ep, false, &m, "node-js", "cf")).Should(HaveOccurred())
			})
			It("Target directory exists as a file", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: "dev",
				}
				os.MkdirAll(filepath.Join(ep.GetTarget(), ".mta_with_zipped_module_mta_build_tmp"), os.ModePerm)
				createFile("result", ".mta_with_zipped_module_mta_build_tmp", "node-js")
				Ω(packModule(&ep, false, &m, "node-js", "cf")).Should(HaveOccurred())
			})
		})

		It("No platforms - no pack", func() {
			ep := dir.Loc{
				SourcePath: getTestPath("mta_with_zipped_module"),
				TargetPath: getResultPath(),
				Descriptor: "dep",
			}
			mNoPlatforms := mta.Module{
				Name: "node-js",
				Path: "node-js",
				BuildParams: map[string]interface{}{
					buildops.SupportedPlatformsParam: []string{},
				},
			}
			Ω(packModule(&ep, false, &mNoPlatforms, "node-js", "cf")).Should(Succeed())
			Ω(getTestPath("result", "mta_with_zipped_module_mta_build_tmp", "node-js", "data.zip")).
				ShouldNot(BeAnExistingFile())
		})

	})

	var _ = Describe("Zip", func() {
		var _ = Describe("Sanity", func() {

			It("Deployment descriptor - Copy only", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: "dep",
				}
				Ω(zipModule(&ep, true, &m, "node-js", "cf")).Should(Succeed())
				Ω(getTestPath("result", ".mta_with_zipped_module_mta_build_tmp", "node-js", "data.zip")).Should(BeAnExistingFile())
			})

			// ep.GetTargetModuleDir(moduleName)
			It("Wrong source", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_unknown"),
					TargetPath: getResultPath(),
					Descriptor: "dev",
				}
				Ω(zipModule(&ep, false, &m, "node-js", "cf")).Should(HaveOccurred())
			})
			It("Target directory exists as a file", func() {
				ep := dir.Loc{
					SourcePath: getTestPath("mta_with_zipped_module"),
					TargetPath: getResultPath(),
					Descriptor: "dev",
				}
				os.MkdirAll(filepath.Join(ep.GetTarget(), ".mta_with_zipped_module_mta_build_tmp"), os.ModePerm)
				createFile("result", ".mta_with_zipped_module_mta_build_tmp", "node-js")
				Ω(zipModule(&ep, false, &m, "node-js", "cf")).Should(HaveOccurred())
			})
		})

		It("No platforms - no pack", func() {
			ep := dir.Loc{
				SourcePath: getTestPath("mta_with_zipped_module"),
				TargetPath: getResultPath(),
				Descriptor: "dep",
			}
			mNoPlatforms := mta.Module{
				Name: "node-js",
				Path: "node-js",
				BuildParams: map[string]interface{}{
					buildops.SupportedPlatformsParam: []string{},
				},
			}
			Ω(zipModule(&ep, false, &mNoPlatforms, "node-js", "cf")).Should(Succeed())
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
				Ω(buildModule("", &ep, &ep, false, "node-js", "cf")).Should(Succeed())
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
				Ω(buildModule("", &ep, &ep, false, "node-js", "cf")).Should(HaveOccurred())
			})

			It("Target folder exists as a file - dev", func() {
				os.MkdirAll(getTestPath("result", ".mta_mta_build_tmp"), os.ModePerm)
				ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
				createFile("result", ".mta_mta_build_tmp", "node-js")
				Ω(buildModule("", &ep, &ep, false, "node-js", "cf")).Should(HaveOccurred())
			})

			It("Target folder exists as a file - dep", func() {
				os.MkdirAll(getTestPath("result", "mta"), os.ModePerm)
				ep := dir.Loc{
					SourcePath:  getTestPath("mta"),
					TargetPath:  getResultPath(),
					Descriptor:  "dep",
					MtaFilename: "mta.yaml",
				}
				createFile("result", "mta", "node-js")
				Ω(buildModule("", &ep, &ep, true, "node-js", "cf")).Should(HaveOccurred())
			})

			It("Deployment Descriptor", func() {
				ep := dir.Loc{
					SourcePath:  getTestPath("mta_with_zipped_module"),
					TargetPath:  getResultPath(),
					MtaFilename: "mta.yaml",
					Descriptor:  "dep"}
				Ω(buildModule("", &ep, &ep, true, "node-js", "cf")).Should(Succeed())
				Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
			})

			var _ = DescribeTable("Invalid inputs", func(projectName, mtaFilename, moduleName string) {
				ep := dir.Loc{SourcePath: getTestPath(projectName), TargetPath: getResultPath(), MtaFilename: mtaFilename}
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
				Ω(buildModule("", &ep, &ep, false, moduleName, "cf")).Should(HaveOccurred())
				Ω(ep.GetTargetTmpDir()).ShouldNot(BeADirectory())
			},
				Entry("Invalid path to application", "mta1", "mta.yaml", "node-js"),
				Entry("Invalid module name", "mta", "mta.yaml", "xxx"),
				Entry("Invalid module name wrong build params", "mtahtml5", "mtaWithWrongBuildParams.yaml", "ui5app"),
			)
		})
	})

	var _ = Describe("copyModuleArchive", func() {

		It("Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module"), TargetPath: getResultPath()}
			Ω(copyModuleArchive(&ep, "node-js", "node-js")).Should(Succeed())
			Ω(ep.GetTargetModuleZipPath("node-js")).Should(BeAnExistingFile())
		})
		It("Invalid - no zip exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta"), TargetPath: getResultPath()}
			Ω(copyModuleArchive(&ep, "node-js", "node-js")).Should(HaveOccurred())
		})
		It("Target directory exists as file", func() {
			ep := dir.Loc{SourcePath: getTestPath("mta_with_zipped_module"), TargetPath: getResultPath()}
			os.MkdirAll(getTestPath("result", ".mta_with_zipped_module_mta_build_tmp"), os.ModePerm)
			createFile("result", ".mta_with_zipped_module_mta_build_tmp", "node-js")
			Ω(copyModuleArchive(&ep, "node-js", "node-js")).Should(HaveOccurred())
		})
	})

	var _ = Describe("CopyMtaContent", func() {
		var source string
		defaultDeploymentDescriptorParam := "dep"
		defaultDeploymentDescriptorName := "mtad.yaml"
		BeforeEach(func() {
			source, _ = ioutil.TempDir("", "testing-mta-content")
		})
		It("Without no deployment descriptor in the source directory", func() {
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, true, os.Getwd)
			Ω(err).Should(Not(BeNil()))
			fmt.Println(err.Error())
			Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(`failed to read the "%s%smtad.yaml" file: open %s%smtad.yaml`, source, pathSep, source, pathSep)))
		})
		It("Location initialization fails", func() {
			err := CopyMtaContent("", source, defaultDeploymentDescriptorParam, false, func() (string, error) {
				return "", errors.New("error")
			})
			Ω(err).Should(Not(BeNil()))
			fmt.Println(err.Error())
			Ω(err.Error()).Should(Equal("copying the MTA content failed during the initialization of deployment descriptor location: failed to initialize the location when getting working directory: error"))
		})
		It("With a deployment descriptor in the source directory with only modules paths as zip archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 0, map[string]string{}, map[string]string{"test-module-0": "zip", "test-module-1": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, true, os.Getwd)
			Ω(err).Should((BeNil()))
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with one module path and one resource path as zip archive and a folder", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 1, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-module-0": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, true, os.Getwd)
			Ω(err).Should((BeNil()))
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only resources with zip and module archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 0, 2, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-resource-1": "folder"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, true, os.Getwd)
			Ω(err).Should((BeNil()))
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true, "test-content": true}, true)).Should(Equal(true))
		})
		It("With a deployment descriptor in the source directory with only resources with zip and module archives", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 2, map[string]string{}, map[string]string{"test-resource-0": "zip", "test-resource-1": "zip", "test-module-0": "zip", "test-module-1": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, false, os.Getwd)
			Ω(err).Should((BeNil()))
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{"." + info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(true))
			Ω(dirContainsAllElements(filepath.Join(source, "."+info.Name()+dir.TempFolderSuffix), map[string]bool{"test.zip": true}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with only one module with zip and one requiredDependency with folder", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{"test-module-0": "test-required"}, map[string]string{"test-module-0": "folder", "test-required": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, false, os.Getwd)
			Ω(err).Should((BeNil()))
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
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, true, os.Getwd)
			Ω(err).Should(HaveOccurred())
		})

		It("With a deployment descriptor in the source directory with only one module with non-existing content", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 1, 0, map[string]string{}, map[string]string{"test-module-0": "not-existing-contet"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, false, os.Getwd)
			Ω(err).Should(Not(BeNil()))
			Ω(err.Error()).Should(Equal(`"not-existing-content" does not exist in the MTA project location`))
			info, _ := os.Stat(source)
			Ω(dirContainsAllElements(source, map[string]bool{info.Name() + dir.TempFolderSuffix: true}, false)).Should(Equal(false))
			Ω(dirContainsAllElements(filepath.Join(source, info.Name()+dir.TempFolderSuffix), map[string]bool{}, true)).Should(Equal(true))
		})

		It("With a deployment descriptor in the source directory with a module with non-existing content and another which has content", func() {
			createFileInGivenPath(filepath.Join(source, defaultDeploymentDescriptorName))
			mta := generateTestMta(source, 2, 0, map[string]string{}, map[string]string{"test-module-0": "not-existing-contet", "test-module-1": "zip"})
			mtaBytes, _ := yaml.Marshal(mta)
			ioutil.WriteFile(filepath.Join(source, defaultDeploymentDescriptorName), mtaBytes, os.ModePerm)
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, false, os.Getwd)
			Ω(err).Should(Not(BeNil()))
			Ω(err.Error()).Should(Equal(`"not-existing-content" does not exist in the MTA project location`))
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
			err := CopyMtaContent(source, source, defaultDeploymentDescriptorParam, false, os.Getwd)
			Ω(err).Should((BeNil()))
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

func createFileInGivenPath(path string) {
	file, err := os.Create(path)
	Ω(err).Should(Succeed())
	file.Close()
}

func createFile(path ...string) {
	createFileInGivenPath(getTestPath(path...))
}
