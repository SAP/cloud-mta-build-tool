package gen

import "testing"

func TestGenerate(t *testing.T) {
	type args struct {
		path string
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Generate(tt.args.path)
		})
	}
}
