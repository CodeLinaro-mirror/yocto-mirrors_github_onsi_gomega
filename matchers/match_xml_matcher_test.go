package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("MatchXMLMatcher", func() {

	var (
		sample_01 = readFileContents("test_data/xml/sample_01.xml")
		sample_02 = readFileContents("test_data/xml/sample_02.xml")
		sample_03 = readFileContents("test_data/xml/sample_03.xml")
		sample_04 = readFileContents("test_data/xml/sample_04.xml")
		sample_05 = readFileContents("test_data/xml/sample_05.xml")
		sample_06 = readFileContents("test_data/xml/sample_06.xml")
		sample_07 = readFileContents("test_data/xml/sample_07.xml")
		sample_08 = readFileContents("test_data/xml/sample_08.xml")
		sample_09 = readFileContents("test_data/xml/sample_09.xml")
		sample_10 = readFileContents("test_data/xml/sample_10.xml")
		sample_11 = readFileContents("test_data/xml/sample_11.xml")
	)

	When("passed stringifiables", func() {
		It("matches documents regardless of the attribute order", func() {
			a := `<a foo="bar" ka="boom"></a>`
			b := `<a ka="boom" foo="bar"></a>`
			Expect(b).Should(MatchXML(a))
			Expect(a).Should(MatchXML(b))
		})

		It("should succeed if the XML matches", func() {
			Expect(sample_01).Should(MatchXML(sample_01))    // same XML
			Expect(sample_01).Should(MatchXML(sample_02))    // same XML with blank lines
			Expect(sample_01).Should(MatchXML(sample_03))    // same XML with different formatting
			Expect(sample_01).ShouldNot(MatchXML(sample_04)) // same structures with different values
			Expect(sample_01).ShouldNot(MatchXML(sample_05)) // different structures
			Expect(sample_06).ShouldNot(MatchXML(sample_07)) // same xml names with different namespaces
			Expect(sample_07).ShouldNot(MatchXML(sample_08)) // same structures with different values
			Expect(sample_09).ShouldNot(MatchXML(sample_10)) // same structures with different attribute values
			Expect(sample_11).Should(MatchXML(sample_11))    // with non UTF-8 encoding
		})
	})

	Context("with XML namespaces (https://github.com/onsi/gomega/issues/932)", func() {
		It("matches documents that bind the same namespace to different prefixes", func() {
			Expect(`<a xmlns:p="urn:u"><p:b/></a>`).Should(MatchXML(`<a xmlns:q="urn:u"><q:b/></a>`))
			Expect(`<a xmlns="urn:u"><b/></a>`).Should(MatchXML(`<p:a xmlns:p="urn:u"><p:b/></p:a>`))
			Expect(`<a xmlns:p="urn:u"><p:b/></a>`).Should(MatchXML(`<a><b xmlns="urn:u"/></a>`))
			Expect(`<a xmlns:p="urn:u" p:x="1"/>`).Should(MatchXML(`<a xmlns:q="urn:u" q:x="1"/>`))
		})

		It("does not match documents whose elements are in different namespaces", func() {
			Expect(`<a xmlns:p="urn:u"><p:b/></a>`).ShouldNot(MatchXML(`<a xmlns:p="urn:v"><p:b/></a>`))
			Expect(`<a xmlns="urn:u"><b/></a>`).ShouldNot(MatchXML(`<a><b/></a>`))
			Expect(`<a xmlns:p="urn:u"><p:b/></a>`).ShouldNot(MatchXML(`<a><b/></a>`))
			Expect(`<a xmlns:p="urn:u"><p:b/></a>`).ShouldNot(MatchXML(`<a xmlns:p="urn:u"><b/></a>`))
		})

		It("compares attributes by namespace, not by prefix", func() {
			Expect(`<a xmlns:p="urn:u" p:x="1"/>`).ShouldNot(MatchXML(`<a xmlns:p="urn:v" p:x="1"/>`))
			Expect(`<a xmlns:p="urn:u" p:x="1"/>`).ShouldNot(MatchXML(`<a xmlns:p="urn:u" x="1"/>`))
			Expect(`<a xmlns:p="urn:u" p:x="1"/>`).ShouldNot(MatchXML(`<a xmlns:p="urn:u" p:x="2"/>`))
		})

		It("matches documents regardless of the order of attributes that share a local name", func() {
			a := `<a xmlns:p="urn:u" xmlns:q="urn:v" p:id="1" q:id="2" id="3"/>`
			b := `<a id="3" q:id="2" xmlns:q="urn:v" p:id="1" xmlns:p="urn:u"/>`
			Expect(b).Should(MatchXML(a))
			Expect(a).Should(MatchXML(b))
			Expect(`<a xmlns:p="urn:u" xmlns:q="urn:v" p:id="2" q:id="1"/>`).ShouldNot(MatchXML(a))
		})

		It("shows the documents as written in the failure message", func() {
			failuresMessages := InterceptGomegaFailures(func() {
				Expect(`<a xmlns:p="urn:u"><p:b/></a>`).To(MatchXML(`<a xmlns:p="urn:v"><p:b/></a>`))
			})
			Expect(failuresMessages).To(Equal([]string{"Expected\n<a xmlns:p=\"urn:u\"><p:b/></a>\nto match XML of\n<a xmlns:p=\"urn:v\"><p:b/></a>"}))
		})
	})

	When("the expected is not valid XML", func() {
		It("should error and explain why", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: sample_01}).Match(`oops`)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Actual 'oops' should be valid XML"))
		})
	})

	When("the actual is not valid XML", func() {
		It("should error and explain why", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: `oops`}).Match(sample_01)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Expected 'oops' should be valid XML"))
		})
	})

	When("the expected is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: 2}).Match(sample_01)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <int>: 2"))

			success, err = (&MatchXMLMatcher{XMLToMatch: nil}).Match(sample_01)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got expected:\n    <nil>: nil"))
		})
	})

	When("the actual is neither a string nor a stringer nor a byte array", func() {
		It("should error", func() {
			success, err := (&MatchXMLMatcher{XMLToMatch: sample_01}).Match(2)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <int>: 2"))

			success, err = (&MatchXMLMatcher{XMLToMatch: sample_01}).Match(nil)
			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("MatchXMLMatcher matcher requires a string, stringer, or []byte.  Got actual:\n    <nil>: nil"))
		})
	})

	It("shows failure message", func() {
		failuresMessages := InterceptGomegaFailures(func() {
			Expect("<xml/>").To(MatchXML("<yml/>"))
		})
		Expect(failuresMessages).To(Equal([]string{"Expected\n<xml/>\nto match XML of\n<yml/>"}))
	})

	It("shows negated failure message", func() {
		failuresMessages := InterceptGomegaFailures(func() {
			Expect("<xml/>").ToNot(MatchXML("<xml/>"))
		})
		Expect(failuresMessages).To(Equal([]string{"Expected\n<xml/>\nnot to match XML of\n<xml/>"}))
	})
})
