package ext

import (
	"cloud-mta-build-tool/cmd/mta/models"
	"reflect"
	"testing"
)

func TestExeCmd(t *testing.T) {
	type args struct {
		m models.Modules
	}
	tests := []struct {
		name string
		args args
		want []Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExeCmd(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExeCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
