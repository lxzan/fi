package internal

import (
	"github.com/stretchr/testify/assert"
	"net"
	"reflect"
	"testing"
	"time"
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

func TestIsPrivate(t *testing.T) {
	assert.Equal(t, true, IsPrivate("name"))
}
