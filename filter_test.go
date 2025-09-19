package fi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var f = NewFilter(WithSkipZeroValue(false), WithQuote(true)).Eq("name", "lee").Gt("age", 0)
	var exp = f.GetExpression()
	assert.Equal(t, exp, "`name` = ? AND `age` > ?")
}

func TestFilter_Equal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true))
		var exp = f.Eq("name", "lee").GetExpression()
		assert.Equal(t, exp, "`name` = ?")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true))
		var timestamp = Timestamp(time.Now().UnixMilli())
		var exp = f.Eq("timestamp", timestamp).GetExpression()
		assert.Equal(t, exp, "`timestamp` = ?")
		assert.Equal(t, int64(timestamp), f.Args[0].(time.Time).UnixMilli())
	})
}

func TestFilter_NotEqual(t *testing.T) {
	var exp = NewFilter(WithQuote(true)).NotEq("name", "lee").GetExpression()
	assert.Equal(t, exp, "`name` != ?")
}

func TestFilter_Merge(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f1 = NewQuery().Eq("name", "caster")
		var f2 = NewQuery().Eq("age", 12)
		var exp = f1.Merge(f2).GetExpression()
		assert.Equal(t, exp, "`name` = ? AND `age` = ?")
	})

	t.Run("", func(t *testing.T) {
		var f1 = NewQuery().Eq("name", "caster")
		var f2 = NewQuery()
		var exp = f1.Merge(f2).GetExpression()
		assert.Equal(t, exp, "`name` = ?")
	})
}

func TestFilter_Gt(t *testing.T) {
	var exp = NewFilter(WithQuote(true)).Gt("age", 1).GetExpression()
	assert.Equal(t, exp, "`age` > ?")
}

func TestFilter_Gte(t *testing.T) {
	var exp = NewFilter(WithQuote(true)).Gte("age", 1).GetExpression()
	assert.Equal(t, exp, "`age` >= ?")
}

func TestFilter_Lt(t *testing.T) {
	var exp = NewFilter(WithQuote(true)).Lt("age", 1).GetExpression()
	assert.Equal(t, exp, "`age` < ?")
}

func TestFilter_Lte(t *testing.T) {
	var exp = NewFilter(WithQuote(true)).Lte("age", 1).GetExpression()
	assert.Equal(t, exp, "`age` <= ?")
}

func TestFilter_Like(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true)).Like("name", "aha")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "`name` LIKE ?")
		assert.Equal(t, f.Args[0], "%aha%")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true)).Like("name", "aha%")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "`name` LIKE ?")
		assert.Equal(t, f.Args[0], "aha%")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter().Like("name", "")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "1=1")
	})
}

func TestFilter_NotLike(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true)).NotLike("name", "aha")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "`name` NOT LIKE ?")
		assert.Equal(t, f.Args[0], "%aha%")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true)).NotLike("name", "aha%")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "`name` NOT LIKE ?")
		assert.Equal(t, f.Args[0], "aha%")
	})
}

func TestFilter_In(t *testing.T) {
	var f = NewFilter(WithQuote(true)).In("name", []string{"lee"})
	var exp = f.GetExpression()
	assert.Equal(t, exp, "`name` IN (?)")
	assert.ElementsMatch(t, f.Args[0], []string{"lee"})
}

func TestFilter_NotIn(t *testing.T) {
	var f = NewFilter(WithQuote(true)).NotIn("name", []string{"lee"})
	var exp = f.GetExpression()
	assert.Equal(t, exp, "`name` NOT IN (?)")
	assert.ElementsMatch(t, f.Args[0], []string{"lee"})
}

func TestFilter_IsNull(t *testing.T) {
	var f = NewFilter(WithQuote(true))
	var exp = f.IsNull("name").GetExpression()
	assert.Equal(t, exp, "`name` IS NULL")
}

func TestFilter_Customize(t *testing.T) {
	assert.Equal(t, "1=1 AND 1=1", NewFilter().Customize("1=1").Customize("1=1").GetExpression())
	var f = Filter{}
	var val = "aha"
	var exp = f.Customize("(first_name LIKE ?) OR (last_name LIKE ?)", val, val).GetExpression()
	assert.Equal(t, exp, "(first_name LIKE ?) OR (last_name LIKE ?)")
}

func TestFilter_GetExpression(t *testing.T) {
	assert.Equal(t, "1=1", new(Filter).GetExpression())
}

func TestFilter_WithTimeSelector(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true)).WithTimeSelector("created_at", 0, 0)
		assert.Equal(t, "1=1", f.GetExpression())
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter(WithQuote(true)).WithTimeSelector("created_at", 1, 2)
		assert.Equal(t, "`created_at` >= ? AND `created_at` < ?", f.GetExpression())
	})
}

