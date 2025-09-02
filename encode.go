package json

func encodeString(s string) string {
	encoded := []rune{'"'}
	for _, char := range s {
		if char == '"' || char == '\\' {
			encoded = append(encoded, '\\', char)
		} else if char == '\b' {
			encoded = append(encoded, '\\', 'b')
		} else if char == '\f' {
			encoded = append(encoded, '\\', 'f')
		} else if char == '\n' {
			encoded = append(encoded, '\\', 'n')
		} else if char == '\r' {
			encoded = append(encoded, '\\', 'r')
		} else if char == '\t' {
			encoded = append(encoded, '\\', 't')
		} else if char >= 0x20 && char <= 0x10ffff {
			encoded = append(encoded, char)
		}
	}
	encoded = append(encoded, '"')
	return string(encoded)
}
