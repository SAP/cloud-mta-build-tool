package dir

import (
	"os"
	"reflect"
	"testing"
)

func TestCreateDirIfNotExist(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDirIfNotExist(tt.args.dir); got != tt.want {
				t.Errorf("CreateDirIfNotExist() = %v, want %v", got, tt.want)
			}
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
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	type args struct {
		src string
		dst string
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
			if err := CopyDir(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_copyFile(t *testing.T) {
	type args struct {
		src string
		dst string
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
			if err := copyFile(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultTempDirFunc(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultTempDirFunc(tt.args.path); got != tt.want {
				t.Errorf("DefaultTempDirFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
