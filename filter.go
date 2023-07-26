package fi

import (
	"github.com/lxzan/fi/internal"
	"strings"
)

type Filter struct {
	Skip        bool          // 是否跳过空值
	Expressions []string      // 表达式
	Args        []interface{} // 参数
}

func NewFilter() *Filter {
	return new(Filter)
}

func (c *Filter) push(key string, val any, cmp string) *Filter {
	if internal.IsNil(val) || (c.Skip && internal.IsZero(val)) {
		return c
	}

	var hasDot = strings.Contains(key, ".")
	builder := strings.Builder{}
	if !hasDot {
		builder.WriteString("`")
	}
	builder.WriteString(key)
	if !hasDot {
		builder.WriteString("`")
	}
	builder.WriteString(" ")
	builder.WriteString(cmp)
	switch cmp {
	case "IN", "NOT IN":
		builder.WriteString(" (?)")
		c.Args = append(c.Args, val)
	case "IS NULL":
	default:
		builder.WriteString(" ?")
		c.Args = append(c.Args, val)
	}
	c.Expressions = append(c.Expressions, builder.String())
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
func (c *Filter) Customize(key string, val ...any) *Filter {
	c.Expressions = append(c.Expressions, key)
	c.Args = append(c.Args, val...)
	return c
}

// GetExpression 获取SQL表达式
func (c *Filter) GetExpression() string {
	if len(c.Expressions) == 0 {
		return "1=1"
	}
	return strings.Join(c.Expressions, " AND ")
}
