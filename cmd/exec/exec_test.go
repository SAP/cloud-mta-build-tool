package exec

import (
	"testing"

	"cloud-mta-build-tool/cmd/logs"
)

func TestExecute_WithEchoTesting(t *testing.T) {
	logs.NewLogger()

	cdParams := [][]string{{"", "echo", "-n", `{"Name": "Bob", "Age": 32}`}}
	Execute(cdParams)
}

func Test_Execute_WithGoTesting(t *testing.T) {
	logs.NewLogger()
	cdParams := [][]string{{"", "go", "test", "exec_dummy_test.go"}}
	Execute(cdParams)
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
