package fi

import (
	"reflect"
	"strings"

	"github.com/lxzan/fi/internal"
)

// field 字段信息结构体
// 用于存储从结构体标签中解析出的字段信息
type field struct {
	column string // 数据库字段名
	cmp    string // 比较操作符
	value  any    // 字段值
}

// GetFilter 通过反射从结构体生成过滤器
// 根据结构体字段的fi标签自动生成查询条件
// 参数:
//   - v: 包含fi标签的结构体实例（支持指针和值类型）
//
// 返回:
//   - *Filter: 生成的过滤器实例
func GetFilter(v any) *Filter {
	var f = NewFilter()
	var g = generator{}
	var vs = reflect.ValueOf(v)

	// 处理指针类型，获取实际的值
	if vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}

	// 遍历结构体字段，解析fi标签并生成查询条件
	g.getField(vs, func(item *field) {
		switch item.cmp {
		case "eq":
			f.Eq(item.column, item.value)
		case "not_eq":
			f.NotEq(item.column, item.value)
		case "gt":
			f.Gt(item.column, item.value)
		case "lt":
			f.Lt(item.column, item.value)
		case "gte":
			f.Gte(item.column, item.value)
		case "lte":
			f.Lte(item.column, item.value)
		case "like":
			if value, ok := item.value.(string); ok {
				f.Like(item.column, value)
			}
		case "not_like":
			if value, ok := item.value.(string); ok {
				f.NotLike(item.column, value)
			}
		case "in":
			f.In(item.column, item.value)
		case "not_in":
			f.NotIn(item.column, item.value)
		}
	})

	return f
}

// generator 标签解析器
// 用于解析结构体字段的fi标签
type generator struct {
	kv [2]string // 键值对缓存，用于解析标签中的key=value格式
}

// split 按分隔符分割字符串
// 参数:
//   - s: 要分割的字符串
//   - sep: 分隔符
//   - f: 对每个分割后的子字符串执行的回调函数
func (c *generator) split(s string, sep byte, f func(str string)) {
	n := len(s)
	start, end := 0, 0
	for i, _ := range s {
		if s[i] == sep || i == n-1 {
			end = i
			if i == n-1 && s[i] != sep {
				end++
			}
			f(s[start:end])
			start = i + 1
		}
	}
}

// getField 递归遍历结构体字段，解析fi标签
// 参数:
//   - vs: 结构体的反射值
//   - callback: 对每个有效字段执行的回调函数
func (c *generator) getField(vs reflect.Value, callback func(*field)) {
	ts := vs.Type()
	f := &field{}

	for i := 0; i < vs.NumField(); i++ {
		var tf = ts.Field(i)
		var vf = vs.Field(i)
		var tag = tf.Tag.Get("fi")

		// 跳过标记为"-"的字段
		if tag == "-" {
			continue
		}

		// 处理指针类型字段
		if vf.Type().Kind() == reflect.Ptr {
			// 跳过nil指针或私有字段
			if vf.IsNil() || !tf.IsExported() {
				continue
			}

			vf = vf.Elem()
			// 如果是结构体指针，递归处理
			if vf.Kind() == reflect.Struct {
				c.getField(vf, callback)
				continue
			}
		}

		// 解析标签并设置字段值
		c.splitTag(tf.Name, tag, f)
		f.value = vf.Interface()
		callback(f)
	}
}

// splitTag 解析fi标签内容
// 支持的标签格式: column=field_name;cmp=operator
// 参数:
//   - name: 字段名（用于生成默认的数据库字段名）
//   - tag: fi标签内容
//   - f: 要填充的field结构体
func (c *generator) splitTag(name string, tag string, f *field) {
	// 设置默认值
	f.cmp = "eq"
	f.column = ""
	f.value = nil

	// 按分号分割标签内容
	c.split(tag, ';', func(item string) {
		item = strings.TrimSpace(item)
		if item == "" {
			return
		}

		// 解析key=value格式
		index := 0
		c.kv[0], c.kv[1] = "", ""
		c.split(item, '=', func(str string) {
			if index < 2 {
				c.kv[index] = str
			}
			index++
		})

		c.kv[0] = strings.TrimSpace(c.kv[0])
		c.kv[1] = strings.TrimSpace(c.kv[1])

		// 根据key设置对应的字段
		switch c.kv[0] {
		case "column":
			f.column = c.kv[1]
		case "cmp":
			f.cmp = c.kv[1]
		}
	})

	// 如果没有指定column，使用字段名的下划线格式
	if f.column == "" {
		f.column = internal.ToSnakeCase(name)
	}
}
