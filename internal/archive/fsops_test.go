package dir

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

func getFullPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, filepath.Join(relPath...))
}

type testMtaYamlStr struct {
	fullpath string
	path     string
	err      error
}

func (t *testMtaYamlStr) GetMtaYamlFilename() string {
	return t.fullpath
}

func (t *testMtaYamlStr) GetMtaYamlPath() string {
	return t.path
}

func (t *testMtaYamlStr) GetMtaExtYamlPath(platform string) string {
	return t.fullpath
}

var _ = Describe("FSOPS", func() {

	var _ = Describe("CreateDir", func() {

		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "level2", "result"))
		})

		var _ = DescribeTable("CreateDir", func(dirPath string) {

			Ω(CreateDirIfNotExist(dirPath)).Should(Succeed())
		},
			Entry("Sanity", getFullPath("testdata", "level2", "result")),
			Entry("DirectoryExists", getFullPath("testdata", "level2", "level3")),
		)
		It("Fails because file with the same name exists", func() {
			err := CreateDirIfNotExist(getFullPath("testdata", "level2", "result"))
			if err != nil {
				fmt.Println("error occurred during directory creation")
			}
			file, _ := os.Create(getFullPath("testdata", "level2", "result", "file"))
			file.Close()
			Ω(CreateDirIfNotExist(getFullPath("testdata", "level2", "result", "file"))).Should(HaveOccurred())
		})
	})

	var _ = Describe("Archive", func() {
		var targetFilePath = getFullPath("testdata", "arch.mbt")

		BeforeEach(func() {
			fileInfoProvider = &mockFileInfoProvider{}
		})

		AfterEach(func() {
			fileInfoProvider = &standardFileInfoProvider{}
			Ω(os.RemoveAll(targetFilePath)).Should(Succeed())
		})

		var _ = DescribeTable("Archive", func(source, target string, ignore []string, fails bool, expectedFiles []string) {
			err := Archive(source, target, ignore)
			if fails {
				Ω(err).Should(HaveOccurred())
			} else {
				Ω(err).Should(Succeed())
				Ω(target).Should(BeAnExistingFile())
				validateArchiveContents(expectedFiles, target)
			}
		},
			Entry("Sanity",
				getFullPath("testdata", "mtahtml5"), targetFilePath, nil, false, []string{
					"mta.sh", "mta.yaml",
					"ui5app/", "ui5app/Gruntfile.js",
					"ui5app/webapp/", "ui5app/webapp/Component.js", "ui5app/webapp/index.html",
					"ui5app/webapp/controller/", "ui5app/webapp/controller/View1.controller.js",
					"ui5app/webapp/css/", "ui5app/webapp/css/style.css",
					"ui5app/webapp/i18n/", "ui5app/webapp/i18n/i18n.properties",
					"ui5app/webapp/model/", "ui5app/webapp/model/models.js",
					"ui5app/webapp/view/", "ui5app/webapp/view/View1.view.xml",
				}),
			Entry("Target is folder",
				getFullPath("testdata", "mtahtml5"), getFullPath("testdata"), nil, true, nil),
			Entry("Source is broken symbolic link",
				getFullPath("testdata", "testsymlink", "symlink_broken"), targetFilePath, nil, true, nil),
			Entry("Sanity - ignore folder",
				getFullPath("testdata", "testproject"), targetFilePath, []string{"ui5app/"}, false, []string{
					"cf-mtaext.yaml", "mta.sh", "mta.yaml",
				}),
			Entry("Sanity - ignore file",
				getFullPath("testdata", "testproject"), targetFilePath, []string{"ui5app/Gr*.js"}, false, []string{
					"cf-mtaext.yaml", "mta.sh", "mta.yaml",
					"ui5app/",
					"ui5app/webapp/", "ui5app/webapp/Component.js", "ui5app/webapp/index.html",
					"ui5app/webapp/controller/", "ui5app/webapp/controller/View1.controller.js",
					"ui5app/webapp/model/", "ui5app/webapp/model/models.js",
					"ui5app/webapp/view/", "ui5app/webapp/view/View1.view.xml",
				}),
			Entry("SourceIsNotFolder",
				getFullPath("testdata", "level2", "level2_one.txt"), targetFilePath, nil, false, []string{"level2_one.txt"}),
			Entry("Target is empty string",
				getFullPath("testdata", "mtahtml5"), "", nil, true, nil),
			Entry("Source is empty string", "", "", nil, true, nil),
			// folder module (which is symbolic link itself) is archived
			// it points to folder moduleNew that consists of symlink pointing to content and package.json
			// etc... thus we check a complex case consisting of normal files/folders and symbolic links that are also
			// files and folders
			Entry("symbolic links",
				getFullPath("testdata", "testsymlink", "symlink_dir_to_moduleNew"), targetFilePath, nil, false,
				[]string{"symlink_dir_to_content/", "package.json",
					"symlink_dir_to_content/test_dir/", "symlink_dir_to_content/test_dir/test1.txt",
					"symlink_dir_to_content/test.txt",
					"symlink_dir_to_content/symlink_dir_to_another_content/", "symlink_dir_to_content/symlink_dir_to_another_content/test3.txt",
					"symlink_dir_to_content/symlink_dir_to_another_content/symlink_to_test4.txt"}),
			Entry("symbolic links with ignore",
				getFullPath("testdata", "testsymlink", "symlink_dir_to_moduleNew"), targetFilePath, []string{"symlink_dir_to_content"}, false,
				[]string{"package.json"}),
		)
	})

	var _ = Describe("utils", func() {
		BeforeEach(func() {
			fileInfoProvider = &mockFileInfoProvider{}
		})

		AfterEach(func() {
			fileInfoProvider = &standardFileInfoProvider{}
		})

		var _ = Describe("getIgnoredEntries", func() {
			It("source path is wrong", func() {
				_, err := getIgnoredEntries([]string{"x"}, getFullPath("testdata", "notexists"))
				Ω(err).Should(HaveOccurred())
			})

			It("ignored symlink to folder", func() {
				entries, err := getIgnoredEntries([]string{"symlink_dir_to_content"}, getFullPath("testdata", "testsymlink", "symlink_dir_to_moduleNew"))
				Ω(err).Should(Succeed())
				Ω(len(entries)).Should(Equal(1))
				_, ok := entries[getFullPath("testdata", "testsymlink", "moduleNew", "symlink_dir_to_content")]
				Ω(ok).Should(BeTrue())
			})

			It("ignored symlink to file", func() {
				entries, err := getIgnoredEntries([]string{filepath.Join("another_content", "symlink_to_test4.txt")}, getFullPath("testdata", "testsymlink"))
				Ω(err).Should(Succeed())
				Ω(len(entries)).Should(Equal(1))
				_, ok := entries[getFullPath("testdata", "testsymlink", "another_content", "symlink_to_test4.txt")]
				Ω(ok).Should(BeTrue())
			})

			It("ignored recursive symlink", func() {
				_, err := getIgnoredEntries([]string{"x"},
					getFullPath("testdata", "testsymlink", "dir_with_recursive_symlink", "subdir", "symlink_dir_to_sibling"))
				Ω(err).Should(HaveOccurred())
			})
		})

		var _ = Describe("dereferenceSymlink", func() {
			It("wrong file path", func() {
				_, _, _, err := dereferenceSymlink(getFullPath("testdata", "notexists"), make(map[string]bool))
				Ω(err).Should(HaveOccurred())
			})
			It("wrong relative path in symlink", func() {
				fileInfoProvider = &mockFileInfoProvider{ReturnRelativePath: true}
				_, _, _, err := dereferenceSymlink(getFullPath("testdata", "testsymlink", "symlink_broken"), make(map[string]bool))
				Ω(err).Should(HaveOccurred())
			})
			It("existing relative path in symlink", func() {
				fileInfoProvider = &mockFileInfoProvider{ReturnRelativePath: true}
				path, _, _, err := dereferenceSymlink(getFullPath("testdata", "testsymlink", "symlink_dir_to_moduleNew"), make(map[string]bool))
				Ω(err).Should(Succeed())
				Ω(path).Should(Equal(getFullPath("testdata", "testsymlink", "moduleNew")))
			})
		})
	})

	var _ = Describe("addSymbolicLinkToArchive - failures", func() {
		var archive *zip.Writer
		var zipFile *os.File
		var err error

		BeforeEach(func() {
			fileInfoProvider = &mockFileInfoProvider{}
			Ω(CreateDirIfNotExist(getFullPath("testdata", "result"))).Should(Succeed())
			zipFile, err = CreateFile(getFullPath("testdata", "result", "arch.zip"))
			Ω(err).Should(Succeed())
			archive = zip.NewWriter(zipFile)
		})

		AfterEach(func() {
			fileInfoProvider = &standardFileInfoProvider{}
			Ω(archive.Close()).Should(Succeed())
			Ω(zipFile.Close()).Should(Succeed())
			Ω(os.RemoveAll(getFullPath("testdata", "result"))).Should(Succeed())
		})

		It("not a symbolic link", func() {
			Ω(addSymbolicLinkToArchive(getFullPath("testdata", "testsymlink", "test4.txt"),
				getFullPath("testdata", "testsymlink"), "", "", nil,
				make(map[string]bool), nil)).Should(HaveOccurred())
		})
		It("broken symbolic link (points to the deleted folder)", func() {
			Ω(addSymbolicLinkToArchive(getFullPath("testdata", "testsymlink", "symlink_broken"),
				getFullPath("testdata", "testsymlink"), "", "", nil,
				make(map[string]bool), nil)).Should(HaveOccurred())
		})
		It("link to folder with broken symbolic link", func() {
			Ω(addSymbolicLinkToArchive(getFullPath("testdata", "testsymlink", "symlink_dir_to_symlink_dir_broken"),
				getFullPath("testdata", "testsymlink", "symlink_dir_to_symlink_dir_broken"), "", "", nil,

				make(map[string]bool), nil)).Should(HaveOccurred())
		})
		var _ = DescribeTable("recursive symbolic link", func(relPath ...string) {
			path := getFullPath("testdata", "testsymlink")
			for _, pathElement := range relPath {
				path = filepath.Join(path, pathElement)
			}

			err := addSymbolicLinkToArchive(path, getFullPath("testdata", "testsymlink"), "", "", archive,
				make(map[string]bool), nil)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal(fmt.Sprintf(recursiveSymLinkMsg, path)))
		},
			Entry("recursion to itself", "symlink_to_itself"),
			Entry("2 steps recursion", "symlink_recursion_2step_a"),
			Entry("3 steps recursion", "symlink_recursion_3step_a"),
			Entry("sibling folders with recursion", "dir_with_recursive_symlink", "subdir", "symlink_dir_to_sibling"),
			Entry("recursion to upper folder", "dir_with_recursive_symlink", "subdir", "symlink_dir_recursion_to_parent_dir"))
	})

	var _ = Describe("Create File", func() {
		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "result.txt"))
		})
		It("Sanity", func() {
			file, err := CreateFile(getFullPath("testdata", "result.txt"))
			Ω(getFullPath("testdata", "result.txt")).Should(BeAnExistingFile())
			file.Close()
			Ω(err).Should(BeNil())
		})
		It("Fails on empty path", func() {
			_, err := CreateFile("")
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Describe("CopyDir", func() {
		var targetPath = getFullPath("testdata", "result")
		AfterEach(func() {
			os.RemoveAll(targetPath)
		})

		It("Sanity - parallel", func() {
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntriesInParallel)).Should(Succeed())
			Ω(countFilesInDir(targetPath)).Should(Equal(countFilesInDir(sourcePath)))
		})

		It("Sanity - not parallel", func() {
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntries)).Should(Succeed())
			Ω(countFilesInDir(targetPath)).Should(Equal(countFilesInDir(sourcePath)))
		})

		It("TargetFileLocked", func() {
			f, _ := os.Create(targetPath)
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntries)).Should(HaveOccurred())
			f.Close()
		})

		It("TargetFileLocked", func() {
			f, _ := os.Create(targetPath)
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntriesInParallel)).Should(HaveOccurred())
			f.Close()
		})

		var _ = DescribeTable("Invalid cases", func(source, target string) {
			Ω(CopyDir(source, targetPath, true, CopyEntries)).Should(HaveOccurred())
		},
			Entry("SourceDirectoryDoesNotExist", getFullPath("testdata", "level5"), targetPath),
			Entry("SourceIsNotDirectory", getFullPath("testdata", "level2", "level2_one.txt"), targetPath),
			Entry("DstDirectoryNotValid", getFullPath("level2"), ":"),
		)

		var _ = DescribeTable("Invalid cases - parallel", func(source, target string) {
			Ω(CopyDir(source, targetPath, true, CopyEntriesInParallel)).Should(HaveOccurred())
		},
			Entry("SourceDirectoryDoesNotExist", getFullPath("testdata", "level5"), targetPath),
			Entry("SourceIsNotDirectory", getFullPath("testdata", "level2", "level2_one.txt"), targetPath),
			Entry("DstDirectoryNotValid", getFullPath("level2"), ":"),
		)

		var _ = DescribeTable("Copy File - Invalid", func(source, target string, matcher GomegaMatcher) {
			Ω(CopyFile(source, target)).Should(matcher)
		},
			Entry("SourceNotExists", getFullPath("testdata", "fileSrc"), targetPath, HaveOccurred()),
			Entry("SourceIsDirectory", getFullPath("testdata", "level2"), targetPath, HaveOccurred()),
			Entry("WrongDestinationName", getFullPath("testdata", "level2", "level2_one.txt"), getFullPath("testdata", "level2", "/"), HaveOccurred()),
			Entry("DestinationExists", getFullPath("testdata", "level2", "level3", "level3_one.txt"), getFullPath("testdata", "level2", "level3", "level3_two.txt"), Succeed()),
		)
		var _ = DescribeTable("Copy File - Invalid", func(source, target string, matcher GomegaMatcher) {
			Ω(CopyFileWithMode(source, target, os.ModePerm)).Should(matcher)
		},
			Entry("TargetPathEmpty", getFullPath("testdata", "fileSrc"), "", HaveOccurred()),
			Entry("SourceIsDirectory", getFullPath("testdata", "level2"), targetPath, HaveOccurred()),
			Entry("DestinationExists", getFullPath("testdata", "level2", "level3", "level3_one.txt"), getFullPath("testdata", "level2", "level3", "level3_two.txt"), Succeed()),
		)
	})

	var _ = Describe("Copy Entries", func() {

		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "result"))
		})

		It("Sanity", func() {
			sourcePath := getFullPath("testdata", "level2", "level3")
			targetPath := getFullPath("testdata", "result")
			err := CreateDirIfNotExist(targetPath)
			if err != nil {
				fmt.Println("error occurred during dir creation")
			}
			files, _ := ioutil.ReadDir(sourcePath)
			// Files wrapped to overwrite their methods
			var filesWrapped []os.FileInfo
			Ω(CopyEntries(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			for _, file := range files {
				filesWrapped = append(filesWrapped, testFile{file: file})
			}
			Ω(CopyEntries(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			Ω(countFilesInDir(sourcePath) - 1).Should(Equal(countFilesInDir(targetPath)))
			os.RemoveAll(targetPath)
		})
		It("Sanity - copy in parallel", func() {
			sourcePath := getFullPath("testdata", "level2", "level3")
			targetPath := getFullPath("testdata", "result")
			err := CreateDirIfNotExist(targetPath)
			if err != nil {
				fmt.Println("error occurred during dir creation")
			}
			files, _ := ioutil.ReadDir(sourcePath)
			// Files wrapped to overwrite their methods
			var filesWrapped []os.FileInfo
			Ω(CopyEntriesInParallel(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			for _, file := range files {
				filesWrapped = append(filesWrapped, testFile{file: file})
			}
			Ω(CopyEntriesInParallel(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			Ω(countFilesInDir(sourcePath) - 1).Should(Equal(countFilesInDir(targetPath)))
			os.RemoveAll(targetPath)
		})
	})

	var _ = Describe("Copy By Patterns", func() {

		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "result"))
		})

		var _ = DescribeTable("Valid Cases", func(modulePath string, patterns, expectedFiles []string) {
			sourcePath := getFullPath("testdata", "testbuildparams", modulePath)
			targetPath := getFullPath("testdata", "result")
			Ω(CopyByPatterns(sourcePath, targetPath, patterns)).Should(Succeed())
			for _, file := range expectedFiles {
				Ω(file).Should(BeAnExistingFile())
			}
		},
			Entry("Single file", "ui2",
				[]string{"deep/folder/inui2/anotherfile.txt"},
				[]string{getFullPath("testdata", "result", "anotherfile.txt")}),
			Entry("Wildcard for 2 files", "ui2",
				[]string{"deep/*/inui2/another*"},
				[]string{getFullPath("testdata", "result", "anotherfile.txt"),
					getFullPath("testdata", "result", "anotherfile2.txt")}),
			Entry("Wildcard for 2 files - dot start", "ui2",
				[]string{"./deep/*/inui2/another*"},
				[]string{getFullPath("testdata", "result", "anotherfile.txt"),
					getFullPath("testdata", "result", "anotherfile2.txt")}),
			Entry("Specific folder of second level", "ui2",
				[]string{"*/folder/*"},
				[]string{
					getFullPath("testdata", "result", "inui2", "anotherfile.txt"),
					getFullPath("testdata", "result", "inui2", "anotherfile2.txt")}),
			Entry("All", "ui1",
				[]string{"*"},
				[]string{getFullPath("testdata", "result", "webapp", "Component.js")}),
			Entry("Dot", "ui1",
				[]string{"."},
				[]string{getFullPath("testdata", "result", "ui1", "webapp", "Component.js")}),
			Entry("Multiple patterns", "ui2", //
				[]string{"deep/folder/inui2/anotherfile.txt", "*/folder/"},
				[]string{
					getFullPath("testdata", "result", "folder", "inui2", "anotherfile.txt"),
					getFullPath("testdata", "result", "anotherfile.txt")}),
			Entry("Empty patterns", "ui2",
				[]string{},
				[]string{}),
		)

		var _ = DescribeTable("Invalid Cases", func(targetPath, modulePath string, patterns []string) {
			sourcePath := getFullPath("testdata", "testbuildparams", modulePath)
			err := CopyByPatterns(sourcePath, targetPath, patterns)
			Ω(err).Should(HaveOccurred())
		},
			Entry("Target path relates to file ",
				getFullPath("testdata", "testbuildparams", "mta.yaml"), "ui2",
				[]string{"deep/folder/inui2/somefile.txt"}),
			Entry("Wrong pattern ",
				getFullPath("testdata", "result"), "ui2", []string{"[a,b"}),
			Entry("Empty target path", "", "ui2", []string{"[a,b"}),
		)
	})

	var _ = Describe("getRelativePath", func() {
		It("not relative path", func() {
			Ω(getRelativePath(getFullPath("sss", "abc", "xyz", "fff"),
				getFullPath("uuu"))).Should(Equal(getFullPath("sss", "abc", "xyz", "fff")))
		})

		It("non empty base path", func() {
			Ω(getRelativePath(getFullPath("abc", "xyz", "fff"),
				getFullPath())).Should(Equal(filepath.Join("abc", "xyz", "fff")))
		})

		It("empty base path", func() {
			Ω(getRelativePath(getFullPath("abc", "xyz", "fff"), "")).Should(Equal(getFullPath("abc", "xyz", "fff")))
		})
	})

	It("copyByPattern - fails because target is file", func() {
		Ω(copyByPattern(getPath("testdata", "mtahtml5"),
			getFullPath("testdata", "level2", "level2_one.txt"), "m*")).Should(HaveOccurred())
	})

	It("copyEntries - fails on entry with empty name", func() {
		Ω(copyEntries([]string{""}, getPath("testdata", "mtahtml5"),
			getFullPath("testdata", "level2", "level2_one.txt"), "m*")).Should(HaveOccurred())
	})

	It("changeTargetMode - fails if source does not exist", func() {
		Ω(changeTargetMode(getPath("testdata", "not-exists"), getPath("testdata", "not-exists-2"))).
			Should(HaveOccurred())
	})

	var _ = Describe("Read", func() {
		It("Sanity", func() {
			test := testMtaYamlStr{
				fullpath: getFullPath("testdata", "testproject", "mta.yaml"),
				path:     getFullPath("testdata", "testproject", "mta.yaml"),
				err:      nil,
			}
			res, resErr := Read(&test)
			Ω(res).ShouldNot(BeNil())
			Ω(resErr).Should(BeNil())
		})
	})

	var _ = Describe("ReadExt", func() {
		It("Sanity", func() {
			test := testMtaYamlStr{
				fullpath: getFullPath("testdata", "testext", "mta.yaml"),
				path:     getFullPath("testdata", "testext", "mta.yaml"),
				err:      nil,
			}
			res, resErr := ReadExt(&test, "cf-mtaext.yaml")
			Ω(res).ShouldNot(BeNil())
			Ω(resErr).Should(BeNil())
		})
	})

	var _ = DescribeTable("CloseFile", func(toFail bool, errorArg error, expectedErr error) {
		testFile := testCloser{fail: toFail}
		if expectedErr == nil {
			Ω(CloseFile(&testFile, errorArg)).Should(BeNil())
		} else {
			Ω(CloseFile(&testFile, errorArg).Error()).Should(Equal(expectedErr.Error()))
		}
	},
		Entry("No error", false, nil, nil),
		Entry("Original error only", false, errors.New("original error"), errors.New("original error")),
		Entry("New error", true, nil, errors.New("failed to close")),
		Entry("Original and new errors", true, errors.New("original error"), errors.New("original error")),
	)

	var _ = Describe("fileInfoProvider", func() {
		It("isSymbolicLink", func() {
			fileInfo, err := os.Stat(getFullPath("testdata"))
			Ω(err).Should(Succeed())
			Ω(fileInfoProvider.isSymbolicLink(fileInfo)).Should(BeFalse())
		})
		It("isDir", func() {
			fileInfo, err := os.Stat(getFullPath("testdata"))
			Ω(err).Should(Succeed())
			Ω(fileInfoProvider.isDir(fileInfo)).Should(BeTrue())
		})
		It("readlink", func() {
			_, err := fileInfoProvider.readlink(getFullPath("testdata"))
			Ω(err).Should(HaveOccurred())
		})
		It("stat", func() {
			_, err := fileInfoProvider.stat(getFullPath("testdata"))
			Ω(err).Should(Succeed())
		})
	})
})

