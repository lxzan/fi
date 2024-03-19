package fi

var __driver = "mysql"

const (
	DriverMySQL      = "mysql"    // 默认给key加反引号
	DriverPostgreSQL = "postgres" // 默认不给key加反引号
)

// SetDriver 设置数据库驱动
func SetDriver(driver string) {
	switch driver {
	case DriverMySQL, DriverPostgreSQL:
		__driver = driver
	}
}

type (
	Option func(*option)

	option struct {
		SkipZeroValue bool // 是否跳过空值
		Size          int  // 预估字段数量
		Quote         bool // 是否加反引号
	}
)

func WithSkipZeroValue(enabled bool) Option {
	return func(o *option) {
		o.SkipZeroValue = enabled
	}
}

func WithSize(size int) Option {
	return func(o *option) {
		o.Size = size
	}
}

func WithQuote(enabled bool) Option {
	return func(o *option) {
		o.Quote = enabled
	}
}
