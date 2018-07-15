package builders

import (
	"reflect"
	"testing"
)

func TestNpmBuilder_Path(t *testing.T) {
	type fields struct {
		path string
		name string
		dir  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NpmBuilder{
				path: tt.fields.path,
				name: tt.fields.name,
				dir:  tt.fields.dir,
			}
			if got := n.Path(); got != tt.want {
				t.Errorf("NpmBuilder.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpmBuilder_Build(t *testing.T) {
	type fields struct {
		path string
		name string
		dir  string
	}
	type args struct {
		pdir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NpmBuilder{
				path: tt.fields.path,
				name: tt.fields.name,
				dir:  tt.fields.dir,
			}
			if err := n.Build(tt.args.pdir); (err != nil) != tt.wantErr {
				t.Errorf("NpmBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewNPMBuilder(t *testing.T) {
	type args struct {
		p string
		n string
	}
	tests := []struct {
		name string
		args args
		want *NpmBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNPMBuilder(tt.args.p, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNPMBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_npmSeq(t *testing.T) {
	type args struct {
		modPath string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := npmSeq(tt.args.modPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("npmSeq() = %v, want %v", got, tt.want)
			}
		})
	}
}
