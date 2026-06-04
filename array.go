package json

import (
	"errors"
	"fmt"
)

type ArrayType []any

func NewArray() ArrayType {
	array := ArrayType{}
	return array
}

// Len returns the number of elements in the array.
func (array ArrayType) Len() int {
	return len(array)
}

// Takes one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
func (array *ArrayType) Add(value any) {
	*array = append(*array, value)
}

// Returns [ErrArrayOutOfBounds] if index is invalid.
// Takes one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
func (array ArrayType) Set(index int, value any) error {
	if index < 0 || index >= len(array) {
		return ErrArrayIndexOutOfBounds
	}
	array[index] = value
	return nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds.
// Returns one of StringType, NumberType, BooleanType, NullType, ObjectType, or ArrayType.
func (array ArrayType) Get(index int) (any, error) {
	if index < 0 || index >= len(array) {
		return nil, ErrArrayIndexOutOfBounds
	}
	return array[index], nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds, or an error if the value isn't a string.
func (array ArrayType) GetString(index int) (StringType, error) {
	value, err := array.Get(index)
	if err != nil {
		return "", err
	}
	jsonString, ok := value.(StringType)
	if !ok {
		return "", fmt.Errorf("value not a string")
	}
	return jsonString, nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds, or an error if the value isn't a number.
func (array ArrayType) GetNumber(index int) (NumberType, error) {
	value, err := array.Get(index)
	if err != nil {
		return "", err
	}
	jsonNumber, ok := value.(NumberType)
	if !ok {
		return "", fmt.Errorf("value not a number")
	}
	return jsonNumber, nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds, or an error if the value isn't a boolean.
func (array ArrayType) GetBoolean(index int) (BooleanType, error) {
	value, err := array.Get(index)
	if err != nil {
		return false, err
	}
	jsonBoolean, ok := value.(BooleanType)
	if !ok {
		return false, fmt.Errorf("value not a boolean")
	}
	return jsonBoolean, nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds, or an error if the value isn't an object.
func (array ArrayType) GetObject(index int) (ObjectType, error) {
	value, err := array.Get(index)
	if err != nil {
		return nil, err
	}
	jsonObject, ok := value.(ObjectType)
	if !ok {
		return nil, fmt.Errorf("value not an object")
	}
	return jsonObject, nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds, or an error if the value isn't an array.
func (array ArrayType) GetArray(index int) (ArrayType, error) {
	value, err := array.Get(index)
	if err != nil {
		return nil, err
	}
	jsonArray, ok := value.(ArrayType)
	if !ok {
		return nil, fmt.Errorf("value not an array")
	}
	return jsonArray, nil
}

// Returns [ErrArrayOutOfBounds] if the index is out of bounds.
func (array ArrayType) IsNull(index int) (bool, error) {
	value, err := array.Get(index)
	if err != nil {
		return false, err
	}
	_, ok := value.(NullType)
	return ok, nil
}

var ErrArrayIndexOutOfBounds = errors.New("array index out of bounds")
