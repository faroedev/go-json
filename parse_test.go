package json

import (
	"bytes"
	"testing"
)

func TestParseString(t *testing.T) {
	successCases := []SuccessTestCase{
		{`"hello"`, "hello"},
		{`"say \"hi\"!"`, `say "hi"!`},
		{`"\\\/\"\b\f\n\r\t"`, "\\/\"\b\f\n\r\t"},
		{`"\uD834\uDD1E"`, "𝄞"},
		{`"\u005C"`, "\\"},
		{`"\u005c"`, "\\"},
	}
	for _, c := range successCases {
		got, err := parseString(bytes.NewReader([]byte(c.input)))
		if err != nil {
			t.Errorf("error on input: %s: %s", c.input, err)
			continue
		}
		if got != c.expected {
			t.Errorf("unexpected output on input %s: %s", c.input, got)
			continue
		}
	}

	failCases := []string{
		`"hello\x"`,
		`"\uD834"`,
		`"\uD834t"`,
		`"\uD834\""`,
	}
	for _, c := range failCases {
		_, err := parseString(bytes.NewReader([]byte(c)))
		if err == nil {
			t.Errorf("expected error on input: %s", c)
			continue
		}
	}
}

type SuccessTestCase struct {
	input    string
	expected string
}
