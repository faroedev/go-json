package json

import (
	"strings"
	"unicode/utf16"
)

// UseCharacter reports whether a character must be written directly in a JSON string without escaping.
// The parameter r must be a character that does not require escaping.
//
// UseShorthandEscapeSequence reports whether a character that is required to be escaped (defined either by the JSON specification or UseCharacter)
// must be escaped using its shorthand form.
// The parameter r must be a character that has a shorthand escape sequence.
//
// UseCharacter should be used before UseShorthandEscapeSequence.
type StringCharacterEscapingBehaviorInterface interface {
	UseCharacter(r rune) bool
	UseShorthandEscapeSequence(r rune) bool
}

// Escapes characters only when necessary, preferring shorthand escape sequences.
var MinimalStringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface = minimalStringCharacterEscapingBehaviorStruct{}

type minimalStringCharacterEscapingBehaviorStruct struct{}

func (minimalStringCharacterEscapingBehaviorStruct) UseCharacter(_ rune) bool {
	return true
}

func (minimalStringCharacterEscapingBehaviorStruct) UseShorthandEscapeSequence(_ rune) bool {
	return true
}

var shorthandStringCharacterEscapeSequences = map[rune]string{
	'"':  `\"`,
	'\\': `\\`,
	'/':  `\/`,
	'\b': `\b`,
	'\f': `\f`,
	'\n': `\n`,
	'\r': `\r`,
	'\t': `\t`,
}

func encodeString(s string, characterEscapingBehavior StringCharacterEscapingBehaviorInterface) string {
	b := strings.Builder{}
	b.WriteRune('"')
	for _, char := range s {
		if char >= 0x0000 && char <= 0x001f && char != '"' && char != '\\' {
			if characterEscapingBehavior.UseCharacter(char) {
				b.WriteRune(char)
				continue
			}
		}

		if shorthandEscapeSequence, ok := shorthandStringCharacterEscapeSequences[char]; ok {
			if characterEscapingBehavior.UseShorthandEscapeSequence(char) {
				b.WriteString(shorthandEscapeSequence)
				continue
			}
		}

		if utf16.RuneLen(char) > 1 {
			pair1, pair2 := utf16.EncodeRune(char)
			b.WriteString(toStringHexEscapeSequence(pair1))
			b.WriteString(toStringHexEscapeSequence(pair2))
		} else {
			b.WriteString(toStringHexEscapeSequence(char))
		}
	}
	b.WriteRune('"')
	return b.String()
}

const hexTable = "0123456789abcdef"

func toStringHexEscapeSequence(r rune) string {
	b := make([]byte, 6)
	b[0] = byte('\\')
	b[1] = byte('u')
	b[2] = hexTable[(r>>12)&0x0f]
	b[3] = hexTable[(r>>8)&0x0f]
	b[4] = hexTable[(r>>4)&0x0f]
	b[5] = hexTable[r&0x0f]
	return string(b)
}
