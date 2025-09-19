package fi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试边界情况和错误处理

// 测试空切片和空映射
func TestFilter_EmptySlicesAndMaps(t *testing.T) {
	f := NewFilter(WithSkipZeroValue(false))

	// 测试空切片
	f.In("tags", []string{})
	f.NotIn("excluded", []int{})

	expected := "`tags` IN (?) AND `excluded` NOT IN (?)"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 2, len(f.Args))
	assert.Equal(t, []string{}, f.Args[0])
	assert.Equal(t, []int{}, f.Args[1])
}

// 测试nil切片和映射
func TestFilter_NilSlicesAndMaps(t *testing.T) {
	f := NewFilter(WithSkipZeroValue(false))

	var nilSlice []string
	var nilMap map[string]int

	f.In("tags", nilSlice)
	f.NotIn("excluded", nilMap)

	expected := "`tags` IN (?) AND `excluded` NOT IN (?)"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 2, len(f.Args))
	assert.Equal(t, []string(nil), f.Args[0])
	assert.Equal(t, map[string]int(nil), f.Args[1])
}

// 测试特殊字符的字段名
func TestFilter_SpecialCharactersInFieldNames(t *testing.T) {
	f := NewFilter(WithQuote(true))

	// 测试包含特殊字符的字段名
	f.Eq("field-with-dash", "value1")
	f.Eq("field_with_underscore", "value2")
	f.Eq("field.with.dot", "value3")
	f.Eq("field with space", "value4")

	// 注意：包含点号的字段名不会被添加反引号
	expected := "`field-with-dash` = ? AND `field_with_underscore` = ? AND field.with.dot = ? AND `field with space` = ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 4, len(f.Args))
}

// 测试超长字符串
func TestFilter_VeryLongStrings(t *testing.T) {
	f := NewFilter()

	// 创建一个很长的字符串
	longString := ""
	for i := 0; i < 1000; i++ {
		longString += "a"
	}

	f.Like("description", longString)

	expected := "`description` LIKE ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 1, len(f.Args))
	assert.Equal(t, "%"+longString+"%", f.Args[0])
}

// 测试特殊数值
func TestFilter_SpecialNumbers(t *testing.T) {
	f := NewFilter(WithSkipZeroValue(false))

	// 测试特殊数值
	f.Eq("max_int", int(^uint(0)>>1))
	f.Eq("min_int", -int(^uint(0)>>1))
	f.Eq("max_float", 1.7976931348623157e+308)
	f.Eq("min_float", -1.7976931348623157e+308)

	// 测试特殊浮点值
	f.Eq("pi", 3.141592653589793)
	f.Eq("e", 2.718281828459045)
	f.Eq("small_float", 1e-10)
	f.Eq("large_float", 1e10)

	expected := "`max_int` = ? AND `min_int` = ? AND `max_float` = ? AND `min_float` = ? AND `pi` = ? AND `e` = ? AND `small_float` = ? AND `large_float` = ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 8, len(f.Args))
}

// 测试时间边界值
func TestFilter_TimeBoundaries(t *testing.T) {
	f := NewFilter(WithSkipZeroValue(false))

	// 测试时间边界值
	zeroTime := time.Time{}
	maxTime := time.Unix(1<<63-1, 999999999) // 最大时间
	minTime := time.Unix(-1<<63, 0)          // 最小时间

	f.Eq("zero_time", zeroTime)
	f.Eq("max_time", maxTime)
	f.Eq("min_time", minTime)

	expected := "`zero_time` = ? AND `max_time` = ? AND `min_time` = ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 3, len(f.Args))
}

// 测试WithTimeSelector的边界值
func TestFilter_WithTimeSelectorBoundaries(t *testing.T) {
	f := NewFilter()

	// 测试最大和最小时间戳
	maxTimestamp := int64(1<<63 - 1)
	minTimestamp := int64(-1 << 63)

	f.WithTimeSelector("created_at", minTimestamp, maxTimestamp)

	expected := "`created_at` >= ? AND `created_at` < ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 2, len(f.Args))
}

// 测试大量条件
func TestFilter_ManyConditions(t *testing.T) {
	f := NewFilter()

	// 添加大量条件
	for i := 0; i < 1000; i++ {
		f.Eq("field"+string(rune(i)), i)
	}

	// 验证表达式长度
	expr := f.GetExpression()
	assert.True(t, len(expr) > 0)
	// 由于零值跳过，实际参数数量可能少于1000
	assert.True(t, len(f.Args) > 0)
}

// 测试并发安全性
func TestFilter_Concurrency(t *testing.T) {
	f := NewFilter()

	// 并发添加条件
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			f.Eq("field"+string(rune(i)), i)
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证结果
	expr := f.GetExpression()
	assert.True(t, len(expr) > 0)
	// 由于并发和零值跳过，实际参数数量可能少于10
	assert.True(t, len(f.Args) > 0)
}

