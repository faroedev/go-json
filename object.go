package json

import (
	"fmt"
	"strconv"
)

// Represents a JSON object.
// The zero value is a ready-to-use empty object.
type ObjectStruct struct {
	strings map[string]string
	numbers map[string]string
	bools   map[string]bool
	nulls   map[string]struct{}
	objects map[string]ObjectStruct
	arrays  map[string]ArrayStruct
	keys    []string
}

func (object *ObjectStruct) Has(key string) bool {
	if _, ok := object.strings[key]; ok {
		return true
	}
	if _, ok := object.numbers[key]; ok {
		return true
	}
	if _, ok := object.bools[key]; ok {
		return true
	}
	if _, ok := object.nulls[key]; ok {
		return true
	}
	if _, ok := object.objects[key]; ok {
		return true
	}
	if _, ok := object.arrays[key]; ok {
		return true
	}
	return false
}

func (object *ObjectStruct) addKey(key string) {
	if _, ok := object.strings[key]; ok {
		delete(object.strings, key)
	} else if _, ok := object.numbers[key]; ok {
		delete(object.numbers, key)
	} else if _, ok := object.nulls[key]; ok {
		delete(object.nulls, key)
	} else if _, ok := object.bools[key]; ok {
		delete(object.bools, key)
	} else if _, ok := object.objects[key]; ok {
		delete(object.objects, key)
	} else if _, ok := object.arrays[key]; ok {
		delete(object.arrays, key)
	} else {
		object.keys = append(object.keys, key)
	}
}

// Set a member with a JSON string value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetString(key string, value string) {
	object.addKey(key)
	if object.strings == nil {
		object.strings = map[string]string{}
	}
	object.strings[key] = value
}

// Returns an error if the key doesn't exist or the value isn't a JSON string.
func (object *ObjectStruct) GetString(key string) (string, error) {
	if object.strings == nil {
		return "", fmt.Errorf("no matching member")
	}
	value, ok := object.strings[key]
	if !ok {
		return "", fmt.Errorf("no matching member")
	}
	return value, nil
}

// Set a member with a JSON number value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetNumber(key string, value string) {
	object.addKey(key)
	if object.numbers == nil {
		object.numbers = map[string]string{}
	}
	object.numbers[key] = value
}

// Returns an error if the key doesn't exist or the value isn't a JSON number
func (object *ObjectStruct) GetNumber(key string) (string, error) {
	if object.numbers == nil {
		return "", fmt.Errorf("no matching member")
	}
	value, ok := object.numbers[key]
	if !ok {
		return "", fmt.Errorf("no matching member")
	}
	return value, nil
}

// Set a member with a JSON number value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetInt(key string, value int) {
	object.SetNumber(key, strconv.Itoa(value))
}

// Returns an error if the key doesn't exist,
// the value isn't a JSON number,
// or the JSON number cannot be represented as an int.
func (object *ObjectStruct) GetInt(key string) (int, error) {
	if object.numbers == nil {
		return 0, fmt.Errorf("no matching member")
	}
	value, ok := object.numbers[key]
	if !ok {
		return 0, fmt.Errorf("no matching member")
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("no matching member")
	}
	return parsed, nil
}

// Set a member with a JSON number value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetInt32(key string, value int32) {
	object.SetNumber(key, strconv.FormatInt(int64(value), 10))
}

// Returns an error if the key doesn't exist,
// the value isn't a JSON number,
// or the JSON number cannot be represented as an int32.
func (object *ObjectStruct) GetInt32(key string) (int32, error) {
	if object.numbers == nil {
		return 0, fmt.Errorf("no matching member")
	}
	value, ok := object.numbers[key]
	if !ok {
		return 0, fmt.Errorf("no matching member")
	}
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("no matching member")
	}
	return int32(parsed), nil
}

// Set a member with a JSON number value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetInt64(key string, value int64) {
	object.SetNumber(key, strconv.FormatInt(value, 10))
}

