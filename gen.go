package fi

import (
	"github.com/lxzan/fi/internal"
	"reflect"
	"strings"
)

type field struct {
	column string
	cmp    string
	value  any
}

func GetFilter(v any) *Filter {
	var f = NewFilter()
	var g = new(generator)
	var fields = g.getFields(v)
	for _, item := range fields {
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
	}
	return f
}

type generator struct {
	kv [2]string
}

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

func (c *generator) getFields(v any) []*field {
	var fields = make([]*field, 0, 5)
	c.doGetFields(reflect.ValueOf(v), &fields)
	return fields
}

func (c *generator) doGetFields(vs reflect.Value, fields *[]*field) {
	if vs.Kind() == reflect.Ptr {
		if vs.IsNil() {
			return
		}
		vs = vs.Elem()
	}

	ts := vs.Type()
	for i := 0; i < vs.NumField(); i++ {
		var tf = ts.Field(i)
		var vf = vs.Field(i)
		var tag = tf.Tag.Get("filter")
		if tag == "-" {
			continue
		}
		if vf.Type().Kind() == reflect.Ptr {
			c.doGetFields(vf, fields)
			continue
		}
		f := c.splitTag(tf.Name, tag)
		f.value = vf.Interface()
		*fields = append(*fields, f)
	}
}

func (c *generator) splitTag(name string, tag string) *field {
	var f = &field{cmp: "eq"}

	c.split(tag, ';', func(item string) {
		item = strings.TrimSpace(item)
		if item == "" {
			return
		}

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
		switch c.kv[0] {
		case "column":
			f.column = c.kv[1]
		case "cmp":
			f.cmp = c.kv[1]
		}
	})

	if f.column == "" {
		f.column = internal.ToSnakeCase(name)
	}
	return f
}
