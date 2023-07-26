package filter

import "strings"

type Filter struct {
	Expressions []string
	Args        []interface{}
}

func NewFilter() *Filter {
	return new(Filter)
}

func (c *Filter) push(key string, val any, op string) *Filter {
	if val == nil {
		return c
	}

	builder := strings.Builder{}
	builder.WriteString(key)
	builder.WriteString(" ")
	builder.WriteString(op)
	switch op {
	case "IN", "NOT IN":
		builder.WriteString(" (?)")
	default:
		builder.WriteString(" ?")
	}
	c.Expressions = append(c.Expressions, builder.String())
	c.Args = append(c.Args, val)
	return c
}

func (c *Filter) Equal(key string, val any) *Filter {
	return c.push(key, val, "=")
}

func (c *Filter) NotEqual(key string, val any) *Filter {
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

func (c *Filter) Like(key string, val string) *Filter {
	return c.push(key, c.addPercent(val), "LIKE")
}

func (c *Filter) NotLike(key string, val string) *Filter {
	return c.push(key, c.addPercent(val), "NOT LIKE")
}

func (c *Filter) In(key string, val ...any) *Filter {
	return c.push(key, val, "IN")
}

func (c *Filter) NotIn(key string, val ...any) *Filter {
	return c.push(key, val, "NOT IN")
}

func (c *Filter) IsNull(key string) *Filter {
	c.Expressions = append(c.Expressions, key+" IS NULL")
	return c
}

func (c *Filter) Customize(key string, val ...any) *Filter {
	c.Expressions = append(c.Expressions, key)
	c.Args = append(c.Args, val...)
	return c
}

func (c *Filter) GetExpression() string {
	return strings.Join(c.Expressions, " AND ")
}
