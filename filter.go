package fi

import (
	"official-bk/internal/utils/fi/internal"
	"strings"
)

type Filter struct {
	builder strings.Builder
	skip    bool          // 是否跳过空值
	quote   bool          // 是否加引号
	Args    []interface{} // 参数
}

func NewFilter(options ...Option) *Filter {
	o := &option{SkipZeroValue: true}
	for _, f := range options {
		f(o)
	}

	f := &Filter{skip: o.SkipZeroValue}
	if o.Size > 0 {
		f.Args = make([]interface{}, 0, o.Size)
		f.builder.Grow(20 * o.Size)
	}
	return f
}

func (c *Filter) push(key string, val any, cmp string) *Filter {
	if v, ok := val.(Valuer); ok {
		val = v.Value()
	}

	if internal.IsNil(val) || (c.skip && internal.IsZero(val)) {
		return c
	}

	if c.builder.Len() > 0 {
		c.builder.WriteString(" AND ")
	}

	var hasDot = strings.Contains(key, ".")
	if !hasDot && c.quote {
		c.builder.WriteString("`")
	}
	c.builder.WriteString(key)
	if !hasDot && c.quote {
		c.builder.WriteString("`")
	}
	c.builder.WriteString(" ")
	c.builder.WriteString(cmp)
	switch cmp {
	case "IN", "NOT IN":
		c.builder.WriteString(" (?)")
		c.Args = append(c.Args, val)
	case "IS NULL":
	default:
		c.builder.WriteString(" ?")
		c.Args = append(c.Args, val)
	}
	return c
}

func (c *Filter) Eq(key string, val any) *Filter {
	return c.push(key, val, "=")
}

func (c *Filter) NotEq(key string, val any) *Filter {
	return c.push(key, val, "!=")
}

func (c *Filter) Gt(key string, val any) *Filter {
	return c.push(key, val, ">")
}

func (c *Filter) Lt(key string, val any) *Filter {
	return c.push(key, val, "<")
}

func (c *Filter) Gte(key string, val any) *Filter {
	return c.push(key, val, ">=")
}

func (c *Filter) Lte(key string, val any) *Filter {
	return c.push(key, val, "<=")
}

func (c *Filter) addPercent(str string) string {
	if str == "" {
		return str
	}
	var n = len(str)
	if str[0] == '%' || str[n-1] == '%' {
		return str
	}
	return "%" + str + "%"
}

// Like 如果value前后不包含百分号, 会自动添加; NotLike同理.
func (c *Filter) Like(key string, val string) *Filter {
	return c.push(key, c.addPercent(val), "LIKE")
}

func (c *Filter) NotLike(key string, val string) *Filter {
	return c.push(key, c.addPercent(val), "NOT LIKE")
}

func (c *Filter) In(key string, val any) *Filter {
	return c.push(key, val, "IN")
}

func (c *Filter) NotIn(key string, val any) *Filter {
	return c.push(key, val, "NOT IN")
}

func (c *Filter) IsNull(key string) *Filter {
	return c.push(key, "", "IS NULL")
}

// Customize 追加自定义SQL
func (c *Filter) Customize(layout string, val ...any) *Filter {
	if c.builder.Len() > 0 {
		c.builder.WriteString(" AND ")
	}
	c.builder.WriteString(layout)
	c.Args = append(c.Args, val...)
	return c
}

// WithTimeSelector 时间选择器
// 区间: [startTime, endTime)
func (c *Filter) WithTimeSelector(key string, startTime int64, endTime int64) *Filter {
	if startTime+endTime == 0 {
		return c
	}
	skip := c.skip
	c.skip = false
	c.Gte(key, startTime).Lt(key, endTime)
	c.skip = skip
	return c
}

// GetExpression 获取SQL表达式
func (c *Filter) GetExpression() string {
	if c.builder.Len() == 0 {
		return "1=1"
	}
	return c.builder.String()
}

type (
	Option func(*option)

	option struct {
		SkipZeroValue bool
		Size          int
		Quote         bool
	}
)

func WithSkipZeroValue(skip bool) Option {
	return func(o *option) {
		o.SkipZeroValue = skip
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
