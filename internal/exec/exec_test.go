package exec

import (
	"fmt"
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
			Entry("EchoTesting", [][]string{{"", "bash", "-c", `echo -n {"Name": "Bob", "Age": 32}`}}),
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
			Ω(executeCommand(cmd, make(chan struct{}))).Should(HaveOccurred())
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

	var _ = DescribeTable("parseTimeoutString",
		func(timeout, expectedTimeout string, isError bool) {
			duration, err := parseTimeoutString(timeout)
			if isError {
				Ω(err).Should(HaveOccurred())
			} else {
				Ω(err).Should(Succeed())
				Ω(duration.String()).Should(Equal(expectedTimeout))
			}

		},
		Entry("parses timeout with seconds", "3s", "3s", false),
		Entry("parses timeout with minutes", "10m", "10m0s", false),
		Entry("parses timeout with hours", "5h", "5h0m0s", false),
		Entry("parses timeout with mixed time units", "10m3s", "10m3s", false),
		Entry("returns default timeout when timeout is empty", "", "5m0s", false),
		Entry("returns error for bad timeout", "abc", "", true),
	)

	DescribeTable("ExecuteWithTimeout",
		func(args [][]string, timeout string, minSeconds, maxSeconds int, isError bool, expectedTimeout string) {
			start := time.Now()
			err := ExecuteWithTimeout(args, timeout)
			elapsed := time.Since(start)
			// Check error
			if isError {
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(ExecTimeoutMsg, expectedTimeout)))
			} else {
				Ω(err).Should(Succeed())
			}

			// Check elapsed time
			Ω(elapsed).Should(BeNumerically(">=", time.Duration(minSeconds)*time.Second))
			Ω(elapsed).Should(BeNumerically("<=", time.Duration(maxSeconds)*time.Second))
		},
		Entry("succeeds when timeout wasn't reached", [][]string{{"", "bash", "-c", "sleep 2"}}, "10s", 2, 5, false, ""),
		Entry("fails when timeout was reached", [][]string{{"", "bash", "-c", "sleep 5"}}, "2s", 2, 3, true, "2s"),
		Entry("fails when timeout was reached in the second command",
			[][]string{{"", "bash", "-c", "sleep 2"}, {"", "bash", "-c", "sleep 3"}}, "4s", 4, 5, true, "4s"),
	)

	It("ExecuteWithTimeout fails when timeout value is invalid", func() {
		err := ExecuteWithTimeout([][]string{{"bash", "-c", "sleep 1"}}, "1234")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(ExecInvalidTimeoutMsg, "1234")))
	})
})
