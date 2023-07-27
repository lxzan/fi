package internal

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestSelectValue(t *testing.T) {
	assert.Equal(t, 1, SelectValue(true, 1, 0))
	assert.Equal(t, 0, SelectValue(false, 1, 0))
}

func TestIsNil(t *testing.T) {
	assert.Equal(t, true, IsNil(nil))
	assert.Equal(t, false, IsNil(1))

	var arr []int
	assert.Equal(t, true, IsNil(arr))

	var num *int
	assert.Equal(t, true, IsNil(num))

	var tcpConn *net.TCPConn
	var conn net.Conn = tcpConn
	assert.Equal(t, true, IsNil(conn))
}

func TestIsZero(t *testing.T) {
	assert.Equal(t, true, IsZero(int(0)))
	assert.Equal(t, true, IsZero(int64(0)))
	assert.Equal(t, true, IsZero(int32(0)))
	assert.Equal(t, true, IsZero(int16(0)))
	assert.Equal(t, true, IsZero(int8(0)))
	assert.Equal(t, true, IsZero(uint(0)))
	assert.Equal(t, true, IsZero(uint64(0)))
	assert.Equal(t, true, IsZero(uint32(0)))
	assert.Equal(t, true, IsZero(uint16(0)))
	assert.Equal(t, true, IsZero(uint8(0)))
	assert.Equal(t, true, IsZero(0.000001))
	assert.Equal(t, true, IsZero(float32(0.000001)))
	assert.Equal(t, true, IsZero(""))
	assert.Equal(t, true, IsZero(false))
	assert.Equal(t, true, IsZero(time.Time{}))

	type Int int
	assert.Equal(t, true, IsZero(Int(0)))
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
