package logs

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()
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

	It("should return Debug Level", func() {
		level := logLevel("debug")
		Expect(level).To(Equal(logrus.DebugLevel))
	})
	It("should return Info Level", func() {
		level := logLevel("info")
		Expect(level).To(Equal(logrus.InfoLevel))
	})
	It("should return Error Level", func() {
		level := logLevel("error")
		Expect(level).To(Equal(logrus.ErrorLevel))
	})
	It("should return Warn Level", func() {
		level := logLevel("warn")
		Expect(level).To(Equal(logrus.WarnLevel))
	})
	It("should return Fatal Level", func() {
		level := logLevel("fatal")
		Expect(level).To(Equal(logrus.FatalLevel))
	})
	It("should return Panic Level", func() {
		level := logLevel("panic")
		Expect(level).To(Equal(logrus.PanicLevel))
	})
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
