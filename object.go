package json

import (
	"errors"
	"fmt"
)

// Represents a JSON object.
type ObjectType map[string]any

func NewObject() ObjectType {
	object := ObjectType{}
	return object
}

func (object ObjectType) Has(key string) bool {
	_, ok := object[key]
	return ok
}

// Takes one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
func (object ObjectType) Set(key string, value any) {
	object[key] = value
}

// Returns [ErrObjectMemberNotFound] if the member doesn't.
// Returns one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
func (object ObjectType) Get(key string) (any, error) {
	value, ok := object[key]
	if !ok {
		return "", ErrObjectMemberNotFound
	}
	return value, nil
}

func (object ObjectType) Delete(key string) {
	delete(object, key)
}

// Returns [ErrObjectMemberNotFound] if the member doesn't exist, or an error if the value isn't a string.
func (object ObjectType) GetString(key string) (StringType, error) {
	value, ok := object[key]
	if !ok {
		return "", ErrObjectMemberNotFound
	}
	jsonString, ok := value.(StringType)
	if !ok {
		return "", fmt.Errorf("value not a string")
	}
	return jsonString, nil
}

// Returns [ErrObjectMemberNotFound] if the member doesn't exist, or an error if the value isn't a number.
func (object ObjectType) GetNumber(key string) (NumberType, error) {
	value, ok := object[key]
	if !ok {
		return "", ErrObjectMemberNotFound
	}
	jsonNumber, ok := value.(NumberType)
	if !ok {
		return "", fmt.Errorf("value not a number")
	}
	return jsonNumber, nil
}

// Returns [ErrObjectMemberNotFound] if the member doesn't exist, or an error if the value isn't a boolean.
func (object ObjectType) GetBoolean(key string) (BooleanType, error) {
	value, ok := object[key]
	if !ok {
		return false, ErrObjectMemberNotFound
	}
	jsonBoolean, ok := value.(BooleanType)
	if !ok {
		return false, fmt.Errorf("value not a boolean")
	}
	return jsonBoolean, nil
}

// Returns [ErrObjectMemberNotFound] if the member doesn't exist, or an error if the value isn't an object.
func (object ObjectType) GetObject(key string) (ObjectType, error) {
	value, ok := object[key]
	if !ok {
		return nil, ErrObjectMemberNotFound
	}
	jsonObject, ok := value.(ObjectType)
	if !ok {
		return nil, fmt.Errorf("value not an object")
	}
	return jsonObject, nil
}

// Returns [ErrObjectMemberNotFound] if the member doesn't exist.
// Returns an error if the value isn't an array.
func (object ObjectType) GetArray(key string) (ArrayType, error) {
	value, ok := object[key]
	if !ok {
		return nil, ErrObjectMemberNotFound
	}
	jsonArray, ok := value.(ArrayType)
	if !ok {
		return nil, fmt.Errorf("value not an array")
	}
	return jsonArray, nil
}

// Returns [ErrObjectMemberNotFound] if the member doesn't exist.
func (object ObjectType) IsNull(key string) (bool, error) {
	value, ok := object[key]
	if !ok {
		return false, ErrObjectMemberNotFound
	}
	_, ok = value.(NullType)
	return ok, nil
}

var ErrObjectMemberNotFound = errors.New("object member not found")
