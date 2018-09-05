package exec

import "testing"

func TestExecute(t *testing.T) {
	type args struct {
		cmdParams [][]string
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Execute(tt.args.cmdParams); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_indicator(t *testing.T) {
	type args struct {
		shutdownCh <-chan struct{}
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indicator(tt.args.shutdownCh)
		})
	}
}
