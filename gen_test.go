package filter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFilter(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		type Template struct {
			A1  string `filter:"column=a1;op=eq"`
			A2  int    `filter:"column=a2;op=not_eq"`
			A3  int    `filter:"column=a3;op=gt"`
			A4  int    `filter:"column=a4;op=lt"`
			A5  int    `filter:"column=a5;op=gte"`
			A6  int    `filter:"column=a6;op=lte"`
			A7  string `filter:"column=a7;op=like"`
			A8  string `filter:"column=a8;op=not_like"`
			A9  []int  `filter:"column=a9;op=in"`
			A10 []int  `filter:"column=a10;op=not_in"`
		}
		var template = &Template{
			A1:  "1",
			A2:  2,
			A3:  3,
			A4:  4,
			A5:  5,
			A6:  6,
			A7:  "7",
			A8:  "8",
			A9:  []int{9},
			A10: []int{10},
		}
		var f1 = GetFilter(template)
		var f2 = NewFilter().
			Eq("a1", template.A1).
			NotEq("a2", template.A2).
			Gt("a3", template.A3).
			Lt("a4", template.A4).
			Gte("a5", template.A5).
			Lte("a6", template.A6).
			Like("a7", template.A7).
			NotLike("a8", template.A8).
			In("a9", template.A9).
			NotIn("a10", template.A10)
		as.ElementsMatch(f1.Expressions, f2.Expressions)
		as.Equal(len(f1.Args), len(f2.Args))
	})

	t.Run("", func(t *testing.T) {
		type TimeSelector struct {
			StartTime int64 `filter:"column=created_at;op=gte"`
			EndTime   int64 `filter:"column=created_at;op=lt"`
		}

		type Template struct {
			*TimeSelector
			Password string  `filter:"-"`
			Age      int     `filter:" column=age; op=lt ;"`
			Name     *string `filter:"column=name;"`
		}
		var template = Template{
			TimeSelector: &TimeSelector{
				StartTime: 1,
				EndTime:   2,
			},
		}
		var f1 = GetFilter(template)
		var f2 = NewFilter().
			Gte("created_at", template.StartTime).
			Lt("created_at", template.EndTime).
			Lt("age", template.Age).
			Eq("name", template.Name)
		as.ElementsMatch(f1.Expressions, f2.Expressions)
		as.Equal(len(f1.Args), len(f2.Args))
	})
}
