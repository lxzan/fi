package fi

import (
	"github.com/lxzan/fi/internal"
	"reflect"
	"strings"
	"time"
)

type Filter struct {
	builder strings.Builder
	conf    *option
	Args    []any // 参数
}

// NewFilter 新建过滤器
// 默认会跳过零值, 用于动态查询条件
func NewFilter(options ...Option) *Filter {
	o := &option{SkipZeroValue: true, Size: 10, Quote: true}
	if __driver == DriverPostgreSQL {
		o.Quote = false
	}
	for _, f := range options {
		f(o)
	}

	f := &Filter{conf: o}
	f.builder.Grow(256)
	if o.Size > 0 {
		f.Args = make([]any, 0, o.Size)
	}
	return f
}

// NewQuery 新建查询
// 默认不跳过零值, 用于拼接静态查询条件
func NewQuery(options ...Option) *Filter {
	n := len(options)
	switch n {
	case 0:
		return NewFilter(WithSkipZeroValue(false))
	default:
		opts := make([]Option, 0, 1+len(options))
		opts = append(opts, WithSkipZeroValue(false))
		opts = append(opts, options...)
		return NewFilter(opts...)
	}
}

// Clone 拷贝一个副本
func (c *Filter) Clone() *Filter {
	var d = &Filter{conf: c.conf}
	d.Args = append(d.Args, c.Args...)
	d.builder.WriteString(c.builder.String())
	return d
}

func (c *Filter) push(key string, val any, cmp string) *Filter {
	if v, ok := val.(Valuer); ok {
		val = v.Value()
	}

	if cmp != "IS NULL" && c.conf.SkipZeroValue {
		var rv = reflect.ValueOf(val)
		if internal.IsNil(val, rv) || internal.IsZero(val, rv) {
			return c
		}
	}

	if c.builder.Len() > 0 {
		c.builder.WriteString(" AND ")
	}

	var hasDot = strings.Contains(key, ".")
	if !hasDot && c.conf.Quote {
		c.builder.WriteString("`")
	}
	c.builder.WriteString(key)
	if !hasDot && c.conf.Quote {
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

// IsEmpty 查询条件是否为空
func (c *Filter) IsEmpty() bool {
	return c.builder.Len() == 0
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

// Merge 合并过滤器
func (c *Filter) Merge(f *Filter) *Filter {
	if !f.IsEmpty() {
		return c.Customize(f.GetExpression(), f.Args...)
	}
	return c
}

// WithTimeSelector 时间选择器, 毫秒时间戳
// 区间: [startTime, endTime)
func (c *Filter) WithTimeSelector(key string, startTime int64, endTime int64) *Filter {
	if startTime+endTime == 0 {
		return c
	}
	skip := c.conf.SkipZeroValue
	c.conf.SkipZeroValue = false
	c.Gte(key, time.UnixMilli(startTime)).Lt(key, time.UnixMilli(endTime))
	c.conf.SkipZeroValue = skip
	return c
}

// GetExpression 获取SQL表达式
func (c *Filter) GetExpression() string {
	if c.IsEmpty() {
		return "1=1"
	}
	return c.builder.String()
}
