package builders

import (
	"reflect"
	"testing"
)

func TestGetPath(t *testing.T) {
	tests := []struct {
		name    string
		wantDir string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := GetPath(); gotDir != tt.wantDir {
				t.Errorf("GetPath() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	type args struct {
		b         Builder
		toPath    string
		mkTempDir TempDirFunc
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
			if err := Build(tt.args.b, tt.args.toPath, tt.args.mkTempDir); (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
