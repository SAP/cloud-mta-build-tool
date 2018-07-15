package builders

import (
	"reflect"
	"testing"
)

func TestBuildProcess(t *testing.T) {
	type args struct {
		options []func(*BuildCfg)
	}
	tests := []struct {
		name         string
		args         args
		wantBuildcfg *BuildCfg
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuildcfg, err := BuildProcess(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuildcfg, tt.wantBuildcfg) {
				t.Errorf("BuildProcess() = %v, want %v", gotBuildcfg, tt.wantBuildcfg)
			}
		})
	}
}

func TestTargetDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TargetDir(); got != tt.want {
				t.Errorf("TargetDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
