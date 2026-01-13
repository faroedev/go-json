package json

import (
	"fmt"
	"strconv"
)

// Represents a JSON array.
type ArrayStruct struct {
	strings map[int]string
	numbers map[int]string
	bools   map[int]bool
	nulls   map[int]struct{}
	objects map[int]ObjectStruct
	arrays  map[int]ArrayStruct
	length  int
}

func NewArray() ArrayStruct {
	array := ArrayStruct{
		strings: map[int]string{},
		numbers: map[int]string{},
		bools:   map[int]bool{},
		nulls:   map[int]struct{}{},
		objects: map[int]ObjectStruct{},
		arrays:  map[int]ArrayStruct{},
		length:  0,
	}
	return array
}

func (array *ArrayStruct) Length() int {
	return array.length
}

func (array *ArrayStruct) removeElement(index int) {
	delete(array.strings, index)
	delete(array.numbers, index)
	delete(array.bools, index)
	delete(array.nulls, index)
	delete(array.objects, index)
	delete(array.arrays, index)
}

// Sets a JSON string value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetString(index int, value string) {
	if index >= array.length {
		panic("out of bounds")
	}
	array.removeElement(index)
	array.strings[index] = value
}

// Appends a JSON string value at the end of the array.
func (array *ArrayStruct) AddString(value string) {
	array.strings[array.length] = value
	array.length++
}

// Returns an error if an item doesn't exist in the index or the value isn't a JSON string.
func (array *ArrayStruct) GetString(index int) (string, error) {
	value, ok := array.strings[index]
	if !ok {
		return "", fmt.Errorf("no matching member")
	}
	return value, nil
}

// Sets a JSON number value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetNumber(index int, value string) {
	if index >= array.length {
		panic("out of bounds")
	}
	array.removeElement(index)
	array.numbers[index] = value
}

// Appends a JSON number value at the end of the array.
func (array *ArrayStruct) AddNumber(value string) {
	array.numbers[array.length] = value
	array.length++
}

// Returns an error if an item doesn't exist in the index or the value isn't a JSON number.
func (array *ArrayStruct) GetNumber(key int) (string, error) {
	value, ok := array.numbers[key]
	if !ok {
		return "", fmt.Errorf("no matching member")
	}
	return value, nil
}

// Sets a JSON number value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetInt(index int, value int) {
	array.SetNumber(index, strconv.Itoa(value))
}

// Appends a JSON number value at the end of the array.
func (array *ArrayStruct) AddInt(value int) {
	array.AddNumber(strconv.Itoa(value))
}

// Returns an error if an item doesn't exist in the index,
// the value isn't a JSON number,
// or the JSON number cannot be represented as an int.
func (array *ArrayStruct) GetInt(key int) (int, error) {
	value, err := array.GetNumber(key)
	if err != nil {
		return 0, fmt.Errorf("failed to get number: %s", err.Error())
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int: %s", err.Error())
	}
	return parsed, nil
}

// Sets a JSON number value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetInt64(index int, value int64) {
	array.SetNumber(index, strconv.FormatInt(value, 10))
}

// Appends a JSON number value at the end of the array.
func (array *ArrayStruct) AddInt64(value int64) {
	array.AddNumber(strconv.FormatInt(value, 10))
}

// Returns an error if an item doesn't exist in the index,
// the value isn't a JSON number,
// or the JSON number cannot be represented as an int64.
func (array *ArrayStruct) GetInt64(key int) (int64, error) {
	value, err := array.GetNumber(key)
	if err != nil {
		return 0, fmt.Errorf("failed to get number: %s", err.Error())
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int64: %s", err.Error())
	}
	return parsed, nil
}

// Sets a JSON number value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetInt32(index int, value int32) {
	array.SetNumber(index, strconv.FormatInt(int64(value), 10))
}

// Appends a JSON number value at the end of the array.
func (array *ArrayStruct) AddInt32(value int32) {
	array.AddNumber(strconv.FormatInt(int64(value), 10))
}

