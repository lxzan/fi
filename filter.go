package fi

import (
	"reflect"
	"strings"
	"time"

	"github.com/lxzan/fi/internal"
)

// Filter 动态SQL查询条件构造器
// 用于构建动态的SQL WHERE条件，支持链式调用和多种比较操作符
type Filter struct {
	builder   strings.Builder // SQL表达式构建器
	keyPrefix string          // 字段名前缀，用于表别名
	conf      *option         // 配置选项
	Args      []any           // SQL参数列表
}

// NewFilter 新建过滤器实例
// 默认会跳过零值，用于动态查询条件构建
// 参数:
//   - options: 可选的配置项，用于自定义过滤器行为
//
// 返回:
//   - *Filter: 新创建的过滤器实例
func NewFilter(options ...Option) *Filter {
	// 设置默认配置
	o := &option{SkipZeroValue: true, Size: 10, Quote: true}

	// PostgreSQL驱动默认不添加反引号
	if __driver == DriverPostgreSQL {
		o.Quote = false
	}

	// 应用用户自定义配置
	for _, f := range options {
		f(o)
	}

	// 创建过滤器实例
	f := &Filter{conf: o}
	f.builder.Grow(256) // 预分配256字节的缓冲区

	// 预分配参数切片容量
	if o.Size > 0 {
		f.Args = make([]any, 0, o.Size)
	}
	return f
}

// NewQuery 新建查询实例
// 默认不跳过零值，用于拼接静态查询条件
// 与NewFilter的区别：不会自动跳过零值，所有条件都会被添加到SQL中
// 参数:
//   - options: 可选的配置项
//
// 返回:
//   - *Filter: 新创建的查询实例
func NewQuery(options ...Option) *Filter {
	n := len(options)
	switch n {
	case 0:
		// 无额外配置，直接创建不跳过零值的过滤器
		return NewFilter(WithSkipZeroValue(false))
	default:
		// 有额外配置，需要确保SkipZeroValue为false
		opts := make([]Option, 0, 1+len(options))
		opts = append(opts, WithSkipZeroValue(false))
		opts = append(opts, options...)
		return NewFilter(opts...)
	}
}

// Clone 创建过滤器的深拷贝副本
// 返回一个包含相同配置和条件的新的过滤器实例
// 返回:
//   - *Filter: 新的过滤器副本
func (c *Filter) Clone() *Filter {
	var d = &Filter{conf: c.conf, keyPrefix: c.keyPrefix}
	d.Args = append(d.Args, c.Args...)        // 复制参数列表
	d.builder.WriteString(c.builder.String()) // 复制SQL表达式
	return d
}

// push 内部方法，用于添加查询条件到SQL表达式中
// 参数:
//   - key: 字段名
//   - val: 字段值
//   - cmp: 比较操作符
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) push(key string, val any, cmp string) *Filter {
	// 添加字段名前缀（表别名）
	if c.keyPrefix != "" {
		key = c.keyPrefix + key
	}

	// 处理实现了Valuer接口的值
	if v, ok := val.(Valuer); ok {
		val = v.Value()
	}

	// 跳过零值检查（IS NULL操作符除外）
	if cmp != "IS NULL" && c.conf.SkipZeroValue {
		var rv = reflect.ValueOf(val)
		if internal.IsNil(val, rv) || internal.IsZero(val, rv) {
			return c
		}
	}

	// 添加AND连接符（非第一个条件）
	if c.builder.Len() > 0 {
		c.builder.WriteString(" AND ")
	}

	// 处理字段名的反引号
	var hasDot = strings.Contains(key, ".")
	if !hasDot && c.conf.Quote {
		c.builder.WriteString("`")
	}
	c.builder.WriteString(key)
	if !hasDot && c.conf.Quote {
		c.builder.WriteString("`")
	}

	// 添加比较操作符
	c.builder.WriteString(" ")
	c.builder.WriteString(cmp)

	// 根据操作符类型添加参数占位符
	switch cmp {
	case "IN", "NOT IN":
		c.builder.WriteString(" (?)")
		c.Args = append(c.Args, val)
	case "IS NULL":
		// IS NULL不需要参数
	default:
		c.builder.WriteString(" ?")
		c.Args = append(c.Args, val)
	}
	return c
}

// IsEmpty 检查查询条件是否为空
// 返回:
//   - bool: true表示没有添加任何查询条件
func (c *Filter) IsEmpty() bool {
	return c.builder.Len() == 0
}

