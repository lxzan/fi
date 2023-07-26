package internal

import (
	"reflect"
	"time"
)

func SelectValue[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr:
		return rv.IsNil()
	case reflect.Slice:
		return rv.Len() == 0
	default:
		return false
	}
}

func IsZero(v any) bool {
	switch value := v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return value == 0
	case float32, float64:
		return value == 0.0
	case string:
		return value == ""
	case time.Time:
		return value.IsZero()
	default:
		return false
	}
}