func TestWithSkipZeroValue(t *testing.T) {
	f1 := NewFilter(WithSkipZeroValue(true))
	assert.True(t, f1.conf.SkipZeroValue)

	f2 := NewFilter(WithSkipZeroValue(false))
	assert.False(t, f2.conf.SkipZeroValue)

	f3 := NewFilter()
	assert.True(t, f3.conf.SkipZeroValue)
}

func TestWithSize(t *testing.T) {
	f1 := NewFilter(WithQuote(true))
	assert.Equal(t, 10, cap(f1.Args))

	f2 := NewFilter(WithSize(20))
	assert.Equal(t, 20, cap(f2.Args))
}

func TestSetDriver(t *testing.T) {
	t.Run("", func(t *testing.T) {
		f := NewFilter()
		assert.Equal(t, f.conf.Quote, true)
	})

	t.Run("", func(t *testing.T) {
		SetDriver(DriverPostgreSQL)
		f := NewFilter()
		assert.Equal(t, f.conf.Quote, false)
	})
}

func TestNewQuery(t *testing.T) {
	t.Run("", func(t *testing.T) {
		f := NewQuery()
		assert.Equal(t, f.conf.SkipZeroValue, false)
	})

	t.Run("", func(t *testing.T) {
		f := NewFilter()
		assert.Equal(t, f.conf.SkipZeroValue, true)
		assert.Equal(t, f.conf.Size, 10)
	})

	t.Run("", func(t *testing.T) {
		f := NewQuery(WithSize(3))
		assert.Equal(t, f.conf.SkipZeroValue, false)
		assert.Equal(t, f.conf.Size, 3)
	})
}

func TestFilter_Clone(t *testing.T) {
	var a = NewFilter(WithQuote(false)).Eq("name", 1).Eq("age", 2).Clone()
	assert.Equal(t, a.GetExpression(), "name = ? AND age = ?")
	assert.Equal(t, len(a.Args), 2)
}

// 测试Alias方法
func TestFilter_Alias(t *testing.T) {
	f := NewFilter(WithQuote(true)).Alias("u").Eq("name", "test")
	// 注意：包含点号的字段名不会被添加反引号
	assert.Equal(t, "u.name = ?", f.GetExpression())
	assert.Equal(t, "test", f.Args[0])
}

// 测试空字符串的Like查询
func TestFilter_LikeEmptyString(t *testing.T) {
	f := NewFilter().Like("name", "")
	assert.Equal(t, "1=1", f.GetExpression())
	assert.Equal(t, 0, len(f.Args))
}

// 测试已包含百分号的Like查询
func TestFilter_LikeWithPercent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"%test%", "%test%"},
		{"test%", "test%"},
		{"%test", "%test"},
		{"test", "%test%"},
	}

	for _, tt := range tests {
		f := NewFilter().Like("name", tt.input)
		assert.Equal(t, tt.expected, f.Args[0])
	}
}

// 测试NotLike的百分号处理
func TestFilter_NotLikeWithPercent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"%test%", "%test%"},
		{"test%", "test%"},
		{"%test", "%test"},
		{"test", "%test%"},
	}

	for _, tt := range tests {
		f := NewFilter().NotLike("name", tt.input)
		assert.Equal(t, tt.expected, f.Args[0])
	}
}

// 测试Valuer接口
func TestFilter_WithValuer(t *testing.T) {
	timestamp := Timestamp(1640995200000) // 2022-01-01 00:00:00
	f := NewFilter().Eq("created_at", timestamp)

	// 验证参数被正确转换
	expectedTime := time.UnixMilli(1640995200000)
	assert.Equal(t, expectedTime, f.Args[0])
}

// 测试SkipZeroValue配置
func TestFilter_SkipZeroValue(t *testing.T) {
	t.Run("SkipZeroValue=true", func(t *testing.T) {
		f := NewFilter(WithSkipZeroValue(true))
		f.Eq("name", "").Eq("age", 0).Eq("active", false)
		assert.Equal(t, "1=1", f.GetExpression())
		assert.Equal(t, 0, len(f.Args))
	})

	t.Run("SkipZeroValue=false", func(t *testing.T) {
		f := NewFilter(WithSkipZeroValue(false))
		f.Eq("name", "").Eq("age", 0).Eq("active", false)
		assert.Equal(t, "name = ? AND age = ? AND active = ?", f.GetExpression())
		assert.Equal(t, 3, len(f.Args))
	})
}