// Returns an error if an item doesn't exist in the index,
// the value isn't a JSON number,
// or the JSON number cannot be represented as an int32.
func (array *ArrayStruct) GetInt32(key int) (int32, error) {
	value, err := array.GetNumber(key)
	if err != nil {
		return 0, fmt.Errorf("failed to get number: %s", err.Error())
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int32: %s", err.Error())
	}
	return int32(parsed), nil
}

// Sets a JSON boolean value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetBool(index int, value bool) {
	if index >= array.length {
		panic("out of bounds")
	}
	array.removeElement(index)
	array.bools[index] = value
}

// Appends a JSON boolean value at the end of the array.
func (array *ArrayStruct) AddBool(value bool) {
	array.bools[array.length] = value
	array.length++
}

// Returns an error if an item doesn't exist in the index or the value isn't a JSON boolean.
func (array *ArrayStruct) GetBool(index int) (bool, error) {
	value, ok := array.bools[index]
	if !ok {
		return false, fmt.Errorf("no matching member")
	}
	return value, nil
}

// Sets a JSON object value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetJSONObject(index int, value ObjectStruct) {
	if index >= array.length {
		panic("out of bounds")
	}
	array.removeElement(index)
	array.objects[index] = value
}

// Appends a JSON object value at the end of the array.
func (array *ArrayStruct) AddJSONObject(value ObjectStruct) {
	array.objects[array.length] = value
	array.length++
}

// Returns an error if an item doesn't exist in the index or the value isn't a JSON object.
func (array *ArrayStruct) GetJSONObject(index int) (ObjectStruct, error) {
	value, ok := array.objects[index]
	if !ok {
		return value, fmt.Errorf("no matching member")
	}
	return value, nil
}

// Sets a JSON array value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetJSONArray(index int, value ArrayStruct) {
	if index >= array.length {
		panic("out of bounds")
	}
	array.removeElement(index)
	array.arrays[index] = value
}

// Appends a JSON array value at the end of the array.
func (array *ArrayStruct) AddJSONArray(value ArrayStruct) {
	array.arrays[array.length] = value
	array.length++
}

// Returns an error if an item doesn't exist in the index or the value isn't a JSON array.
func (array *ArrayStruct) GetJSONArray(index int) (ArrayStruct, error) {
	value, ok := array.arrays[index]
	if !ok {
		return value, fmt.Errorf("no matching member")
	}
	return value, nil
}

// Sets a JSON null value at index.
// Panics if the index is out of bounds.
func (array *ArrayStruct) SetNull(index int) {
	if index >= array.length {
		panic("out of bounds")
	}
	array.removeElement(index)
	array.nulls[index] = struct{}{}
}

// Appends a JSON null value at the end of the array.
func (array *ArrayStruct) AddNull() {
	array.nulls[array.length] = struct{}{}
	array.length++
}

// Returns an error if an item doesn't exist in the index.
func (array *ArrayStruct) IsNull(index int) (bool, error) {
	_, ok := array.nulls[index]
	if !ok {
		return false, fmt.Errorf("no matching member")
	}
	return true, nil
}

// Returns true if the value at the index is null.
func (array *ArrayStruct) ExistsAndIsNull(index int) bool {
	_, ok := array.nulls[index]
	return ok
}

// Encodes the array using ArrayBuilderStruct.
// Embedded objects are encoded with ObjectStruct.String().
// Embedded arrays are encoded with ArrayStruct.String().
func (array *ArrayStruct) String(stringCharacterEscapingBehavior StringCharacterEscapingBehaviorInterface) string {
	builder := NewArrayBuilder(stringCharacterEscapingBehavior)
	for i := range array.length {
		if value, ok := array.strings[i]; ok {
			builder.AddString(value)
			continue
		}
		if value, ok := array.numbers[i]; ok {
			builder.AddJSON(value)
			continue
		}
		if value, ok := array.bools[i]; ok {
			builder.AddBool(value)
			continue
		}
		if _, ok := array.nulls[i]; ok {
			builder.AddNull()
			continue
		}
		if value, ok := array.objects[i]; ok {
			builder.AddJSON(value.String(stringCharacterEscapingBehavior))
			continue
		}
		if value, ok := array.arrays[i]; ok {
			builder.AddJSON(value.String(stringCharacterEscapingBehavior))
			continue
		}
	}
	return builder.Done()
}
