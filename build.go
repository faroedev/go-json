package json

import (
	"strconv"
	"strings"
)

// Use [NewObjectBuilder].
type ObjectBuilderStruct struct {
	b                               *strings.Builder
	stringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface
	memberCount                     int
}

func NewObjectBuilder(stringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface) *ObjectBuilderStruct {
	objectBuilder := &ObjectBuilderStruct{b: &strings.Builder{}, stringCharacterEscapingBehavior: stringCharacterEscapingBehavior}
	return objectBuilder
}

// Encodes the name to a JSON string and adds a new object member with the value untouched.
// The value is assumed to be valid JSON.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddJSON(name string, value string) {
	if objectBuilder.b.Len() == 0 {
		objectBuilder.b.WriteRune('{')
	}
	if objectBuilder.memberCount > 0 {
		objectBuilder.b.WriteRune(',')
	}
	encodedName := encodeString(name, objectBuilder.stringCharacterEscapingBehavior)
	objectBuilder.b.WriteString(encodedName)
	objectBuilder.b.WriteRune(':')
	objectBuilder.b.WriteString(value)
	objectBuilder.memberCount++
}

// Encodes the name and value to JSON strings, and adds a new object member.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddString(name string, value string) {
	encoded := encodeString(value, objectBuilder.stringCharacterEscapingBehavior)
	objectBuilder.AddJSON(name, encoded)
}

// Encodes the name to a JSON string and value to a JSON number, and adds a new object member.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddInt(name string, value int) {
	encoded := strconv.FormatInt(int64(value), 10)
	objectBuilder.AddJSON(name, encoded)
}

// Encodes the name to a JSON string and value to a JSON number, and adds a new object member.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddInt64(name string, value int64) {
	encoded := strconv.FormatInt(value, 10)
	objectBuilder.AddJSON(name, encoded)
}

// Encodes the name to a JSON string and value to a JSON number, and adds a new object member.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddInt32(name string, value int32) {
	objectBuilder.AddInt64(name, int64(value))
}

// Encodes the name to a JSON string and value to a JSON boolean, and adds a new object member.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddBool(name string, value bool) {
	if value {
		objectBuilder.AddJSON(name, "true")
	} else {
		objectBuilder.AddJSON(name, "false")
	}
}

// Encodes the name to a JSON string and adds a new object member with a null value.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (objectBuilder *ObjectBuilderStruct) AddNull(name string) {
	objectBuilder.AddJSON(name, `null`)
}

// Returns the built JSON.
// The builder can no longer be used.
func (objectBuilder *ObjectBuilderStruct) Done() string {
	if objectBuilder.b.Len() == 0 {
		return "{}"
	}
	objectBuilder.b.WriteRune('}')
	return objectBuilder.b.String()
}

// Use [NewArrayBuilder].
type ArrayBuilderStruct struct {
	b                               *strings.Builder
	stringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface
	elementCount                    int
}

func NewArrayBuilder(stringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface) *ArrayBuilderStruct {
	arrayBuilder := &ArrayBuilderStruct{b: &strings.Builder{}, stringCharacterEscapingBehavior: stringCharacterEscapingBehavior}
	return arrayBuilder
}

// Adds the JSON value as a new array element.
// The value is assumed to be valid JSON.
func (arrayBuilder *ArrayBuilderStruct) AddJSON(value string) {
	if arrayBuilder.b.Len() == 0 {
		arrayBuilder.b.WriteRune('[')
	}
	if arrayBuilder.elementCount > 0 {
		arrayBuilder.b.WriteRune(',')
	}
	arrayBuilder.b.WriteString(value)
	arrayBuilder.elementCount++
}

// Encodes the value to a JSON string and adds it as a new array element.
// Control characters not allowed in JSON strings are ignored when encoding.
func (arrayBuilder *ArrayBuilderStruct) AddString(value string) {
	encoded := encodeString(value, arrayBuilder.stringCharacterEscapingBehavior)
	arrayBuilder.AddJSON(encoded)
}

// Encodes the value to a JSON number and adds it as a new array element.
func (arrayBuilder *ArrayBuilderStruct) AddInt(value int) {
	encoded := strconv.FormatInt(int64(value), 10)
	arrayBuilder.AddJSON(encoded)
}

// Encodes the value to a JSON number and adds it as a new array element.
func (arrayBuilder *ArrayBuilderStruct) AddInt64(value int64) {
	encoded := strconv.FormatInt(value, 10)
	arrayBuilder.AddJSON(encoded)
}

// Encodes the value to a JSON number and adds it as a new array element.
func (arrayBuilder *ArrayBuilderStruct) AddInt32(key string, value int32) {
	encoded := strconv.FormatInt(int64(value), 10)
	arrayBuilder.AddJSON(encoded)
}

// Encodes the value to a JSON boolean and adds it as a new array element.
func (arrayBuilder *ArrayBuilderStruct) AddBool(value bool) {
	if value {
		arrayBuilder.AddJSON("true")
	} else {
		arrayBuilder.AddJSON("false")
	}
}

// Adds null to the array.
func (arrayBuilder *ArrayBuilderStruct) AddNull() {
	arrayBuilder.AddJSON("null")
}

// Returns the built JSON.
// The builder can no longer be used.
func (arrayBuilder *ArrayBuilderStruct) Done() string {
	if arrayBuilder.b.Len() == 0 {
		return "[]"
	}
	arrayBuilder.b.WriteRune(']')
	return arrayBuilder.b.String()
}
