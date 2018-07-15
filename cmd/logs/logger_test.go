package logs

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("logger", func() {
	Describe("logger", func() {
		AfterEach(func() {
			Logger = nil
		})

		Describe("#NewLogger", func() {
			It("should return a pointer to a Logger object ('info' level)", func() {
				logger := NewLogger()

				Expect(logger.Out).To(Equal(os.Stderr))
				Expect(logger.Level).To(Equal(logrus.InfoLevel))
				Expect(Logger).To(Equal(logger))
			})

			It("should return a pointer to a Logger object ('debug' level)", func() {
				logger := NewLogger()

				Expect(logger.Out).To(Equal(os.Stderr))
				Expect(logger.Level).To(Equal(logrus.DebugLevel))
				Expect(Logger).To(Equal(logger))
			})

			It("should return a pointer to a Logger object ('error' level)", func() {
				logger := NewLogger()

				Expect(logger.Out).To(Equal(os.Stderr))
				Expect(logger.Level).To(Equal(logrus.ErrorLevel))
				Expect(Logger).To(Equal(logger))
			})
		})

	})
})
