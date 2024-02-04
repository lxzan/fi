package internal

import (
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func SelectValue[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}

func IsNil(v any, rv reflect.Value) bool {
	if v == nil {
		return true
	}

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	default:
		return v == nil
	}
}

func IsZero(v any, rv reflect.Value) bool {
	switch value := v.(type) {
	case int:
		return value == 0
	case int64:
		return value == 0
	case int32:
		return value == 0
	case int16:
		return value == 0
	case int8:
		return value == 0
	case uint:
		return value == 0
	case uint64:
		return value == 0
	case uint32:
		return value == 0
	case uint16:
		return value == 0
	case uint8:
		return value == 0
	case float64:
		return math.Abs(value) <= 1.0e-5
	case float32:
		return math.Abs(float64(value)) <= 1.0e-5
	case string:
		return value == ""
	case time.Time:
		return value.IsZero()
	case bool:
		return value == false
	default:
		return rv.IsZero()
	}
}

var reUpperString = regexp.MustCompile(`^[A-Z]+$`)

// ToSnakeCase 驼峰转下划线
func ToSnakeCase(name string) string {
	if reUpperString.MatchString(name) {
		return strings.ToLower(name)
	}

	buffer := strings.Builder{}
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

func IsPrivate(s string) bool {
	return len(s) > 0 && s[0] >= 'a' && s[0] <= 'z'
}
