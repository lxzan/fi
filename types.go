package fi

import "time"

// Valuer 值转换接口
// 用于将自定义类型转换为过滤器可以使用的值
type Valuer interface {
	// Value 返回转换后的值
	Value() any
}

// Timestamp 毫秒时间戳类型
// 实现了Valuer接口，会被自动转换为time.Time类型
// 用于在查询条件中处理毫秒时间戳
type Timestamp int64

// Value 将毫秒时间戳转换为time.Time类型
// 实现Valuer接口，用于在查询条件中自动转换时间格式
// 返回:
//   - any: 转换后的time.Time值
func (c Timestamp) Value() any {
	return time.UnixMilli(int64(c))
}
