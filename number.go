package json

import (
	"strconv"
)

type NumberType string

func NewNumber[T Integer](i T) NumberType {
	switch any(i).(type) {
	case int, int8, int16, int32, int64:
		return NumberType(strconv.FormatInt(int64(i), 10))
	default:
		return NumberType(strconv.FormatUint(uint64(i), 10))
	}
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func (number NumberType) Int() (int, bool) {
	value, err := strconv.ParseInt(string(number), 10, 0)
	if err != nil {
		return 0, false
	}
	return int(value), true
}

func (number NumberType) Int8() (int8, bool) {
	value, err := strconv.ParseInt(string(number), 10, 8)
	if err != nil {
		return 0, false
	}
	return int8(value), true
}

func (number NumberType) Int16() (int16, bool) {
	value, err := strconv.ParseInt(string(number), 10, 16)
	if err != nil {
		return 0, false
	}
	return int16(value), true
}

func (number NumberType) Int32() (int32, bool) {
	value, err := strconv.ParseInt(string(number), 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(value), true
}

func (number NumberType) Int64() (int64, bool) {
	value, err := strconv.ParseInt(string(number), 10, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func (number NumberType) Uint() (uint, bool) {
	value, err := strconv.ParseUint(string(number), 10, 0)
	if err != nil {
		return 0, false
	}
	return uint(value), true
}

func (number NumberType) Uint8() (uint8, bool) {
	value, err := strconv.ParseUint(string(number), 10, 8)
	if err != nil {
		return 0, false
	}
	return uint8(value), true
}

func (number NumberType) Uint16() (uint16, bool) {
	value, err := strconv.ParseUint(string(number), 10, 16)
	if err != nil {
		return 0, false
	}
	return uint16(value), true
}

func (number NumberType) Uint32() (uint32, bool) {
	value, err := strconv.ParseUint(string(number), 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(value), true
}

func (number NumberType) Uint64() (uint64, bool) {
	value, err := strconv.ParseUint(string(number), 10, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}
