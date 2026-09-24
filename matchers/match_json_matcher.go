package matchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/onsi/gomega/format"
)

type MatchJSONMatcher struct {
	JSONToMatch      any
	firstFailurePath []any
}

func (matcher *MatchJSONMatcher) Match(actual any) (success bool, err error) {
	actualString, expectedString, err := matcher.prettyPrint(actual)
	if err != nil {
		return false, err
	}

	var aval any
	var eval any

	// prettyPrint has already checked the syntax, so decoding is not expected to fail
	if aval, err = decodeJSON(actualString); err != nil {
		return false, fmt.Errorf("Actual '%s' should be valid JSON, but it is not.\nUnderlying error:%s", actualString, err)
	}
	if eval, err = decodeJSON(expectedString); err != nil {
		return false, fmt.Errorf("Expected '%s' should be valid JSON, but it is not.\nUnderlying error:%s", expectedString, err)
	}
	var equal bool
	equal, matcher.firstFailurePath = deepEqual(aval, eval)
	return equal, nil
}

// decodeJSON decodes s with every number in canonical form, so that numbers
// are compared exactly (a float64 would lose precision) and by value (so that
// 1, 1.0 and 1e0 are all equal).
func decodeJSON(s string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return canonicalizeJSONNumbers(value), nil
}

// canonicalJSONNumber is a JSON number written as its significant digits
// followed by a base-10 exponent, e.g. -15e-1 for -1.50 or -15E-1. Two JSON
// numbers have the same value exactly when they have the same canonical form.
type canonicalJSONNumber string

func canonicalizeJSONNumbers(value any) any {
	switch v := value.(type) {
	case []any:
		for i, element := range v {
			v[i] = canonicalizeJSONNumbers(element)
		}
	case map[string]any:
		for key, element := range v {
			v[key] = canonicalizeJSONNumbers(element)
		}
	case json.Number:
		return canonicalizeJSONNumber(v)
	}
	return value
}

func canonicalizeJSONNumber(n json.Number) canonicalJSONNumber {
	s, sign := string(n), ""
	if rest, negative := strings.CutPrefix(s, "-"); negative {
		s, sign = rest, "-"
	}
	mantissa, exponentString, _ := strings.Cut(strings.ToLower(s), "e")
	integerPart, fractionPart, _ := strings.Cut(mantissa, ".")

	digits := strings.TrimLeft(integerPart+fractionPart, "0")
	if digits == "" {
		return "0" // so that -0 and 0 are equal
	}
	significantDigits := strings.TrimRight(digits, "0")

	// the exponent is parsed as a big.Int as JSON puts no limit on its size
	exponent := new(big.Int)
	if exponentString != "" {
		exponent.SetString(exponentString, 10) // the syntax has been checked by the decoder
	}
	exponent.Add(exponent, big.NewInt(int64(len(digits)-len(significantDigits)-len(fractionPart))))

	return canonicalJSONNumber(sign + significantDigits + "e" + exponent.String())
}

func (matcher *MatchJSONMatcher) FailureMessage(actual any) (message string) {
	actualString, expectedString, _ := matcher.prettyPrint(actual)
	return formattedMessage(format.Message(actualString, "to match JSON of", expectedString), matcher.firstFailurePath)
}

func (matcher *MatchJSONMatcher) NegatedFailureMessage(actual any) (message string) {
	actualString, expectedString, _ := matcher.prettyPrint(actual)
	return formattedMessage(format.Message(actualString, "not to match JSON of", expectedString), matcher.firstFailurePath)
}

func (matcher *MatchJSONMatcher) prettyPrint(actual any) (actualFormatted, expectedFormatted string, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return "", "", fmt.Errorf("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok := toString(matcher.JSONToMatch)
	if !ok {
		return "", "", fmt.Errorf("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.JSONToMatch, 1))
	}

	abuf := new(bytes.Buffer)
	ebuf := new(bytes.Buffer)

	if err := json.Indent(abuf, []byte(actualString), "", "  "); err != nil {
		return "", "", fmt.Errorf("Actual '%s' should be valid JSON, but it is not.\nUnderlying error:%s", actualString, err)
	}

	if err := json.Indent(ebuf, []byte(expectedString), "", "  "); err != nil {
		return "", "", fmt.Errorf("Expected '%s' should be valid JSON, but it is not.\nUnderlying error:%s", expectedString, err)
	}

	return abuf.String(), ebuf.String(), nil
}
