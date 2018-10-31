package logs

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test logger")
}

func logLevelWithPanic() {
	logLevel("other")
}

var _ = Describe("logger", func() {
	AfterEach(func() {
		Logger = nil
	})

	DescribeTable("Log Level", func(input string, expected logrus.Level) {
		level := logLevel(input)
		Ω(level).To(Equal(expected))
	},
		Entry("Debug Level", "debug", logrus.DebugLevel),
		Entry("Info Level", "info", logrus.InfoLevel),
		Entry("Error Level", "error", logrus.ErrorLevel),
		Entry("Warn Level", "warn", logrus.WarnLevel),
		Entry("Fatal Level", "fatal", logrus.FatalLevel),
		Entry("Panic Level", "panic", logrus.PanicLevel),
	)

	It("should Panic", func() {
		Ω(logLevelWithPanic).Should(Panic())
	})

	It("should return a pointer to a Logger object ('info' level)", func() {
		logger := NewLogger()

		Expect(logger.Out).To(Equal(os.Stdout))
		Expect(logger.Level).To(Equal(logLevel(getLogLevel())))
		Expect(Logger).To(Equal(logger))
	})

})
