package dir

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"cloud-mta-build-tool/cmd/logs"
	"time"
	"strings"
)

func TestCreateDirIfNotExist(t *testing.T) {
	tests := []struct {
		name      string
		dirName   string
		validator func(t *testing.T, dirName string, err error)
	}{
		{
			name:    "SanityTest",
			dirName: filepath.Join(GetPath(), "testdata", "level2", "result"),
			validator: func(t *testing.T, dirName string, err error) {
				assert.Nil(t, err)
				err = os.RemoveAll(dirName)
				assert.Nil(t, err)
			},
		},
		{
			name:    "DirectoryExists",
			dirName: filepath.Join(GetPath(), "testdata", "level2", "level3"),
			validator: func(t *testing.T, dirName string, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name:    "BadDirectoryName",
			dirName: filepath.Join(GetPath(), "testdata", "level2", "/"),
			validator: func(t *testing.T, dirName string, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateDirIfNotExist(tt.dirName);
			tt.validator(t, tt.dirName, err)
		})
	}
}

func TestArchive(t *testing.T) {
	type args struct {
		params []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Archive(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Archive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		validator func(t *testing.T, filename string, file *os.File, err error)
	}{
		{
			name:     "SanityTest",
			filename: filepath.Join(GetPath(), "testdata", "level2", "newFile"),
			validator: func(t *testing.T, filename string, file *os.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				file.Close()
				err = os.Remove(filename)
				assert.Nil(t, err)
			},
		},
		{
			name:     "BadNameOfFile",
			filename: filepath.Join(GetPath(), "testdata", "level2", "/.txt"),
			validator: func(t *testing.T, filename string, file *os.File, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	logs.NewLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := CreateFile(tt.filename)
			tt.validator(t, tt.filename, file, err)
		})
	}
}

func countFilesInDir(name string) int {
	files, _ := ioutil.ReadDir(name)
	return len(files)
}

func TestCopyDir(t *testing.T) {

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name         string
		args         args
		preprocessor func(t *testing.T, args args)
		validator    func(t *testing.T, args args, err error)
	}{
		{
			name:         "SanityTest",
			args:         args{filepath.Join(GetPath(), "testdata", "level2"), filepath.Join(GetPath(), "testdata", "result")},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.Nil(t, err)
				assert.Equal(t, countFilesInDir(args.src), countFilesInDir(args.dst))
				os.RemoveAll(args.dst)
			},
		},
		{
			name:         "SourceDirectoryDoesNotExist",
			args:         args{filepath.Join(GetPath(), "testdata", "level5"), filepath.Join(GetPath(), "testdata", "result")},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:         "SourceIsNotDirectory",
			args:         args{filepath.Join(GetPath(), "testdata", "level2", "level2_one.txt"), filepath.Join(GetPath(), "testdata", "result")},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:         "DstDirectoryNotValid",
			args:         args{filepath.Join(GetPath(), "testdata", "level2"), "/"},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	logs.Logger = logs.NewLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preprocessor(t, tt.args)
			err := CopyDir(tt.args.src, tt.args.dst)
			tt.validator(t, tt.args, err)
		})
	}
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

func Test_copyEntries(t *testing.T) {
	srcPath := filepath.Join(GetPath(), "testdata", "level2", "level3")
	dstPath := filepath.Join(GetPath(), "testdata", "result")
	os.MkdirAll(dstPath, os.ModePerm)
	files, _ := ioutil.ReadDir(srcPath)
	var filesWrapped [2]os.FileInfo
	for i, file := range files {
		filesWrapped[i] = testFile{file: file}
	}
	copyEntries(filesWrapped[:], srcPath, dstPath)
	assert.Equal(t, countFilesInDir(srcPath)-1, countFilesInDir(dstPath))
	os.RemoveAll(dstPath)

	dstPath = filepath.Join(GetPath(), "testdata", "//")
	err := copyEntries(filesWrapped[:], filepath.Join(GetPath(), "testdata", "level2", "levelx"), dstPath)
	assert.NotNil(t, err)
}

func Test_copyFile(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name      string
		args      args
		validator func(t *testing.T, args args, err error)
	}{
		{
			name: "SourceNotExists",
			args: args{filepath.Join(GetPath(), "testdata", "fileSrc"), filepath.Join(GetPath(), "testdata", "fileDst")},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "SourceIsDirectory",
			args: args{filepath.Join(GetPath(), "testdata", "level2"), filepath.Join(GetPath(), "testdata", "level2", "fileDst")},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
				os.RemoveAll(args.dst)
			},
		},
		{
			name: "WrongDestinationName",
			args: args{filepath.Join(GetPath(), "testdata", "level2", "level2_one.txt"), filepath.Join(GetPath(), "testdata", "level2", "/")},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "DestinationExists",
			args: args{filepath.Join(GetPath(), "testdata", "level2", "level3", "level3_one.txt"), filepath.Join(GetPath(), "testdata", "level2", "level3", "level3_two.txt")},
			validator: func(t *testing.T, args args, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyFile(tt.args.src, tt.args.dst)
			tt.validator(t, tt.args, err)
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		validator func(t *testing.T, filename string, fileContent []byte, err error)
	}{
		{
			name:     "SanityTest",
			filename: filepath.Join(GetPath(), "testdata", "level2", "level2_one.txt"),
			validator: func(t *testing.T, filename string, fileContent []byte, err error) {
				assert.Nil(t, err)
				s := string(fileContent)
				assert.Equal(t, "level2_one.txt", s)
			},
		},
		{
			name:     "FileNotExists",
			filename: filepath.Join(GetPath(), "testdata", "level2", "level2_one__.txt"),
			validator: func(t *testing.T, filename string, fileContent []byte, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	logs.NewLogger()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileContent, err := Load(tt.filename);
			tt.validator(t, tt.filename, fileContent, err)

		})
	}
}
