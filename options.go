package fi

// 全局数据库驱动类型
var __driver = "mysql"

const (
	DriverMySQL      = "mysql"    // MySQL驱动，默认给字段名加反引号
	DriverPostgreSQL = "postgres" // PostgreSQL驱动，默认不给字段名加反引号
)

// SetDriver 设置数据库驱动类型
// 影响字段名的引用方式（是否添加反引号）
// 参数:
//   - driver: 数据库驱动类型，支持mysql和postgres
func SetDriver(driver string) {
	switch driver {
	case DriverMySQL, DriverPostgreSQL:
		__driver = driver
	}
}

type (
	// Option 配置选项函数类型
	// 用于自定义过滤器的行为
	Option func(*option)

	// option 过滤器配置选项
	option struct {
		SkipZeroValue bool // 是否跳过零值字段
		Size          int  // 预估字段数量，用于预分配参数切片容量
		Quote         bool // 是否给字段名添加反引号
	}
)

// WithSkipZeroValue 设置是否跳过零值字段
// 当enabled为true时，零值字段不会被添加到查询条件中
// 参数:
//   - enabled: 是否跳过零值
//
// 返回:
//   - Option: 配置选项函数
func WithSkipZeroValue(enabled bool) Option {
	return func(o *option) {
		o.SkipZeroValue = enabled
	}
}

// WithSize 设置预估字段数量
// 用于预分配参数切片容量，提高性能
// 参数:
//   - size: 预估的字段数量
//
// 返回:
//   - Option: 配置选项函数
func WithSize(size int) Option {
	return func(o *option) {
		o.Size = size
	}
}

// WithQuote 设置是否给字段名添加反引号
// 用于兼容不同数据库的字段名引用方式
// 参数:
//   - enabled: 是否添加反引号
//
// 返回:
//   - Option: 配置选项函数
func WithQuote(enabled bool) Option {
	return func(o *option) {
		o.Quote = enabled
	}
}
