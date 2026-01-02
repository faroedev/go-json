package json

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf16"
)

// Parses a JSON object. Ignores any leading and trailing whitespace.
// Returns an error if the string is an invalid JSON object or
// an object has duplicate member names.
//
// JSON object member names are compared after resolving any escaped characters.
func ParseObject(s string) (ObjectStruct, error) {
	r := strings.NewReader(s)

	parsed, err := parseEmbeddedObject(r)
	if err != nil {
		return ObjectStruct{}, fmt.Errorf("failed to parse embedded object: %s", err.Error())
	}

	err = parseEnd(r)
	if err != nil {
		return ObjectStruct{}, fmt.Errorf("failed to parse end: %s", err.Error())
	}

	return parsed, nil
}

// Parses a JSON array. Ignores any leading and trailing whitespace.
// Returns an error if the string is an invalid JSON array or
// an object has duplicate member names.
//
// JSON object member names are compared after resolving any escaped characters.
func ParseArray(s string) (ArrayStruct, error) {
	r := strings.NewReader(s)

	parsed, err := parseEmbeddedArray(r)
	if err != nil {
		return ArrayStruct{}, fmt.Errorf("failed to parse embedded array: %s", err.Error())
	}

	err = parseEnd(r)
	if err != nil {
		return ArrayStruct{}, fmt.Errorf("failed to parse end: %s", err.Error())
	}

	return parsed, nil
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

func parseEmbeddedObject(r io.RuneScanner) (ObjectStruct, error) {
	object := NewObject()

	err := skipWhitespace(r)
	if err != nil {
		return ObjectStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
	}

	char, _, err := r.ReadRune()
	if err != nil {
		return ObjectStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return ObjectStruct{}, fmt.Errorf("invalid encoding")
	}
	if char != '{' {
		return ObjectStruct{}, fmt.Errorf("unexpected character %s", string(char))
	}

	for {
		err := skipWhitespace(r)
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err := r.ReadRune()
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ObjectStruct{}, fmt.Errorf("invalid encoding")
		}
		if char == '}' {
			break
		}
		err = r.UnreadRune()
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to unread rune: %s", err.Error())
		}

		key, err := parseString(r)
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to parse member name: %s", err.Error())
		}
		if object.Has(key) {
			return ObjectStruct{}, fmt.Errorf("duplicate member name %s", key)
		}

		err = skipWhitespace(r)
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ObjectStruct{}, fmt.Errorf("invalid encoding")
		}
		if char != ':' {
			return ObjectStruct{}, fmt.Errorf("unexpected character %s", string(char))
		}

		err = skipWhitespace(r)
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		nextChar, _, err := r.ReadRune()
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ObjectStruct{}, fmt.Errorf("invalid encoding")
		}
		err = r.UnreadRune()
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to unread rune: %s", err.Error())
		}
		if nextChar == '{' {
			value, err := parseEmbeddedObject(r)
			if err != nil {
				return ObjectStruct{}, fmt.Errorf("failed to parse object: %s", err.Error())
			}
			object.SetJSONObject(key, value)
		} else if nextChar == '[' {
			value, err := parseEmbeddedArray(r)
			if err != nil {
				return ObjectStruct{}, fmt.Errorf("failed to parse array: %s", err.Error())
			}
			object.SetJSONArray(key, value)
		} else if nextChar == '"' {
			value, err := parseString(r)
			if err != nil {
				return ObjectStruct{}, fmt.Errorf("failed to parse string: %s", err.Error())
			}
			object.SetString(key, value)
		} else if isDigitCharacter(nextChar) {
			value, err := extractNumber(r)
			if err != nil {
				return ObjectStruct{}, fmt.Errorf("failed to extract number: %s", err.Error())
			}
			object.SetNumber(key, value)
		} else {
			value, err := extractIdentifier(r)
			if err != nil {
				return ObjectStruct{}, fmt.Errorf("failed to extract identifier: %s", err.Error())
			}
			switch value {
			case "true":
				object.SetBool(key, true)
			case "false":
				object.SetBool(key, false)
			case "null":
				object.SetNull(key)
			default:
				return ObjectStruct{}, fmt.Errorf("unexpected identifier %s", value)
			}
		}

		err = skipWhitespace(r)
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("Failed to skip whitespace: %s", err.Error())
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return ObjectStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ObjectStruct{}, fmt.Errorf("invalid encoding")
		}
		if char == '}' {
			break
		}
		if char != ',' {
			return ObjectStruct{}, fmt.Errorf("unexpected character %s", string(char))
		}
	}

	return object, nil
}