// Alias 设置字段名前缀（表别名）
// 参数:
//   - name: 表别名
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Alias(name string) *Filter {
	c.keyPrefix = name + "."
	return c
}

// Eq 添加等于条件 (field = value)
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Eq(key string, val any) *Filter {
	return c.push(key, val, "=")
}

// NotEq 添加不等于条件 (field != value)
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) NotEq(key string, val any) *Filter {
	return c.push(key, val, "!=")
}

// Gt 添加大于条件 (field > value)
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Gt(key string, val any) *Filter {
	return c.push(key, val, ">")
}

// Lt 添加小于条件 (field < value)
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Lt(key string, val any) *Filter {
	return c.push(key, val, "<")
}

// Gte 添加大于等于条件 (field >= value)
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Gte(key string, val any) *Filter {
	return c.push(key, val, ">=")
}

// Lte 添加小于等于条件 (field <= value)
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Lte(key string, val any) *Filter {
	return c.push(key, val, "<=")
}

// addPercent 为LIKE查询添加百分号通配符
// 如果字符串前后不包含百分号，会自动添加
// 参数:
//   - str: 原始字符串
//
// 返回:
//   - string: 添加了百分号的字符串
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

// Like 添加LIKE模糊查询条件 (field LIKE %value%)
// 如果value前后不包含百分号，会自动添加
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Like(key string, val string) *Filter {
	return c.push(key, c.addPercent(val), "LIKE")
}

// NotLike 添加NOT LIKE模糊查询条件 (field NOT LIKE %value%)
// 如果value前后不包含百分号，会自动添加
// 参数:
//   - key: 字段名
//   - val: 字段值
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) NotLike(key string, val string) *Filter {
	return c.push(key, c.addPercent(val), "NOT LIKE")
}

// In 添加IN查询条件 (field IN (values...))
// 参数:
//   - key: 字段名
//   - val: 值列表（切片或数组）
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) In(key string, val any) *Filter {
	return c.push(key, val, "IN")
}

// NotIn 添加NOT IN查询条件 (field NOT IN (values...))
// 参数:
//   - key: 字段名
//   - val: 值列表（切片或数组）
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) NotIn(key string, val any) *Filter {
	return c.push(key, val, "NOT IN")
}

// IsNull 添加IS NULL查询条件 (field IS NULL)
// 参数:
//   - key: 字段名
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) IsNull(key string) *Filter {
	return c.push(key, "", "IS NULL")
}

// Customize 追加自定义SQL条件
// 用于添加复杂的自定义查询条件
// 参数:
//   - layout: 自定义SQL模板，支持占位符
//   - val: 参数值列表
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Customize(layout string, val ...any) *Filter {
	if c.builder.Len() > 0 {
		c.builder.WriteString(" AND ")
	}
	c.builder.WriteString(layout)
	c.Args = append(c.Args, val...)
	return c
}

// Merge 合并另一个过滤器
// 将另一个过滤器的条件合并到当前过滤器中
// 参数:
//   - f: 要合并的过滤器
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) Merge(f *Filter) *Filter {
	if !f.IsEmpty() {
		return c.Customize(f.GetExpression(), f.Args...)
	}
	return c
}

// WithTimeSelector 添加时间范围查询条件
// 使用毫秒时间戳，查询区间为 [startTime, endTime)
// 参数:
//   - key: 时间字段名
//   - startTime: 开始时间（毫秒时间戳）
//   - endTime: 结束时间（毫秒时间戳）
//
// 返回:
//   - *Filter: 返回自身，支持链式调用
func (c *Filter) WithTimeSelector(key string, startTime int64, endTime int64) *Filter {
	if startTime+endTime == 0 {
		return c
	}
	// 临时关闭零值跳过，确保时间条件被添加
	skip := c.conf.SkipZeroValue
	c.conf.SkipZeroValue = false
	c.Gte(key, time.UnixMilli(startTime)).Lt(key, time.UnixMilli(endTime))
	c.conf.SkipZeroValue = skip
	return c
}

// GetExpression 获取完整的SQL WHERE表达式
// 如果没有添加任何条件，返回 "1=1"
// 返回:
//   - string: SQL WHERE表达式
func (c *Filter) GetExpression() string {
	if c.IsEmpty() {
		return "1=1"
	}
	return c.builder.String()
}
