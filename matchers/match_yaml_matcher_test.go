package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchYAMLMatcher", func() {
	Context("When passed stringifiables", func() {
		It("should succeed if the YAML matches", func() {
			Expect("---").Should(MatchYAML(""))
			Expect("a: 1").Should(MatchYAML(`{"a":1}`))
			Expect("a: 1\nb: 2").Should(MatchYAML(`{"b":2, "a":1}`))
		})

		It("should explain if the YAML does not match when it should", func() {
			message := (&MatchYAMLMatcher{YAMLToMatch: "a: 1"}).FailureMessage("b: 2")
			Expect(message).To(MatchRegexp(`Expected\s+<string>: b: 2\s+to match YAML of\s+<string>: a: 1`))
		})

		It("should normalise the expected and actual when explaining if the YAML does not match when it should", func() {
			message := (&MatchYAMLMatcher{YAMLToMatch: "a: 'one'"}).FailureMessage("{b: two}")
			Expect(message).To(MatchRegexp(`Expected\s+<string>: b: two\s+to match YAML of\s+<string>: a: one`))
		})

		It("should explain if the YAML matches when it should not", func() {
			message := (&MatchYAMLMatcher{YAMLToMatch: "a: 1"}).NegatedFailureMessage("a: 1")
			Expect(message).To(MatchRegexp(`Expected\s+<string>: a: 1\s+not to match YAML of\s+<string>: a: 1`))
		})

		It("should normalise the expected and actual when explaining if the YAML matches when it should not", func() {
			message := (&MatchYAMLMatcher{YAMLToMatch: "a: 'one'"}).NegatedFailureMessage("{a: one}")
			Expect(message).To(MatchRegexp(`Expected\s+<string>: a: one\s+not to match YAML of\s+<string>: a: one`))
		})

		It("should fail if the YAML does not match", func() {
			Expect("a: 1").ShouldNot(MatchYAML(`{"b":2, "a":1}`))
		})

		It("should work with byte arrays", func() {
			Expect([]byte("a: 1")).Should(MatchYAML([]byte("a: 1")))
			Expect("a: 1").Should(MatchYAML([]byte("a: 1")))
			Expect([]byte("a: 1")).Should(MatchYAML("a: 1"))
		})
	})

	Context("with multi-document YAML streams (https://github.com/onsi/gomega/issues/933)", func() {
		It("should compare every document, not just the first", func() {
			Expect("a: 1\n---\nb: 2").ShouldNot(MatchYAML("a: 1"))
			Expect("a: 1").ShouldNot(MatchYAML("a: 1\n---\nb: 2"))
			Expect("a: 1\n---\nb: 3").ShouldNot(MatchYAML("a: 1\n---\nb: 2"))
			Expect("a: 2\n---\nb: 2").ShouldNot(MatchYAML("a: 1\n---\nb: 2"))
			Expect("{a: 1}\n---\n{b: 2}\n").Should(MatchYAML("a: 1\n---\nb: 2"))
		})

		It("should not confuse a stream of documents with a single document holding a list", func() {
			Expect("- a: 1\n- b: 2").ShouldNot(MatchYAML("a: 1\n---\nb: 2"))
			Expect("a: 1\n---\nb: 2").ShouldNot(MatchYAML("- a: 1\n- b: 2"))
		})

		It("should ignore empty documents, such as those created by leading or trailing document separators", func() {
			Expect("---\na: 1").Should(MatchYAML("a: 1"))
			Expect("a: 1\n---\n").Should(MatchYAML("a: 1"))
			Expect("---\na: 1\n---\n# nothing to see here\n").Should(MatchYAML("a: 1"))
			Expect("---\na: 1\n---\n---\nb: 2\n---\n").Should(MatchYAML("a: 1\n---\nb: 2"))
			Expect("---\n---\n").Should(MatchYAML(""))
		})

		It("should treat explicit null documents as documents", func() {
			Expect("a: 1\n---\n~").ShouldNot(MatchYAML("a: 1"))
			Expect("a: 1\n---\nnull").Should(MatchYAML("a: 1\n---\n~"))
			Expect("null").Should(MatchYAML(""))
		})

		It("should error if any document is invalid", func() {
			success, err := (&MatchYAMLMatcher{YAMLToMatch: "a: 1"}).Match("a: 1\n---\ngood:\nbad")
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError(ContainSubstring("Actual 'a: 1\n---\ngood:\nbad' should be valid YAML")))

			success, err = (&MatchYAMLMatcher{YAMLToMatch: "a: 1\n---\ngood:\nbad"}).Match("a: 1")
			Expect(success).Should(BeFalse())
			Expect(err).Should(MatchError(ContainSubstring("Expected 'a: 1\n---\ngood:\nbad' should be valid YAML")))
		})

		It("should show every document, and say which document mismatched, when explaining a failure", func() {
			matcher := &MatchYAMLMatcher{YAMLToMatch: "a: 1\n---\nb: {c: 2}"}
			Expect(matcher.Match("a: 1\n---\nb: {c: 3}")).To(BeFalse())
			Expect(matcher.FailureMessage("a: 1\n---\nb: {c: 3}")).To(Equal(`Expected
    <string>: a: 1
    ---
    b:
        c: 3
to match YAML of
    <string>: a: 1
    ---
    b:
        c: 2

first mismatched document: 2 (counting from 1)

first mismatched key: "b"."c"`))
		})

		It("should not mention documents when explaining a failure between single documents", func() {
			matcher := &MatchYAMLMatcher{YAMLToMatch: "b: {c: 2}"}
			Expect(matcher.Match("---\nb: {c: 3}")).To(BeFalse())
			Expect(matcher.FailureMessage("---\nb: {c: 3}")).To(Equal(`Expected
    <string>: b:
        c: 3
to match YAML of
    <string>: b:
        c: 2

first mismatched key: "b"."c"`))
		})
	})

	When("the expected is not valid YAML", func() {
		It("should error and explain why", func() {
			success, err := (&MatchYAMLMatcher{YAMLToMatch: ""}).Match("good:\nbad")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Actual 'good:\nbad' should be valid YAML"))
		})
	})

	When("the actual is not valid YAML", func() {
		It("should error and explain why", func() {
			success, err := (&MatchYAMLMatcher{YAMLToMatch: "good:\nbad"}).Match("")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Expected 'good:\nbad' should be valid YAML"))
		})

		It("errors when passed directly to Message", func() {
			Expect(func() {
				matcher := MatchYAMLMatcher{YAMLToMatch: "good"}
				matcher.FailureMessage("good:\nbad")
			}).To(Panic())
		})
	})

	When("the expected is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchYAMLMatcher{YAMLToMatch: 2}).Match("")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchYAMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <int>: 2"))

			success, err = (&MatchYAMLMatcher{YAMLToMatch: nil}).Match("")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchYAMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <nil>: nil"))
		})
	})

	When("the actual is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchYAMLMatcher{YAMLToMatch: ""}).Match(2)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchYAMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <int>: 2"))

			success, err = (&MatchYAMLMatcher{YAMLToMatch: ""}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchYAMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <nil>: nil"))
		})
	})
})