func parseEmbeddedArray(r io.RuneScanner) (ArrayStruct, error) {
	array := ArrayStruct{}

	err := skipWhitespace(r)
	if err != nil {
		return ArrayStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
	}

	char, _, err := r.ReadRune()
	if err != nil {
		return ArrayStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return ArrayStruct{}, fmt.Errorf("invalid encoding")
	}
	if char != '[' {
		return ArrayStruct{}, fmt.Errorf("unexpected character %s", string(char))
	}

	for {
		err := skipWhitespace(r)
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err := r.ReadRune()
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ArrayStruct{}, fmt.Errorf("invalid encoding")
		}
		if char == ']' {
			break
		}
		err = r.UnreadRune()
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to unread rune: %s", err.Error())
		}

		nextChar, _, err := r.ReadRune()
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ArrayStruct{}, fmt.Errorf("invalid encoding")
		}
		err = r.UnreadRune()
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to unread rune: %s", err.Error())
		}
		if nextChar == '{' {
			value, err := parseEmbeddedObject(r)
			if err != nil {
				return ArrayStruct{}, fmt.Errorf("failed to parse embedded object: %s", err.Error())
			}
			array.AddJSONObject(value)
		} else if nextChar == '[' {
			value, err := parseEmbeddedArray(r)
			if err != nil {
				return ArrayStruct{}, fmt.Errorf("failed to parse embedded array: %s", err.Error())
			}
			array.AddJSONArray(value)
		} else if nextChar == '"' {
			value, err := parseString(r)
			if err != nil {
				return ArrayStruct{}, fmt.Errorf("failed to parse string: %s", err.Error())
			}
			array.AddString(value)
		} else if isDigitCharacter(nextChar) {
			value, err := extractNumber(r)
			if err != nil {
				return ArrayStruct{}, fmt.Errorf("failed to extract number: %s", err.Error())
			}
			array.AddNumber(value)
		} else {
			value, err := extractIdentifier(r)
			if err != nil {
				return ArrayStruct{}, fmt.Errorf("failed to extract identifier: %s", err.Error())
			}

			switch value {
			case "true":
				array.AddBool(true)
			case "false":
				array.AddBool(false)
			case "null":
				array.AddNull()
			default:
				return ArrayStruct{}, fmt.Errorf("unexpected identifier %s", value)
			}
		}

		err = skipWhitespace(r)
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to skip whitespace: %s", err.Error())
		}

		char, _, err = r.ReadRune()
		if err != nil {
			return ArrayStruct{}, fmt.Errorf("failed to read rune: %s", err.Error())
		}
		if char == unicode.ReplacementChar {
			return ArrayStruct{}, fmt.Errorf("invalid encoding")
		}
		if char == ']' {
			break
		}
		if char != ',' {
			return ArrayStruct{}, fmt.Errorf("unexpected character %s", string(char))
		}
	}

	return array, nil
}

func parseString(r io.RuneScanner) (string, error) {
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

	return b.String(), nil
}

func extractNumber(r io.RuneScanner) (string, error) {
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
			if !isDigitCharacter(char) {
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
		return string(extracted), nil
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
			if !isDigitCharacter(char) {
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
		return string(extracted), nil
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
		if !isDigitCharacter(char) {
			return "", fmt.Errorf("unexpected character %s", string(char))
		}
		extracted = append(extracted, char)

		for {
			char, _, err = r.ReadRune()
			if err != nil && errors.Is(err, io.EOF) {
				return string(extracted), nil
			}
			if err != nil {
				return "", fmt.Errorf("failed to read rune: %s", err.Error())
			}
			if char == unicode.ReplacementChar {
				return "", fmt.Errorf("invalid encoding")
			}
			if !isDigitCharacter(char) {
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

	return string(extracted), nil
}

func extractIdentifier(r io.RuneScanner) (string, error) {
	extracted := []rune{}
	char, _, err := r.ReadRune()
	if err != nil {
		return "", fmt.Errorf("failed to read rune: %s", err.Error())
	}
	if char == unicode.ReplacementChar {
		return "", fmt.Errorf("invalid encoding")
	}
	if !isIdentifierCharacter(char) {
		return "", fmt.Errorf("unexpected character %s", string(char))
	}
	extracted = append(extracted, char)

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
		if !isIdentifierCharacter(char) {
			err = r.UnreadRune()
			if err != nil {
				return "", fmt.Errorf("failed to unread rune: %s", err.Error())
			}
			break
		}
		extracted = append(extracted, char)
	}

	return string(extracted), nil
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

func isIdentifierCharacter(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

func isDigitCharacter(r rune) bool {
	return r >= '0' && r <= '9'
}