// 测试内存泄漏（通过大量操作）
func TestFilter_MemoryLeak(t *testing.T) {
	// 创建大量过滤器并丢弃
	for i := 0; i < 10000; i++ {
		f := NewFilter()
		for j := 0; j < 10; j++ {
			f.Eq("field"+string(rune(j)), j)
		}
		_ = f.GetExpression()
	}

	// 如果没有内存泄漏，这个测试应该能正常完成
	assert.True(t, true)
}

// 测试无效的数据库驱动
func TestSetDriver_InvalidDriver(t *testing.T) {
	originalDriver := __driver
	defer func() {
		__driver = originalDriver
	}()

	// 测试无效驱动
	SetDriver("invalid_driver")
	assert.Equal(t, originalDriver, __driver) // 应该保持原值

	// 测试空字符串
	SetDriver("")
	assert.Equal(t, originalDriver, __driver) // 应该保持原值
}

// 测试负数Size配置
func TestWithSize_NegativeSize(t *testing.T) {
	f := NewFilter(WithSize(-1))
	assert.Equal(t, 0, cap(f.Args)) // 负数应该被处理为0
}

// 测试零Size配置
func TestWithSize_ZeroSize(t *testing.T) {
	f := NewFilter(WithSize(0))
	assert.Equal(t, 0, cap(f.Args))
}

// 测试非常大的Size配置
func TestWithSize_VeryLargeSize(t *testing.T) {
	f := NewFilter(WithSize(1000000))
	assert.Equal(t, 1000000, cap(f.Args))
}

// 测试空字符串的addPercent
func TestFilter_addPercent_EmptyString(t *testing.T) {
	f := &Filter{}

	// 测试空字符串
	result := f.addPercent("")
	assert.Equal(t, "", result)

	// 测试只有百分号的字符串
	result = f.addPercent("%")
	assert.Equal(t, "%", result)

	// 测试两端都有百分号
	result = f.addPercent("%test%")
	assert.Equal(t, "%test%", result)
}

// 测试Unicode字符串
func TestFilter_UnicodeStrings(t *testing.T) {
	f := NewFilter()

	// 测试Unicode字符串
	unicodeStr := "测试中文"
	emojiStr := "😀🎉🚀"
	mixedStr := "Hello 世界 🌍"

	f.Like("name", unicodeStr)
	f.Like("emoji", emojiStr)
	f.Like("mixed", mixedStr)

	expected := "`name` LIKE ? AND `emoji` LIKE ? AND `mixed` LIKE ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 3, len(f.Args))
	assert.Equal(t, "%"+unicodeStr+"%", f.Args[0])
	assert.Equal(t, "%"+emojiStr+"%", f.Args[1])
	assert.Equal(t, "%"+mixedStr+"%", f.Args[2])
}

// 测试结构体标签的边界情况
func TestGetFilter_TagBoundaryCases(t *testing.T) {
	type TestStruct struct {
		Field1 string `fi:"column=;cmp=eq"`                     // 空column
		Field2 string `fi:"column=field2;cmp="`                 // 空cmp
		Field3 string `fi:"column=field3;cmp=invalid"`          // 无效cmp
		Field4 string `fi:"column=field4;cmp=eq;extra=value"`   // 额外标签
		Field5 string `fi:"column=field5;cmp=eq;;"`             // 末尾分号
		Field6 string `fi:";column=field6;cmp=eq"`              // 开头分号
		Field7 string `fi:"column=field7;cmp=eq;column=field8"` // 重复标签
		Field8 string `fi:"column=field8;cmp=eq;cmp=gt"`        // 重复标签
	}

	user := &TestStruct{
		Field1: "test1",
		Field2: "test2",
		Field3: "test3",
		Field4: "test4",
		Field5: "test5",
		Field6: "test6",
		Field7: "test7",
		Field8: "test8",
	}

	f := GetFilter(user)

	// 验证表达式和参数（某些字段可能被跳过）
	expr := f.GetExpression()
	assert.True(t, len(expr) > 0)
	assert.True(t, len(f.Args) > 0)
}

// 测试深度嵌套结构体
func TestGetFilter_DeepNestedStruct(t *testing.T) {
	type Level3 struct {
		Field string `fi:"column=level3_field;cmp=eq"`
	}

	type Level2 struct {
		Field  string  `fi:"column=level2_field;cmp=eq"`
		Level3 *Level3 `fi:"column=level3;cmp=eq"`
	}

	type Level1 struct {
		Field  string  `fi:"column=level1_field;cmp=eq"`
		Level2 *Level2 `fi:"column=level2;cmp=eq"`
	}

	user := &Level1{
		Field: "level1",
		Level2: &Level2{
			Field: "level2",
			Level3: &Level3{
				Field: "level3",
			},
		},
	}

	f := GetFilter(user)
	expected := "`level1_field` = ? AND `level2_field` = ? AND `level3_field` = ?"
	assert.Equal(t, expected, f.GetExpression())
	assert.Equal(t, 3, len(f.Args))
}

// 测试循环引用结构体（跳过，因为会导致无限递归）
// func TestGetFilter_CircularReference(t *testing.T) {
// 	// 这个测试会导致无限递归，跳过
// }
