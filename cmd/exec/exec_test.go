package exec

import (
	"cloud-mta-build-tool/cmd/logs"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestExecute_WithEchoTesting(t *testing.T) {
	logs.NewLogger()
	cdParams := [][]string{{"", "echo", "-n", `{"Name": "Bob", "Age": 32}`}}
	Execute(cdParams)
}

func Test_Execute_WithGoTesting(t *testing.T) {
	logs.NewLogger()
	cdParams := [][]string{{"", "go", "test", "exec_dummy_test.go"}}
	err := Execute(cdParams)
	assert.NilError(t, err)
}

func Test_Execute_WithGoTestingNegative(t *testing.T) {
	logs.NewLogger()
	cdParams := [][]string{{"", "go", "test", "exec_unknown_test.go"}}
	err := Execute(cdParams)
	assert.Equal(t, err != nil, true)

	cdParams = [][]string{{"", "dateXXX"}}
	err = Execute(cdParams)
	assert.Equal(t, err != nil, true)
}

func Test_Indicator(t *testing.T) {

	shutdownCh := make(chan struct{})
	start := time.Now()
	go indicator(shutdownCh)
	time.Sleep(1 * time.Second)
	close(shutdownCh)
	sec := time.Since(start).Seconds()
	switch int(sec) {
	case 0:
		// Output:
	case 1:
		// Output: .
	case 2:
		// Output: ..
	case 3:
		// Output: ...
	default:
		t.Error("Sleeping time is more than 3 seconds")
	}


}
