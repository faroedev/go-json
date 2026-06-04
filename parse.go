package json

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf16"
)

// Parses a JSON-encoded string. Ignores any leading and trailing whitespace.
// Returns one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
// Returns an error if the string is an invalid JSON object or
// an object has duplicate member names.
//
// JSON object member names are compared after resolving any escaped characters.
func Parse(s string) (any, error) {
	r := strings.NewReader(s)

	jsonValue, err := parseEmbeddedValue(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value: %s", err.Error())
	}

	err = parseEnd(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end: %s", err.Error())
	}

	return jsonValue, nil
}

// Parses a JSON-encoded string as an object with [Parse].
func ParseObject(s string) (ObjectType, error) {
	jsonValue, err := Parse(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %s", err.Error())
	}

	jsonObject, ok := jsonValue.(ObjectType)
	if !ok {
		return nil, fmt.Errorf("value not an object")
	}
	return jsonObject, nil
}

// Parses a JSON-encoded string as an array with [Parse].
func ParseArray(s string) (ArrayType, error) {
	jsonValue, err := Parse(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %s", err.Error())
	}

	jsonArray, ok := jsonValue.(ArrayType)
	if !ok {
		return nil, fmt.Errorf("value not an array")
	}
	return jsonArray, nil
}

func parseEnd(r io.RuneScanner) error {
	for {
		char, _, err := r.ReadRune()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return fmt.Errorf("invalid encoding")
		}
		if char == '	' || char == '\n' || char == ' ' || char == '\r' {
			continue
		}
		return errors.New("unexpected character")
	}
	return nil
}

func parseEmbeddedValue(r io.RuneScanner) (any, error) {
	err := skipWhitespace(r)
	if err != nil {
		return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
	}

	nextChar, _, err := r.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if nextChar == unicode.ReplacementChar {
		return nil, fmt.Errorf("invalid encoding")
	}
	err = r.UnreadRune()
	if err != nil {
		return nil, fmt.Errorf("failed to unread rune: %s", err.Error())
	}

	if nextChar == '{' {
		jsonObject, err := parseEmbeddedObject(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse object: %s", err.Error())
		}
		return jsonObject, nil
	}
	if nextChar == '[' {
		jsonArray, err := parseEmbeddedArray(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse array: %s", err.Error())
		}
		return jsonArray, nil
	}
	if nextChar == '"' {
		jsonString, err := parseString(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string: %s", err.Error())
		}
		return jsonString, nil
	}
	if isNumberDigitCharacter(nextChar) {
		jsonNumber, err := parseNumber(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %s", err.Error())
		}
		return jsonNumber, nil
	}

	literalName, err := parseLiteralName(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse literal name: %s", err.Error())
	}

	switch literalName {
	case "true":
		return True, nil
	case "false":
		return False, nil
	case "null":
		return Null, nil
	default:
		return "", fmt.Errorf("invalid literal name: %s", literalName)
	}
}

func parseEmbeddedObject(r io.RuneScanner) (ObjectType, error) {
	object := NewObject()

	err := skipWhitespace(r)
	if err != nil {
		return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
	}

	char, _, err := r.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return nil, fmt.Errorf("invalid encoding")
	}
	if char != '{' {
		return nil, fmt.Errorf("unexpected character %s", string(char))
	}

	for {
		err := skipWhitespace(r)
		if err != nil {
			return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err := r.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return nil, fmt.Errorf("invalid encoding")
		}
		if char == '}' {
			break
		}
		err = r.UnreadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to unread rune: %s", err.Error())
		}

		jsonMemberString, err := parseString(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse member name: %s", err.Error())
		}
		memberName := string(jsonMemberString)
		if object.Has(memberName) {
			return nil, fmt.Errorf("duplicate member name %s", memberName)
		}

		err = skipWhitespace(r)
		if err != nil {
			return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return nil, fmt.Errorf("invalid encoding")
		}
		if char != ':' {
			return nil, fmt.Errorf("unexpected character %s", string(char))
		}

		jsonMemberValue, err := parseEmbeddedValue(r)
		object.Set(memberName, jsonMemberValue)

		err = skipWhitespace(r)
		if err != nil {
			return nil, fmt.Errorf("Failed to skip whitespace: %s", err.Error())
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return nil, fmt.Errorf("invalid encoding")
		}
		if char == '}' {
			break
		}
		if char != ',' {
			return nil, fmt.Errorf("unexpected character %s", string(char))
		}
	}

	return object, nil
}

