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
		return reflect.ValueOf(v).IsZero()
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
