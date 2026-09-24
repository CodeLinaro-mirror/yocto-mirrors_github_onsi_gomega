package matchers_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("BeNumerically", func() {
	When("passed a number", func() {
		It("should support ==", func() {
			Expect(uint32(5)).Should(BeNumerically("==", 5))
			Expect(float64(5.0)).Should(BeNumerically("==", 5))
			Expect(int8(5)).Should(BeNumerically("==", 5))
		})

		It("should not have false positives", func() {
			Expect(5.1).ShouldNot(BeNumerically("==", 5))
			Expect(5).ShouldNot(BeNumerically("==", 5.1))
		})

		It("should support >", func() {
			Expect(uint32(5)).Should(BeNumerically(">", 4))
			Expect(float64(5.0)).Should(BeNumerically(">", 4.9))
			Expect(int8(5)).Should(BeNumerically(">", 4))

			Expect(uint32(5)).ShouldNot(BeNumerically(">", 5))
			Expect(float64(5.0)).ShouldNot(BeNumerically(">", 5.0))
			Expect(int8(5)).ShouldNot(BeNumerically(">", 5))
		})

		It("should support <", func() {
			Expect(uint32(5)).Should(BeNumerically("<", 6))
			Expect(float64(5.0)).Should(BeNumerically("<", 5.1))
			Expect(int8(5)).Should(BeNumerically("<", 6))

			Expect(uint32(5)).ShouldNot(BeNumerically("<", 5))
			Expect(float64(5.0)).ShouldNot(BeNumerically("<", 5.0))
			Expect(int8(5)).ShouldNot(BeNumerically("<", 5))
		})

		It("should support >=", func() {
			Expect(uint32(5)).Should(BeNumerically(">=", 4))
			Expect(float64(5.0)).Should(BeNumerically(">=", 4.9))
			Expect(int8(5)).Should(BeNumerically(">=", 4))

			Expect(uint32(5)).Should(BeNumerically(">=", 5))
			Expect(float64(5.0)).Should(BeNumerically(">=", 5.0))
			Expect(int8(5)).Should(BeNumerically(">=", 5))

			Expect(uint32(5)).ShouldNot(BeNumerically(">=", 6))
			Expect(float64(5.0)).ShouldNot(BeNumerically(">=", 5.1))
			Expect(int8(5)).ShouldNot(BeNumerically(">=", 6))
		})

		It("should support <=", func() {
			Expect(uint32(5)).Should(BeNumerically("<=", 6))
			Expect(float64(5.0)).Should(BeNumerically("<=", 5.1))
			Expect(int8(5)).Should(BeNumerically("<=", 6))

			Expect(uint32(5)).Should(BeNumerically("<=", 5))
			Expect(float64(5.0)).Should(BeNumerically("<=", 5.0))
			Expect(int8(5)).Should(BeNumerically("<=", 5))

			Expect(uint32(5)).ShouldNot(BeNumerically("<=", 4))
			Expect(float64(5.0)).ShouldNot(BeNumerically("<=", 4.9))
			Expect(int8(5)).Should(BeNumerically("<=", 5))
		})

		When("passed ~", func() {
			When("passed a float", func() {
				Context("and there is no precision parameter", func() {
					It("should default to 1e-8", func() {
						Expect(5.00000001).Should(BeNumerically("~", 5.00000002))
						Expect(5.00000001).ShouldNot(BeNumerically("~", 5.0000001))
					})

					It("should show failure message", func() {
						actual := BeNumerically("~", 4.98).FailureMessage(123)
						expected := "Expected\n    <int>: 123\nto be ~\n    <float64>: 4.98"
						Expect(actual).To(Equal(expected))
					})

					It("should show negated failure message", func() {
						actual := BeNumerically("~", 4.98).NegatedFailureMessage(123)
						expected := "Expected\n    <int>: 123\nnot to be ~\n    <float64>: 4.98"
						Expect(actual).To(Equal(expected))
					})
				})

				Context("and there is a precision parameter", func() {
					It("should use the precision parameter", func() {
						Expect(5.1).Should(BeNumerically("~", 5.19, 0.1))
						Expect(5.1).Should(BeNumerically("~", 5.01, 0.1))
						Expect(5.1).ShouldNot(BeNumerically("~", 5.22, 0.1))
						Expect(5.1).ShouldNot(BeNumerically("~", 4.98, 0.1))
					})

					It("should show precision in failure message", func() {
						actual := BeNumerically("~", 4.98, 0.1).FailureMessage(123)
						expected := "Expected\n    <int>: 123\nto be within 0.1 of ~\n    <float64>: 4.98"
						Expect(actual).To(Equal(expected))
					})

					It("should show precision in negated failure message", func() {
						actual := BeNumerically("~", 4.98, 0.1).NegatedFailureMessage(123)
						expected := "Expected\n    <int>: 123\nnot to be within 0.1 of ~\n    <float64>: 4.98"
						Expect(actual).To(Equal(expected))
					})
				})
			})

			When("passed an int/uint", func() {
				Context("and there is no precision parameter", func() {
					It("should just do strict equality", func() {
						Expect(5).Should(BeNumerically("~", 5))
						Expect(5).ShouldNot(BeNumerically("~", 6))
						Expect(uint(5)).ShouldNot(BeNumerically("~", 6))
					})
				})

				Context("and there is a precision parameter", func() {
					It("should use precision parameter", func() {
						Expect(5).Should(BeNumerically("~", 6, 2))
						Expect(5).ShouldNot(BeNumerically("~", 8, 2))
						Expect(uint(5)).Should(BeNumerically("~", 6, 1))
					})
				})
			})
		})

		It("compares signed and unsigned integers by value (https://github.com/onsi/gomega/issues/925)", func() {
			Expect(-1).ShouldNot(BeNumerically("==", uint64(math.MaxUint64)))
			Expect(uint64(math.MaxUint64)).ShouldNot(BeNumerically("==", -1))
			Expect(-1).ShouldNot(BeNumerically("~", uint64(math.MaxUint64), 1))
			Expect(uint64(math.MaxUint64)).ShouldNot(BeNumerically("~", -1, uint(1)))

			Expect(uint(5)).Should(BeNumerically(">", -3))
			Expect(uint(5)).Should(BeNumerically(">=", -3))
			Expect(uint(5)).ShouldNot(BeNumerically("<", -3))
			Expect(uint(5)).ShouldNot(BeNumerically("<=", -3))
			Expect(-3).Should(BeNumerically("<", uint(5)))
			Expect(-3).Should(BeNumerically("<=", uint(5)))
			Expect(-3).ShouldNot(BeNumerically(">", uint(5)))
			Expect(-3).ShouldNot(BeNumerically(">=", uint(5)))

			Expect(int64(-1)).Should(BeNumerically("<", uint64(math.MaxUint64)))
			Expect(uint64(math.MaxUint64)).Should(BeNumerically(">", int64(math.MinInt64)))
			Expect(uint64(math.MaxInt64 + 1)).Should(BeNumerically(">", int64(math.MaxInt64)))
			Expect(int64(math.MaxInt64)).Should(BeNumerically("<", uint64(math.MaxInt64+1)))

			Expect(uint(3)).Should(BeNumerically("==", 3))
			Expect(3).Should(BeNumerically("==", uint(3)))
			Expect(uint(2)).Should(BeNumerically("~", -1, 3))
			Expect(uint(2)).ShouldNot(BeNumerically("~", -1, 2))
			Expect(-1).Should(BeNumerically("~", uint(2), uint(3)))
			Expect(-1).ShouldNot(BeNumerically("~", uint(2), uint(2)))
			Expect(uint64(math.MaxUint64)).ShouldNot(BeNumerically("~", -1, uint64(math.MaxUint64)))
			Expect(uint64(math.MaxUint64)).Should(BeNumerically("~", -1, math.Exp2(64)))
			Expect(uint64(math.MaxUint64)).ShouldNot(BeNumerically("~", int64(math.MinInt64), uint64(math.MaxUint64)))

			// a negative threshold never matches, whatever the signedness of the operands
			Expect(uint(5)).ShouldNot(BeNumerically("~", uint(5), -1))
			Expect(5).ShouldNot(BeNumerically("~", 5, -1))
		})
	})

	When("passed a non-number", func() {
		It("should error", func() {
			success, err := (&BeNumericallyMatcher{Comparator: "==", CompareTo: []any{5}}).Match("foo")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&BeNumericallyMatcher{Comparator: "=="}).Match(5)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&BeNumericallyMatcher{Comparator: "~", CompareTo: []any{3.0, "foo"}}).Match(5.0)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("foo"))

			success, err = (&BeNumericallyMatcher{Comparator: "==", CompareTo: []any{"bar"}}).Match(5)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&BeNumericallyMatcher{Comparator: "==", CompareTo: []any{"bar"}}).Match("foo")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&BeNumericallyMatcher{Comparator: "==", CompareTo: []any{nil}}).Match(0)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&BeNumericallyMatcher{Comparator: "==", CompareTo: []any{0}}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	When("passed an unsupported comparator", func() {
		It("should error", func() {
			success, err := (&BeNumericallyMatcher{Comparator: "!=", CompareTo: []any{5}}).Match(4)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})
})
