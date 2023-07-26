package filter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	var f = NewFilter().Equal("name", "lee").Gt("age", 0)
	var exp = f.GetExpression()
	assert.Equal(t, exp, "name = ? AND age > ?")
}

func TestFilter_Equal(t *testing.T) {
	var f = Filter{}
	var exp = f.Equal("name", "lee").GetExpression()
	assert.Equal(t, exp, "name = ?")
}

func TestFilter_NotEqual(t *testing.T) {
	var exp = NewFilter().NotEqual("name", "lee").GetExpression()
	assert.Equal(t, exp, "name != ?")
}

func TestFilter_Gt(t *testing.T) {
	var exp = NewFilter().Gt("age", 1).GetExpression()
	assert.Equal(t, exp, "age > ?")
}

func TestFilter_Gte(t *testing.T) {
	var exp = NewFilter().Gte("age", 1).GetExpression()
	assert.Equal(t, exp, "age >= ?")
}

func TestFilter_Lt(t *testing.T) {
	var exp = NewFilter().Lt("age", 1).GetExpression()
	assert.Equal(t, exp, "age < ?")
}

func TestFilter_Lte(t *testing.T) {
	var exp = NewFilter().Lte("age", 1).GetExpression()
	assert.Equal(t, exp, "age <= ?")
}

func TestFilter_Like(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f = NewFilter().Like("name", "aha")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "name LIKE ?")
		assert.Equal(t, f.Args[0], "%aha%")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter().Like("name", "aha%")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "name LIKE ?")
		assert.Equal(t, f.Args[0], "aha%")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter().Like("name", "")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "name LIKE ?")
	})
}

func TestFilter_NotLike(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var f = NewFilter().NotLike("name", "aha")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "name NOT LIKE ?")
		assert.Equal(t, f.Args[0], "%aha%")
	})

	t.Run("", func(t *testing.T) {
		var f = NewFilter().NotLike("name", "aha%")
		var exp = f.GetExpression()
		assert.Equal(t, exp, "name NOT LIKE ?")
		assert.Equal(t, f.Args[0], "aha%")
	})
}

func TestFilter_In(t *testing.T) {
	var f = NewFilter().In("name", "lee")
	var exp = f.GetExpression()
	assert.Equal(t, exp, "name IN (?)")
	assert.ElementsMatch(t, f.Args[0], []string{"lee"})
}

func TestFilter_NotIn(t *testing.T) {
	var f = NewFilter().NotIn("name", "lee")
	var exp = f.GetExpression()
	assert.Equal(t, exp, "name NOT IN (?)")
	assert.ElementsMatch(t, f.Args[0], []string{"lee"})
}

func TestFilter_IsNull(t *testing.T) {
	var f = Filter{}
	var exp = f.IsNull("name").GetExpression()
	assert.Equal(t, exp, "name IS NULL")
}

func TestFilter_With(t *testing.T) {
	var f = Filter{}
	var val = "aha"
	var exp = f.Customize("(first_name LIKE ?) OR (last_name LIKE ?)", val, val).GetExpression()
	assert.Equal(t, exp, "(first_name LIKE ?) OR (last_name LIKE ?)")
}
