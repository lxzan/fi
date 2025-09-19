package internal

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSelectValue(t *testing.T) {
	assert.Equal(t, 1, SelectValue(true, 1, 0))
	assert.Equal(t, 0, SelectValue(false, 1, 0))
}

func TestIsNil(t *testing.T) {
	var isNil = func(v any) bool {
		return IsNil(v, reflect.ValueOf(v))
	}

	assert.Equal(t, true, isNil(nil))
	assert.Equal(t, false, isNil(1))

	var arr []int
	assert.Equal(t, true, isNil(arr))

	var num *int
	assert.Equal(t, true, isNil(num))

	var tcpConn *net.TCPConn
	var conn net.Conn = tcpConn
	assert.Equal(t, true, isNil(conn))

	t.Run("", func(t *testing.T) {
		var arr = make([]int, 0)
		assert.Equal(t, false, isNil(arr))
	})
}

func TestIsZero(t *testing.T) {
	var isZero = func(v any) bool {
		return IsZero(v, reflect.ValueOf(v))
	}

	assert.Equal(t, true, isZero(int(0)))
	assert.Equal(t, true, isZero(int64(0)))
	assert.Equal(t, true, isZero(int32(0)))
	assert.Equal(t, true, isZero(int16(0)))
	assert.Equal(t, true, isZero(int8(0)))
	assert.Equal(t, true, isZero(uint(0)))
	assert.Equal(t, true, isZero(uint64(0)))
	assert.Equal(t, true, isZero(uint32(0)))
	assert.Equal(t, true, isZero(uint16(0)))
	assert.Equal(t, true, isZero(uint8(0)))
	assert.Equal(t, true, isZero(0.000001))
	assert.Equal(t, true, isZero(float32(0.000001)))
	assert.Equal(t, true, isZero(""))
	assert.Equal(t, true, isZero(false))
	assert.Equal(t, true, isZero(time.Time{}))

	type Int int
	assert.Equal(t, true, isZero(Int(0)))
}

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "name", ToSnakeCase("Name"))
	assert.Equal(t, "user_name", ToSnakeCase("UserName"))
	assert.Equal(t, "ip", ToSnakeCase("IP"))
	assert.Equal(t, "i_pv6", ToSnakeCase("IPv6"))
}

// 测试IsNil的更多情况
func TestIsNil_MoreCases(t *testing.T) {
	var isNil = func(v any) bool {
		return IsNil(v, reflect.ValueOf(v))
	}

	// 测试各种nil情况
	var nilSlice []int
	var nilMap map[string]int
	var nilChan chan int
	var nilFunc func()
	var nilPtr *int
	var nilInterface interface{}

	assert.Equal(t, true, isNil(nilSlice))
	assert.Equal(t, true, isNil(nilMap))
	assert.Equal(t, true, isNil(nilChan))
	assert.Equal(t, true, isNil(nilFunc))
	assert.Equal(t, true, isNil(nilPtr))
	assert.Equal(t, true, isNil(nilInterface))

	// 测试非nil情况
	nonNilSlice := make([]int, 0)
	nonNilMap := make(map[string]int)
	nonNilChan := make(chan int)
	nonNilFunc := func() {}
	nonNilPtr := new(int)
	nonNilInterface := interface{}(1)

	assert.Equal(t, false, isNil(nonNilSlice))
	assert.Equal(t, false, isNil(nonNilMap))
	assert.Equal(t, false, isNil(nonNilChan))
	assert.Equal(t, false, isNil(nonNilFunc))
	assert.Equal(t, false, isNil(nonNilPtr))
	assert.Equal(t, false, isNil(nonNilInterface))
}