func parseEmbeddedArray(r io.RuneScanner) (ArrayType, error) {
	array := NewArray()

	err := skipWhitespace(r)
	if err != nil {
		return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
	}

	char, _, err := r.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return nil, fmt.Errorf("invalid encoding")
	}
	if char != '[' {
		return nil, fmt.Errorf("unexpected character %s", string(char))
	}

	for {
		err := skipWhitespace(r)
		if err != nil {
			return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err := r.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return nil, fmt.Errorf("invalid encoding")
		}
		if char == ']' {
			break
		}
		err = r.UnreadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to unread rune: %s", err.Error())
		}

		jsonMemberValue, err := parseEmbeddedValue(r)
		array.Add(jsonMemberValue)

		err = skipWhitespace(r)
		if err != nil {
			return nil, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return nil, fmt.Errorf("invalid encoding")
		}
		if char == ']' {
			break
		}
		if char != ',' {
			return nil, fmt.Errorf("unexpected character %s", string(char))
		}
	}

	return array, nil
}

func parseString(r io.RuneScanner) (StringType, error) {
	b := strings.Builder{}

	char, _, err := r.ReadRune()
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if char != '"' {
		return "", fmt.Errorf("unexpected character %s", string(char))
	}

	var prevHex rune = 0
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			return "", fmt.Errorf("failed to read rune: %s", err.Error())
		}

		if char == '"' {
			if prevHex > 0 {
				return "", fmt.Errorf("unexpected character %s", string(char))
			}
			break
		}

		if char == '\\' {
			char, _, err := r.ReadRune()
			if err != nil {
				return "", fmt.Errorf("failed to read rune: %s", err.Error())
			}
			if char == 'u' {
				var decoded rune = 0
				for i := range 4 {
					char, _, err := r.ReadRune()
					if err != nil {
						return "", fmt.Errorf("failed to read rune: %s", err.Error())
					}

					var b rune
					if char >= '0' && char <= '9' {
						b = (char) - '0'
					} else if char >= 'A' && char <= 'F' {
						b = (char) - 'A' + 10
					} else if char >= 'a' && char <= 'f' {
						b = (char) - 'a' + 10
					} else {
						return "", fmt.Errorf("invalid hex encoding")
					}
					decoded |= b << ((3 - i) * 4)
				}
				if prevHex > 0 {
					decoded = utf16.DecodeRune(prevHex, decoded)
					if decoded == unicode.ReplacementChar {
						return "", fmt.Errorf("invalid character encoding")
					}
					b.WriteRune(decoded)
					prevHex = 0
				} else if utf16.IsSurrogate(decoded) {
					prevHex = decoded
				} else {
					b.WriteRune(decoded)
				}
				continue
			}
			if prevHex > 0 {
				return "", errors.New("expected hex encoding")
			}
			switch char {
			case '"', '\\', '/':
				b.WriteRune(char)
			case 'b':
				b.WriteRune('\b')
			case 'f':
				b.WriteRune('\f')
			case 'n':
				b.WriteRune('\n')
			case 'r':
				b.WriteRune('\r')
			case 't':
				b.WriteRune('\t')
			default:
				return "", fmt.Errorf("unexpected escape character %s", string(char))
			}
			continue
		}

		if prevHex > 0 {
			return "", errors.New("expected hex encoding")
		}

		if char < 0x20 || char > 0x10ffff {
			return "", fmt.Errorf("invalid character")
		}

		b.WriteRune(char)

	}

	return StringType(b.String()), nil
}

