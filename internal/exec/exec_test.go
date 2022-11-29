package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
			Ω(Execute(args, false)).Should(Succeed())
		},
			Entry("EchoTesting", [][]string{{"", "sh", "-c", `echo -n {"Name": "Bob", "Age": 32}`}}),
			Entry("Dummy Go Testing", [][]string{{"", "go", "test", "exec_dummy_test.go"}}))

		var _ = DescribeTable("Invalid input", func(args [][]string) {
			Ω(Execute(args, false)).Should(HaveOccurred())
		},
			Entry("Valid command fails on input", [][]string{{"", "go", "test", "exec_unknown_test.go"}}),
			Entry("Invalid command", [][]string{{"", "dateXXX"}}),
		)
	})

	var _ = DescribeTable("executeCommand Failures",
		func(cmd *exec.Cmd) {
			Ω(executeCommand(cmd, make(chan struct{}), true)).Should(HaveOccurred())
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
		Entry("returns default timeout when timeout is empty", "", "10m0s", false),
		Entry("returns error for bad timeout", "abc", "", true),
	)

	var executeTester = func(executor func() error, minSeconds, maxSeconds int, isError bool, expectedTimeout string) {
		start := time.Now()
		err := executor()
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
	}

	DescribeTable("ExecuteWithTimeout",
		func(args [][]string, timeout string, minSeconds, maxSeconds int, isError bool, expectedTimeout string) {
			executeTester(func() error {
				return ExecuteWithTimeout(args, timeout, true)
			}, minSeconds, maxSeconds, isError, expectedTimeout)
		},
		Entry("succeeds when timeout wasn't reached", [][]string{{"", "sh", "-c", "sleep 2"}}, "10s", 2, 5, false, ""),
		Entry("fails when timeout was reached", [][]string{{"", "sh", "-c", "sleep 5"}}, "2s", 2, 3, true, "2s"),
		Entry("fails when timeout was reached in the second command",
			[][]string{{"", "sh", "-c", "sleep 2"}, {"", "sh", "-c", "sleep 3"}}, "4s", 4, 5, true, "4s"),
	)

	It("ExecuteWithTimeout fails when timeout value is invalid", func() {
		err := ExecuteWithTimeout([][]string{{"sh", "-c", "sleep 1"}}, "1234", true)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(ExecInvalidTimeoutMsg, "1234")))
	})

	DescribeTable("ExecuteCommandsWithTimeout",
		func(args []string, timeout string, minSeconds, maxSeconds int, isError bool, expectedTimeout string) {
			executeTester(func() error {
				return ExecuteCommandsWithTimeout(args, timeout, "", true)
			}, minSeconds, maxSeconds, isError, expectedTimeout)
		},
		Entry("succeeds when timeout wasn't reached", []string{`sh -c "sleep 2"`}, "10s", 2, 5, false, ""),
		Entry("fails when timeout was reached", []string{`sh -c 'sleep 5'`}, "2s", 2, 3, true, "2s"),
		Entry("fails when timeout was reached in the second command",
			[]string{`sh -c "sleep 2"`, `sh -c 'sleep 3'`}, "4s", 4, 5, true, "4s"),
	)

	Describe("ExecuteCommandsWithTimeout tests with cleanup", func() {
		wd, _ := os.Getwd()
		path := filepath.Join(wd, "testdata")
		AfterEach(func() {
			Ω(os.RemoveAll(filepath.Join(path, "b.txt"))).Should(Succeed())
		})
		It("ExecuteCommandsWithTimeout is executed in the requested directory", func() {
			Ω(ExecuteCommandsWithTimeout([]string{`sh -c 'cp a.txt b.txt'`}, "10m", path, true)).Should(Succeed())
			Ω(filepath.Join(path, "b.txt")).Should(BeAnExistingFile())
		})
	})

	It("ExecuteCommandsWithTimeout fails when timeout value is invalid", func() {
		err := ExecuteCommandsWithTimeout([]string{`sh -c "sleep 1"`}, "1234", ".", true)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(ExecInvalidTimeoutMsg, "1234")))
	})
})
