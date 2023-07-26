package filter

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
	var fields = getFields(v)
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

func getFields(v any) []*field {
	var fields []*field
	doGetFields(reflect.ValueOf(v), &fields)
	return fields
}

func doGetFields(vs reflect.Value, fields *[]*field) {
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
			doGetFields(vf, fields)
			continue
		}
		f := splitTag(tf.Name, tag)
		f.value = vf.Interface()
		*fields = append(*fields, f)
	}
}

func splitTag(name string, tag string) *field {
	var f = &field{cmp: "eq"}
	var arr = strings.Split(tag, ";")
	for _, item := range arr {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		var row = strings.SplitN(item, "=", 2)
		row[0] = strings.TrimSpace(row[0])
		row[1] = strings.TrimSpace(row[1])
		switch row[0] {
		case "column":
			f.column = row[1]
		case "cmp":
			f.cmp = row[1]
		}
	}
	if f.column == "" {
		f.column = internal.ToSnakeCase(name)
	}
	return f
}
