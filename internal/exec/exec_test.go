package exec

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

type testStr struct {
}

func (s *testStr) Read(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func (s *testStr) Write(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func (s *testStr) Close() error {
	return errors.New("err")
}

var _ = Describe("Execute", func() {

	var _ = Describe("Execute call", func() {

		var _ = DescribeTable("Valid input", func(args [][]string) {
			Ω(Execute(args)).Should(Succeed())
		},
			Entry("EchoTesting", [][]string{{"", "echo", "-n", `{"Name": "Bob", "Age": 32}`}}),
			Entry("Dummy Go Testing", [][]string{{"", "go", "test", "exec_dummy_test.go"}}))

		var _ = DescribeTable("Invalid input", func(args [][]string) {
			Ω(Execute(args)).Should(HaveOccurred())
		},
			Entry("Valid command fails on input", [][]string{{"", "go", "test", "exec_unknown_test.go"}}),
			Entry("Invalid command", [][]string{{"", "dateXXX"}}),
		)
	})

	var _ = DescribeTable("executeCommand Failures",
		func(cmd *exec.Cmd) {
			Ω(executeCommand(cmd)).Should(HaveOccurred())
		},

		Entry("fails on StdoutPipe", &exec.Cmd{Stdout: &testStr{}}),
		Entry("fails on StderrPipe", &exec.Cmd{Stderr: &testStr{}}),
	)

	It("Indicator", func() {
		// var wg sync.WaitGroup
		// wg.Add(1)
		shutdownCh := make(chan struct{})
		start := time.Now()
		go indicator(shutdownCh)
		time.Sleep(3 * time.Second)
		// close(shutdownCh)
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
		}

		shutdownCh <- struct{}{}
		// wg.Wait()
	})

})