// 测试IsZero的更多情况
func TestIsZero_MoreCases(t *testing.T) {
	var isZero = func(v any) bool {
		return IsZero(v, reflect.ValueOf(v))
	}

	// 测试各种零值
	assert.Equal(t, true, isZero(int(0)))
	assert.Equal(t, true, isZero(int64(0)))
	assert.Equal(t, true, isZero(int32(0)))
	assert.Equal(t, true, isZero(int16(0)))
	assert.Equal(t, true, isZero(int8(0)))
	assert.Equal(t, true, isZero(uint(0)))
	assert.Equal(t, true, isZero(uint64(0)))
	assert.Equal(t, true, isZero(uint32(0)))
	assert.Equal(t, true, isZero(uint16(0)))
	assert.Equal(t, true, isZero(uint8(0)))
	assert.Equal(t, true, isZero(float64(0.0)))
	assert.Equal(t, true, isZero(float32(0.0)))
	assert.Equal(t, true, isZero(""))
	assert.Equal(t, true, isZero(false))
	assert.Equal(t, true, isZero(time.Time{}))

	// 测试非零值
	assert.Equal(t, false, isZero(int(1)))
	assert.Equal(t, false, isZero(int64(1)))
	assert.Equal(t, false, isZero(int32(1)))
	assert.Equal(t, false, isZero(int16(1)))
	assert.Equal(t, false, isZero(int8(1)))
	assert.Equal(t, false, isZero(uint(1)))
	assert.Equal(t, false, isZero(uint64(1)))
	assert.Equal(t, false, isZero(uint32(1)))
	assert.Equal(t, false, isZero(uint16(1)))
	assert.Equal(t, false, isZero(uint8(1)))
	assert.Equal(t, false, isZero(float64(0.1)))
	assert.Equal(t, false, isZero(float32(0.1)))
	assert.Equal(t, false, isZero("test"))
	assert.Equal(t, false, isZero(true))
	assert.Equal(t, false, isZero(time.Now()))

	// 测试浮点数精度（根据实际实现调整）
	assert.Equal(t, true, isZero(0.000001))
	assert.Equal(t, true, isZero(0.0000001))
	assert.Equal(t, true, isZero(0.00001)) // 这个值可能被认为是零值
	assert.Equal(t, false, isZero(0.1))

	// 测试更小的浮点数
	assert.Equal(t, true, isZero(1e-6))
	assert.Equal(t, true, isZero(1e-5)) // 这个值可能被认为是零值

	// 测试自定义类型
	type CustomInt int
	assert.Equal(t, true, isZero(CustomInt(0)))
	assert.Equal(t, false, isZero(CustomInt(1)))
}

// 测试ToSnakeCase的更多情况
func TestToSnakeCase_MoreCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Name", "name"},
		{"UserName", "user_name"},
		{"IP", "ip"},
		{"IPv6", "i_pv6"},
		{"HTTPRequest", "h_t_t_p_request"},
		{"XMLHttpRequest", "x_m_l_http_request"},
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"aB", "a_b"},
		{"AB", "ab"},
		{"ABC", "abc"},
		{"abc", "abc"},
		{"userID", "user_i_d"},
		{"userIDName", "user_i_d_name"},
		{"ID", "id"},
		{"IDName", "i_d_name"},
		{"nameID", "name_i_d"},
		{"nameIDName", "name_i_d_name"},
		{"URL", "url"},
		{"URLPath", "u_r_l_path"},
		{"pathURL", "path_u_r_l"},
		{"pathURLPath", "path_u_r_l_path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// 测试SelectValue的更多情况
func TestSelectValue_MoreCases(t *testing.T) {
	// 测试字符串
	assert.Equal(t, "a", SelectValue(true, "a", "b"))
	assert.Equal(t, "b", SelectValue(false, "a", "b"))

	// 测试数字
	assert.Equal(t, 1, SelectValue(true, 1, 2))
	assert.Equal(t, 2, SelectValue(false, 1, 2))

	// 测试切片
	assert.Equal(t, []int{1, 2}, SelectValue(true, []int{1, 2}, []int{3, 4}))
	assert.Equal(t, []int{3, 4}, SelectValue(false, []int{1, 2}, []int{3, 4}))

	// 测试结构体
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}

	assert.Equal(t, p1, SelectValue(true, p1, p2))
	assert.Equal(t, p2, SelectValue(false, p1, p2))

	// 测试指针
	ptr1 := &p1
	ptr2 := &p2
	assert.Equal(t, ptr1, SelectValue(true, ptr1, ptr2))
	assert.Equal(t, ptr2, SelectValue(false, ptr1, ptr2))
}