// Returns an error if the key doesn't exist,
// the value isn't a JSON number,
// or the JSON number cannot be represented as an int64.
func (object *ObjectStruct) GetInt64(key string) (int64, error) {
	if object.numbers == nil {
		return 0, fmt.Errorf("no matching member")
	}
	value, ok := object.numbers[key]
	if !ok {
		return 0, fmt.Errorf("no matching member")
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("no matching member")
	}
	return parsed, nil
}

// Set a member with a JSON boolean value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetBool(key string, value bool) {
	object.addKey(key)
	if object.bools == nil {
		object.bools = map[string]bool{}
	}
	object.bools[key] = value
}

// Returns an error if the key doesn't exist or the value isn't a JSON boolean.
func (object *ObjectStruct) GetBool(key string) (bool, error) {
	if object.bools == nil {
		return false, fmt.Errorf("no matching member")
	}
	value, ok := object.bools[key]
	if !ok {
		return false, fmt.Errorf("no matching member")
	}
	return value, nil
}

// Set a member with a JSON object value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetJSONObject(key string, value ObjectStruct) {
	object.addKey(key)
	if object.objects == nil {
		object.objects = map[string]ObjectStruct{}
	}
	object.objects[key] = value
}

// Returns an error if the key doesn't exist or the value isn't a JSON object.
func (object *ObjectStruct) GetJSONObject(key string) (ObjectStruct, error) {
	if object.objects == nil {
		return ObjectStruct{}, fmt.Errorf("no matching member")
	}
	value, ok := object.objects[key]
	if !ok {
		return ObjectStruct{}, fmt.Errorf("no matching member")
	}
	return value, nil
}

// Set a member with a JSON array value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetJSONArray(key string, value ArrayStruct) {
	object.addKey(key)
	if object.arrays == nil {
		object.arrays = map[string]ArrayStruct{}
	}
	object.arrays[key] = value
}

// Returns an error if the key doesn't exist or the value isn't a JSON array.
func (object *ObjectStruct) GetJSONArray(key string) (ArrayStruct, error) {
	if object.arrays == nil {
		return ArrayStruct{}, fmt.Errorf("no matching member")
	}
	value, ok := object.arrays[key]
	if !ok {
		return ArrayStruct{}, fmt.Errorf("no matching member")
	}
	return value, nil
}

// Set a member with a JSON null value.
// Overrides any member with the same name.
func (object *ObjectStruct) SetNull(key string) {
	object.addKey(key)
	if object.nulls == nil {
		object.nulls = map[string]struct{}{}
	}
	object.nulls[key] = struct{}{}
}

// Returns an error if the key doesn't exist.
func (object *ObjectStruct) IsNull(key string) (bool, error) {
	_, ok := object.nulls[key]
	if !ok {
		return false, fmt.Errorf("no matching member")
	}
	return true, nil
}

// Returns true if the key exists and the value is null.
func (object *ObjectStruct) ExistsAndIsNull(key string) bool {
	_, ok := object.nulls[key]
	if !ok {
		return false
	}
	return ok
}

// Encodes the object using ObjectBuilderStruct.
// Embedded objects are encoded with ObjectStruct.String().
// Embedded arrays are encoded with ArrayStruct.String().
func (object *ObjectStruct) String() string {
	builder := ObjectBuilderStruct{}
	for _, key := range object.keys {
		if value, ok := object.strings[key]; ok {
			builder.AddString(key, value)
			continue
		}
		if value, ok := object.numbers[key]; ok {
			builder.AddJSON(key, value)
			continue
		}
		if value, ok := object.bools[key]; ok {
			builder.AddBool(key, value)
			continue
		}
		if _, ok := object.nulls[key]; ok {
			builder.AddNull(key)
			continue
		}
		if value, ok := object.objects[key]; ok {
			builder.AddJSON(key, value.String())
			continue
		}
		if value, ok := object.arrays[key]; ok {
			builder.AddJSON(key, value.String())
			continue
		}
	}
	return builder.Done()
}
