package json

import (
	"fmt"
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
func (objectBuilder *ObjectBuilderStruct) AddJSON(name string, valueJSON string) {
	if objectBuilder.b.Len() == 0 {
		objectBuilder.b.WriteRune('{')
	}
	if objectBuilder.memberCount > 0 {
		objectBuilder.b.WriteRune(',')
	}
	encodedName := encodeString(StringType(name), objectBuilder.stringCharacterEscapingBehavior)
	objectBuilder.b.WriteString(encodedName)
	objectBuilder.b.WriteRune(':')
	objectBuilder.b.WriteString(valueJSON)
	objectBuilder.memberCount++
}

// Encodes the name and value to JSON and adds a new object member.
// Takes one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
// Returns an error if the value is not a JSON type defined by the package.
// Succeeds even if a member with the same name already exists.
func (objectBuilder *ObjectBuilderStruct) Add(name string, value any) error {
	encoded, err := Encode(value, objectBuilder.stringCharacterEscapingBehavior)
	if err != nil {
		return fmt.Errorf("failed to encode value: %s", err)
	}
	objectBuilder.AddJSON(name, encoded)
	return nil
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
	memberCount                     int
}

func NewArrayBuilder(stringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface) *ArrayBuilderStruct {
	arrayBuilder := &ArrayBuilderStruct{b: &strings.Builder{}, stringCharacterEscapingBehavior: stringCharacterEscapingBehavior}
	return arrayBuilder
}

// Adds a new object member with the value untouched.
// The value is assumed to be valid JSON.
// Succeeds even if a member with the same name already exists.
//
// Control characters not allowed in JSON strings are ignored when encoding values to JSON strings.
func (arrayBuilder *ArrayBuilderStruct) AddJSON(valueJSON string) {
	if arrayBuilder.b.Len() == 0 {
		arrayBuilder.b.WriteRune('[')
	}
	if arrayBuilder.memberCount > 0 {
		arrayBuilder.b.WriteRune(',')
	}
	arrayBuilder.b.WriteString(valueJSON)
	arrayBuilder.memberCount++
}

// Encodes the value to JSON and adds a new object member.
// Takes one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
// Returns an error if the value is not a JSON type defined by the package.
// Succeeds even if a member with the same name already exists.
func (arrayBuilder *ArrayBuilderStruct) Add(value any) error {
	encoded, err := Encode(value, arrayBuilder.stringCharacterEscapingBehavior)
	if err != nil {
		return fmt.Errorf("failed to encode value: %s", err)
	}
	arrayBuilder.AddJSON(encoded)
	return nil
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