func parseNumber(r io.RuneScanner) (NumberType, error) {
	extracted := []rune{}
	char, _, err := r.ReadRune()
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if char == '-' {
		extracted = append(extracted, char)
	} else {
		err = r.UnreadRune()
		if err != nil {
			return "", fmt.Errorf("failed to unread rune: %s", err.Error())
		}
	}

	char, _, err = r.ReadRune()
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if char == '0' {
		extracted = append(extracted, char)
	} else if char >= '1' && char <= '9' {
		extracted = append(extracted, char)
		for {
			char, _, err = r.ReadRune()
			if err != nil {
				return "", fmt.Errorf("failed to read rune: %s", err.Error())
			}
			if char == unicode.ReplacementChar {
				return "", fmt.Errorf("invalid character encoding")
			}
			if !isNumberDigitCharacter(char) {
				err = r.UnreadRune()
				if err != nil {
					return "", fmt.Errorf("failed to unread rune: %s", err.Error())
				}
				break
			}
			extracted = append(extracted, char)
		}
	} else {
		return "", fmt.Errorf("unexpected character %s", string(char))
	}

	char, _, err = r.ReadRune()
	if err != nil && errors.Is(err, io.EOF) {
		return NumberType(string(extracted)), nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if char == '.' {
		extracted = append(extracted, char)
		for {
			char, _, err = r.ReadRune()
			if err != nil {
				return "", fmt.Errorf("failed to read rune: %s", err.Error())
			}
			if char == unicode.ReplacementChar {
				return "", fmt.Errorf("invalid encoding")
			}
			if !isNumberDigitCharacter(char) {
				err = r.UnreadRune()
				if err != nil {
					return "", fmt.Errorf("failed to unread rune: %s", err.Error())
				}
				break
			}
			extracted = append(extracted, char)
		}
	} else {
		err = r.UnreadRune()
		if err != nil {
			return "", fmt.Errorf("failed to unread rune: %s", err.Error())
		}
	}

	char, _, err = r.ReadRune()
	if err != nil && errors.Is(err, io.EOF) {
		return NumberType(string(extracted)), nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if char == 'E' || char == 'e' {
		extracted = append(extracted, char)

		char, _, err = r.ReadRune()
		if err != nil {
			return "", fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return "", fmt.Errorf("invalid encoding")
		}
		if char == '-' || char == '+' {
			extracted = append(extracted, char)
		} else {
			err = r.UnreadRune()
			if err != nil {
				return "", fmt.Errorf("failed to unread rune: %s", err.Error())
			}
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return "", fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return "", fmt.Errorf("invalid encoding")
		}
		if !isNumberDigitCharacter(char) {
			return "", fmt.Errorf("unexpected character %s", string(char))
		}
		extracted = append(extracted, char)

		for {
			char, _, err = r.ReadRune()
			if err != nil && errors.Is(err, io.EOF) {
				return NumberType(string(extracted)), nil
			}
			if err != nil {
				return "", fmt.Errorf("failed to read rune: %s", err.Error())
			}
			if char == unicode.ReplacementChar {
				return "", fmt.Errorf("invalid encoding")
			}
			if !isNumberDigitCharacter(char) {
				err = r.UnreadRune()
				if err != nil {
					return "", fmt.Errorf("failed to unread rune: %s", err.Error())
				}
				break
			}
			extracted = append(extracted, char)
		}
	} else {
		err = r.UnreadRune()
		if err != nil {
			return "", fmt.Errorf("failed to unread rune: %s", err.Error())
		}
	}

	return NumberType(string(extracted)), nil
}

func parseLiteralName(r io.RuneScanner) (string, error) {
	literalNameCharacters := []rune{}
	char, _, err := r.ReadRune()
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if !isLiteralNameCharacter(char) {
		return "", fmt.Errorf("unexpected character %s", string(char))
	}
	literalNameCharacters = append(literalNameCharacters, char)

	for {
		char, _, err := r.ReadRune()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return "", fmt.Errorf("invalid encoding")
		}
		if !isLiteralNameCharacter(char) {
			err = r.UnreadRune()
			if err != nil {
				return "", fmt.Errorf("failed to unread rune: %s", err.Error())
			}
			break
		}
		literalNameCharacters = append(literalNameCharacters, char)
	}

	literalName := string(literalNameCharacters)
	return literalName, nil
}

func skipWhitespace(r io.RuneScanner) error {
	for {
		char, _, err := r.ReadRune()
		if err != nil && errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return fmt.Errorf("invalid encoding")
		}
		if char == '	' || char == '\n' || char == ' ' || char == '\r' {
			continue
		}
		err = r.UnreadRune()
		if err != nil {
			return fmt.Errorf("failed to unread rune: %s", err.Error())
		}
		return nil
	}
}

func isLiteralNameCharacter(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

func isNumberDigitCharacter(r rune) bool {
	return r >= '0' && r <= '9'
}
