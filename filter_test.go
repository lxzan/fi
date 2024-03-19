package fi

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
