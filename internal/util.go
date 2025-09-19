package internal

import (
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// SelectValue 根据条件选择值
// 泛型函数，根据ok条件返回a或b
// 参数:
//   - ok: 选择条件
//   - a: 条件为true时返回的值
//   - b: 条件为false时返回的值
//
// 返回:
//   - T: 根据条件选择的值
func SelectValue[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}

// IsNil 检查值是否为nil
// 支持检查各种类型的nil值，包括指针、切片、映射、通道、函数和接口
// 参数:
//   - v: 要检查的值
//   - rv: 值的反射对象
//
// 返回:
//   - bool: true表示值为nil
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

// IsZero 检查值是否为零值
// 支持检查各种基本类型的零值，包括数字、字符串、时间、布尔值等
// 对于浮点数，使用精度阈值判断是否为零值
// 参数:
//   - v: 要检查的值
//   - rv: 值的反射对象
//
// 返回:
//   - bool: true表示值为零值
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

// 全大写字符串的正则表达式
var reUpperString = regexp.MustCompile(`^[A-Z]+$`)

// ToSnakeCase 将驼峰命名转换为下划线命名
// 例如: UserName -> user_name, IP -> ip, IPv6 -> i_pv6
// 参数:
//   - name: 驼峰命名的字符串
//
// 返回:
//   - string: 下划线命名的字符串
func ToSnakeCase(name string) string {
	// 如果全是大写字母，直接转小写
	if reUpperString.MatchString(name) {
		return strings.ToLower(name)
	}

	buffer := strings.Builder{}
	for i, r := range name {
		if unicode.IsUpper(r) {
			// 不是第一个字符时添加下划线
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
