package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/matchers/internal/miter"
)

var _ = Describe("HaveKeyWithValue", func() {
	var (
		stringKeys map[string]int
		intKeys    map[int]string
		objKeys    map[*myCustomType]*myCustomType

		customA *myCustomType
		customB *myCustomType
	)
	BeforeEach(func() {
		stringKeys = map[string]int{"foo": 2, "bar": 3}
		intKeys = map[int]string{2: "foo", 3: "bar"}

		customA = &myCustomType{s: "a", n: 2, f: 2.3, arr: []string{"ice", "cream"}}
		customB = &myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"cake"}}
		objKeys = map[*myCustomType]*myCustomType{customA: customA, customB: customA}
	})

	When("passed a map", func() {
		It("should do the right thing", func() {
			Expect(stringKeys).Should(HaveKeyWithValue("foo", 2))
			Expect(stringKeys).ShouldNot(HaveKeyWithValue("foo", 1))
			Expect(stringKeys).ShouldNot(HaveKeyWithValue("baz", 2))
			Expect(stringKeys).ShouldNot(HaveKeyWithValue("baz", 1))

			Expect(intKeys).Should(HaveKeyWithValue(2, "foo"))
			Expect(intKeys).ShouldNot(HaveKeyWithValue(4, "foo"))
			Expect(intKeys).ShouldNot(HaveKeyWithValue(2, "baz"))

			Expect(objKeys).Should(HaveKeyWithValue(customA, customA))
			Expect(objKeys).Should(HaveKeyWithValue(&myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"cake"}}, &myCustomType{s: "a", n: 2, f: 2.3, arr: []string{"ice", "cream"}}))
			Expect(objKeys).ShouldNot(HaveKeyWithValue(&myCustomType{s: "b", n: 4, f: 3.1, arr: []string{"apple", "pie"}}, customA))
		})
	})

	When("passed a correctly typed nil", func() {
		It("should operate successfully on the passed in value", func() {
			var nilMap map[int]string
			Expect(nilMap).ShouldNot(HaveKeyWithValue("foo", "bar"))
		})
	})

	When("the passed in key or value is actually a matcher", func() {
		It("should pass each element through the matcher", func() {
			Expect(stringKeys).Should(HaveKeyWithValue(ContainSubstring("oo"), 2))
			Expect(intKeys).Should(HaveKeyWithValue(2, ContainSubstring("oo")))
			Expect(stringKeys).ShouldNot(HaveKeyWithValue(ContainSubstring("foobar"), 2))
		})

		It("should fail if the matcher ever fails", func() {
			actual := map[int]string{1: "a", 3: "b", 2: "c"}
			success, err := (&HaveKeyWithValueMatcher{Key: ContainSubstring("ar"), Value: 2}).Match(actual)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			otherActual := map[string]int{"a": 1, "b": 2, "c": 3}
			success, err = (&HaveKeyWithValueMatcher{Key: "a", Value: ContainSubstring("1")}).Match(otherActual)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})

		It("succeeds if any matching key has a matching value, even if the matchers error on other entries, and errors only if nothing matches (https://github.com/onsi/gomega/issues/926)", func() {
			keyErrors := map[any]int{"a": 1, 2: 2, "b": 3, "c": 4}
			valueErrors := map[string]any{"aFoo": "x", "bFoo": 2, "cFoo": "y", "dFoo": "z"}
			for range 100 {
				success, err := (&HaveKeyWithValueMatcher{Key: BeNumerically(">", 1), Value: 2}).Match(keyErrors)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(success).Should(BeTrue())

				success, err = (&HaveKeyWithValueMatcher{Key: BeNumerically(">", 1), Value: 3}).Match(keyErrors)
				Expect(err).Should(MatchError(ContainSubstring("HaveKeyWithValue's key matcher failed with")))
				Expect(success).Should(BeFalse())

				success, err = (&HaveKeyWithValueMatcher{Key: HaveSuffix("Foo"), Value: BeNumerically(">", 1)}).Match(valueErrors)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(success).Should(BeTrue())

				success, err = (&HaveKeyWithValueMatcher{Key: HaveSuffix("Foo"), Value: BeNumerically(">", 2)}).Match(valueErrors)
				Expect(err).Should(MatchError(ContainSubstring("HaveKeyWithValue's value matcher failed with")))
				Expect(success).Should(BeFalse())
			}
		})

		It("succeeds if any matching key has a matching value (https://github.com/onsi/gomega/issues/929)", func() {
			actual := map[string]string{"aFoo": "x", "bFoo": "Bar", "cFoo": "y", "dFoo": "z"}
			for range 100 {
				Expect(actual).Should(HaveKeyWithValue(MatchRegexp(`.+Foo$`), "Bar"))
				Expect(actual).ShouldNot(HaveKeyWithValue(MatchRegexp(`.+Foo$`), "Baz"))
			}
		})
	})

	When("passed something that is not a map", func() {
		It("should error", func() {
			success, err := (&HaveKeyWithValueMatcher{Key: "foo", Value: "bar"}).Match([]string{"foo"})
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())

			success, err = (&HaveKeyWithValueMatcher{Key: "foo", Value: "bar"}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("iterators", func() {
		BeforeEach(func() {
			if !miter.HasIterators() {
				Skip("iterators not available")
			}
		})

		When("passed an iter.Seq2", func() {
			It("should do the right thing", func() {
				Expect(universalMapIter2).Should(HaveKeyWithValue("foo", 0))
				Expect(universalMapIter2).ShouldNot(HaveKeyWithValue("foo", 1))
				Expect(universalMapIter2).ShouldNot(HaveKeyWithValue("baz", 2))
				Expect(universalMapIter2).ShouldNot(HaveKeyWithValue("baz", 1))

				Expect(universalMapIter2).Should(HaveKeyWithValue("bar", 42))
				Expect(universalMapIter2).Should(HaveKeyWithValue("baz", 666))

				Expect(universalMapIter2).ShouldNot(HaveKeyWithValue("bar", "abc"))
				Expect(universalMapIter2).ShouldNot(HaveKeyWithValue(555, "abc"))
			})
		})

		When("passed a correctly typed nil", func() {
			It("should operate successfully on the passed in value", func() {
				var nilIter2 func(func(string, int) bool)
				Expect(nilIter2).ShouldNot(HaveKeyWithValue("foo", 0))
			})
		})

		When("the passed in key or value is actually a matcher", func() {
			It("should pass each element through the matcher", func() {
				Expect(universalMapIter2).Should(HaveKeyWithValue(ContainSubstring("oo"), BeNumerically("<", 1)))
				Expect(universalMapIter2).Should(HaveKeyWithValue(ContainSubstring("foo"), 0))
			})

			It("should fail if the matcher ever fails", func() {
				success, err := (&HaveKeyWithValueMatcher{Key: "bar", Value: ContainSubstring("argh")}).Match(universalMapIter2)
				Expect(success).Should(BeFalse())
				Expect(err).Should(HaveOccurred())

				success, err = (&HaveKeyWithValueMatcher{Key: "foo", Value: ContainSubstring("1")}).Match(universalMapIter2)
				Expect(success).Should(BeFalse())
				Expect(err).Should(HaveOccurred())
			})

			It("succeeds if any matching key has a matching value, even if the matchers error on other entries, and errors only if nothing matches (https://github.com/onsi/gomega/issues/926)", func() {
				keyErrorFirst := func(yield func(any, int) bool) {
					_ = yield("a", 1) && yield(2, 2) && yield("b", 3)
				}
				keyMatchFirst := func(yield func(any, int) bool) {
					_ = yield(2, 2) && yield("a", 1) && yield("b", 3)
				}
				for _, actual := range []any{keyErrorFirst, keyMatchFirst} {
					success, err := (&HaveKeyWithValueMatcher{Key: BeNumerically(">", 1), Value: 2}).Match(actual)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(success).Should(BeTrue())

					success, err = (&HaveKeyWithValueMatcher{Key: BeNumerically(">", 1), Value: 3}).Match(actual)
					Expect(err).Should(MatchError(ContainSubstring("HaveKeyWithValue's key matcher failed with")))
					Expect(success).Should(BeFalse())
				}

				valueErrorFirst := func(yield func(string, any) bool) {
					_ = yield("aFoo", "x") && yield("bFoo", 2) && yield("cFoo", "y")
				}
				valueMatchFirst := func(yield func(string, any) bool) {
					_ = yield("bFoo", 2) && yield("aFoo", "x") && yield("cFoo", "y")
				}
				for _, actual := range []any{valueErrorFirst, valueMatchFirst} {
					success, err := (&HaveKeyWithValueMatcher{Key: HaveSuffix("Foo"), Value: BeNumerically(">", 1)}).Match(actual)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(success).Should(BeTrue())

					success, err = (&HaveKeyWithValueMatcher{Key: HaveSuffix("Foo"), Value: BeNumerically(">", 2)}).Match(actual)
					Expect(err).Should(MatchError(ContainSubstring("HaveKeyWithValue's value matcher failed with")))
					Expect(success).Should(BeFalse())
				}
			})

			It("succeeds if any matching key has a matching value (https://github.com/onsi/gomega/issues/929)", func() {
				matchingLast := func(yield func(string, string) bool) {
					_ = yield("aFoo", "x") && yield("cFoo", "y") && yield("bFoo", "Bar")
				}
				matchingFirst := func(yield func(string, string) bool) {
					_ = yield("bFoo", "Bar") && yield("aFoo", "x") && yield("cFoo", "y")
				}
				Expect(matchingLast).Should(HaveKeyWithValue(MatchRegexp(`.+Foo$`), "Bar"))
				Expect(matchingFirst).Should(HaveKeyWithValue(MatchRegexp(`.+Foo$`), "Bar"))
				Expect(matchingLast).ShouldNot(HaveKeyWithValue(MatchRegexp(`.+Foo$`), "Baz"))
				Expect(matchingFirst).ShouldNot(HaveKeyWithValue(MatchRegexp(`.+Foo$`), "Baz"))
			})
		})

		When("passed something that is not an iter.Seq2", func() {
			It("should error", func() {
				success, err := (&HaveKeyWithValueMatcher{Key: "foo", Value: "bar"}).Match(universalIter)
				Expect(success).Should(BeFalse())
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