func countFilesInDir(name string) int {
	files, _ := ioutil.ReadDir(name)
	return len(files)
}

type testFile struct {
	file os.FileInfo
}

func (file testFile) Name() string {
	return file.file.Name()
}

func (file testFile) Size() int64 {
	return file.file.Size()
}

func (file testFile) Mode() os.FileMode {
	if strings.Contains(file.file.Name(), "level3_one.txt") {
		return os.ModeSymlink
	}
	return file.file.Mode()
}

func (file testFile) ModTime() time.Time {
	return file.file.ModTime()
}

func (file testFile) IsDir() bool {
	return file.file.IsDir()
}

func (file testFile) Sys() interface{} {
	return nil
}

type testCloser struct {
	fail bool
}

func (f *testCloser) Close() error {
	if f.fail {
		return errors.New("failed to close")
	}
	return nil
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

func contains(element string, elements []string) bool {
	for _, el := range elements {
		if el == element {
			return true
		}
	}
	return false
}

type mockFileInfoProvider struct {
	ReturnRelativePath bool
}

func (provider *mockFileInfoProvider) isSymbolicLink(file os.FileInfo) bool {
	return strings.HasPrefix(file.Name(), "symlink_")
}

func (provider *mockFileInfoProvider) isDir(file os.FileInfo) bool {
	if provider.isSymbolicLink(file) {
		return strings.HasPrefix(file.Name(), "symlink_dir_")
	}
	return file.IsDir()
}

func (provider *mockFileInfoProvider) readlink(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	text := scanner.Text()
	if text == "error" {
		scanner.Scan()
		return "", errors.New(scanner.Text())
	}
	textSplit := strings.Split(text, "/")
	if provider.ReturnRelativePath {
		return filepath.Join(textSplit...), nil
	}
	// Resolve path
	fullPath := filepath.Dir(path)
	for _, textElement := range textSplit {
		fullPath = filepath.Join(fullPath, textElement)
	}
	return fullPath, nil
}

func (provider *mockFileInfoProvider) stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
