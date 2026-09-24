package matchers_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchJSONMatcher", func() {
	Context("When passed stringifiables", func() {
		It("should succeed if the JSON matches", func() {
			Expect("{}").Should(MatchJSON("{}"))
			Expect(`{"a":1}`).Should(MatchJSON(`{"a":1}`))
			Expect(`{
			             "a":1
			         }`).Should(MatchJSON(`{"a":1}`))
			Expect(`{"a":1, "b":2}`).Should(MatchJSON(`{"b":2, "a":1}`))
			Expect(`{"a":1}`).ShouldNot(MatchJSON(`{"b":2, "a":1}`))

			Expect(`{"a":"a", "b":"b"}`).ShouldNot(MatchJSON(`{"a":"a", "b":"b", "c":"c"}`))
			Expect(`{"a":"a", "b":"b", "c":"c"}`).ShouldNot(MatchJSON(`{"a":"a", "b":"b"}`))

			Expect(`{"a":null, "b":null}`).ShouldNot(MatchJSON(`{"c":"c", "d":"d"}`))
			Expect(`{"a":null, "b":null, "c":null}`).ShouldNot(MatchJSON(`{"a":null, "b":null, "d":null}`))
		})

		It("should work with byte arrays", func() {
			Expect([]byte("{}")).Should(MatchJSON([]byte("{}")))
			Expect("{}").Should(MatchJSON([]byte("{}")))
			Expect([]byte("{}")).Should(MatchJSON("{}"))
		})

		It("should work with json.RawMessage", func() {
			Expect([]byte(`{"a": 1}`)).Should(MatchJSON(json.RawMessage(`{"a": 1}`)))
		})
	})

	When("a key mismatch is found", func() {
		It("reports the first found mismatch", func() {
			subject := MatchJSONMatcher{JSONToMatch: `5`}
			actual := `7`
			subject.Match(actual)

			failureMessage := subject.FailureMessage(`7`)
			Expect(failureMessage).ToNot(ContainSubstring("first mismatched key"))

			subject = MatchJSONMatcher{JSONToMatch: `{"a": 1, "b.g": {"c": 2, "1": ["hello", "see ya"]}}`}
			actual = `{"a": 1, "b.g": {"c": 2, "1": ["hello", "goodbye"]}}`
			subject.Match(actual)

			failureMessage = subject.FailureMessage(actual)
			Expect(failureMessage).To(ContainSubstring(`first mismatched key: "b.g"."1"[1]`))
		})
	})

	When("the expected is not valid JSON", func() {
		It("should error and explain why", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: `{}`}).Match(`oops`)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Actual 'oops' should be valid JSON"))
		})
	})

	When("the actual is not valid JSON", func() {
		It("should error and explain why", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: `oops`}).Match(`{}`)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Expected 'oops' should be valid JSON"))
		})
	})

	When("the JSON contains numbers too large for a float64", func() {
		It("does not match distinct out-of-range numbers (https://github.com/onsi/gomega/issues/930)", func() {
			success, _ := (&MatchJSONMatcher{JSONToMatch: `1e401`}).Match(`1e400`)
			Expect(success).Should(BeFalse())

			success, _ = (&MatchJSONMatcher{JSONToMatch: `{"a": 1e999}`}).Match(`{"a": 1e400}`)
			Expect(success).Should(BeFalse())

			success, _ = (&MatchJSONMatcher{JSONToMatch: `1e400`}).Match(`-1e400`)
			Expect(success).Should(BeFalse())
		})
	})

	When("comparing numbers", func() {
		It("compares numbers exactly, without losing precision (https://github.com/onsi/gomega/issues/931)", func() {
			Expect(`12345678901234567890`).ShouldNot(MatchJSON(`12345678901234567891`))
			Expect(`{"id": 9007199254740993}`).ShouldNot(MatchJSON(`{"id": 9007199254740992}`))
			Expect(`[9007199254740993]`).ShouldNot(MatchJSON(`[9007199254740992]`))
			Expect(`0.1`).ShouldNot(MatchJSON(`0.10000000000000001`))
			Expect(`1e400`).ShouldNot(MatchJSON(`1e401`))
			Expect(`1e999999999`).ShouldNot(MatchJSON(`1e999999998`))

			Expect(`12345678901234567890`).Should(MatchJSON(`12345678901234567890`))
			Expect(`{"id": 9007199254740993}`).Should(MatchJSON(`{"id": 9007199254740993}`))
			Expect(`1e400`).Should(MatchJSON(`1e400`))
			Expect(`-1e999999999`).Should(MatchJSON(`-10e999999998`))
			Expect(`1e99999999999999999999999`).Should(MatchJSON(`0.1e100000000000000000000000`))
		})

		It("compares numbers by value, not by how they are written", func() {
			Expect(`1`).Should(MatchJSON(`1.0`))
			Expect(`100`).Should(MatchJSON(`1e2`))
			Expect(`1E2`).Should(MatchJSON(`100`))
			Expect(`100`).Should(MatchJSON(`1E+2`))
			Expect(`0.01`).Should(MatchJSON(`1e-2`))
			Expect(`1.5`).Should(MatchJSON(`15e-1`))
			Expect(`-1.50`).Should(MatchJSON(`-15E-1`))
			Expect(`0`).Should(MatchJSON(`-0`))
			Expect(`0`).Should(MatchJSON(`0.0`))
			Expect(`-0`).Should(MatchJSON(`0.0`))
			Expect(`0e10`).Should(MatchJSON(`-0.0e-10`))
			Expect(`{"a": 1}`).Should(MatchJSON(`{"a": 1.0}`))
			Expect(`[1, {"b": [100, 0.5]}]`).Should(MatchJSON(`[1.0, {"b": [1e2, 5e-1]}]`))

			Expect(`1`).ShouldNot(MatchJSON(`-1`))
			Expect(`1`).ShouldNot(MatchJSON(`10`))
			Expect(`1.5`).ShouldNot(MatchJSON(`15`))
			Expect(`100`).ShouldNot(MatchJSON(`1e3`))
			Expect(`1`).ShouldNot(MatchJSON(`"1"`))
			Expect(`[1, {"b": [100, 0.5]}]`).ShouldNot(MatchJSON(`[1.0, {"b": [1e2, 5e-2]}]`))
		})

		It("reports the path to the first mismatched number", func() {
			subject := MatchJSONMatcher{JSONToMatch: `{"a": [1, {"b": 9007199254740992}]}`}
			actual := `{"a": [1.0, {"b": 9007199254740993}]}`
			Expect(subject.Match(actual)).Should(BeFalse())
			Expect(subject.FailureMessage(actual)).Should(ContainSubstring(`first mismatched key: "a"[1]."b"`))
			Expect(subject.FailureMessage(actual)).Should(ContainSubstring(`9007199254740993`))
		})
	})

	When("the expected is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: 2}).Match("{}")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <int>: 2"))

			success, err = (&MatchJSONMatcher{JSONToMatch: nil}).Match("{}")
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <nil>: nil"))
		})
	})

	When("the actual is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchJSONMatcher{JSONToMatch: "{}"}).Match(2)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <int>: 2"))

			success, err = (&MatchJSONMatcher{JSONToMatch: "{}"}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <nil>: nil"))
		})
	})

	It("shows negated failure message", func() {
		failuresMessages := InterceptGomegaFailures(func() {
			Expect("1").ToNot(MatchJSON("1"))
		})
		Expect(failuresMessages).To(Equal([]string{"Expected\n    <string>: 1\nnot to match JSON of\n    <string>: 1"}))
	})
})