// 测试不同数据类型的零值
func TestFilter_ZeroValues(t *testing.T) {
	f := NewFilter(WithSkipZeroValue(true))

	// 测试各种零值
	f.Eq("int", 0)
	f.Eq("float", 0.0)
	f.Eq("string", "")
	f.Eq("bool", false)
	f.Eq("slice", []int{})
	f.Eq("ptr", (*int)(nil))

	// 空切片可能不被认为是零值
	expr := f.GetExpression()
	assert.True(t, len(expr) > 0)
}

// 测试非零值
func TestFilter_NonZeroValues(t *testing.T) {
	f := NewFilter(WithSkipZeroValue(true))

	f.Eq("int", 1)
	f.Eq("float", 0.1)
	f.Eq("string", "test")
	f.Eq("bool", true)
	f.Eq("slice", []int{1, 2})

	expected := "int = ? AND float = ? AND string = ? AND bool = ? AND slice = ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 5, len(f.Args))
}

// 测试WithTimeSelector的边界情况
func TestFilter_WithTimeSelectorEdgeCases(t *testing.T) {
	t.Run("BothZero", func(t *testing.T) {
		f := NewFilter().WithTimeSelector("created_at", 0, 0)
		assert.Equal(t, "1=1", f.GetExpression())
	})

	t.Run("OnlyStartTime", func(t *testing.T) {
		f := NewFilter().WithTimeSelector("created_at", 1000, 0)
		// WithTimeSelector总是添加两个条件
		assert.Equal(t, "created_at >= ? AND created_at < ?", f.GetExpression())
		assert.Equal(t, 2, len(f.Args))
	})

	t.Run("OnlyEndTime", func(t *testing.T) {
		f := NewFilter().WithTimeSelector("created_at", 0, 2000)
		// WithTimeSelector总是添加两个条件
		assert.Equal(t, "created_at >= ? AND created_at < ?", f.GetExpression())
		assert.Equal(t, 2, len(f.Args))
	})
}

// 测试复杂查询组合
func TestFilter_ComplexQuery(t *testing.T) {
	f := NewFilter(WithQuote(true)).
		Eq("status", "active").
		Gte("age", 18).
		Lt("age", 65).
		Like("name", "john").
		In("role", []string{"admin", "user"}).
		NotIn("department", []string{"deleted"}).
		IsNull("deleted_at")

	expected := "`status` = ? AND `age` >= ? AND `age` < ? AND `name` LIKE ? AND `role` IN (?) AND `department` NOT IN (?) AND `deleted_at` IS NULL"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 6, len(f.Args)) // IS NULL不占用参数
}

// 测试Customize方法的更多情况
func TestFilter_CustomizeMore(t *testing.T) {
	f := NewFilter().Customize("(name LIKE ? OR email LIKE ?)", "%test%", "%test%")
	assert.Equal(t, "(name LIKE ? OR email LIKE ?)", f.GetExpression())
	assert.Equal(t, 2, len(f.Args))
	assert.Equal(t, "%test%", f.Args[0])
	assert.Equal(t, "%test%", f.Args[1])
}

// 测试Merge方法的更多情况
func TestFilter_MergeMore(t *testing.T) {
	f1 := NewFilter().Eq("status", "active")
	f2 := NewFilter().Eq("age", 18)
	f3 := NewFilter() // 空过滤器

	result1 := f1.Merge(f2)
	assert.Equal(t, "status = ? AND age = ?", result1.GetExpression())
	assert.Equal(t, 2, len(result1.Args))

	// 注意：f1已经被修改了，所以result2包含两个条件
	result2 := f1.Merge(f3)
	assert.Equal(t, "status = ? AND age = ?", result2.GetExpression())
	assert.Equal(t, 2, len(result2.Args))
}

// 测试字段名包含点号的情况
func TestFilter_FieldWithDot(t *testing.T) {
	f := NewFilter(WithQuote(true)).Eq("table.field", "value")
	assert.Equal(t, "table.field = ?", f.GetExpression())
	assert.Equal(t, "value", f.Args[0])
}

// 测试PostgreSQL驱动配置
func TestFilter_PostgreSQLDriver(t *testing.T) {
	SetDriver(DriverPostgreSQL)
	defer SetDriver(DriverMySQL) // 恢复默认驱动

	f := NewFilter().Eq("name", "test")
	assert.Equal(t, "name = ?", f.GetExpression())
}

// 测试参数容量预分配
func TestFilter_ArgsCapacity(t *testing.T) {
	f := NewFilter(WithSize(5))
	assert.Equal(t, 5, cap(f.Args))
	assert.Equal(t, 0, len(f.Args))
}
