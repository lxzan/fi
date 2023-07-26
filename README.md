# fi

SQL 条件过滤器

### 快速开始

- 结构体标签生成过滤器 (自动跳过空值)

```go
package main

import (
	"fmt"
	"github.com/lxzan/fi"
)

type Template struct {
	Name   string `filter:"cmp=like"`
	Age    int    `filter:"cmp=lte"`
	Height int
}

func main() {
	var v = &Template{
		Name: "caster",
		Age:  18,
	}
	var filter = fi.GetFilter(v)
	fmt.Printf("%s %v\n", filter.GetExpression(), filter.Args)
}

// `name` LIKE ? AND `age` <= ? [%caster% 18]
```

- 手动构造过滤器 (不自动跳过空值)

```go
package main

import (
	"fmt"
	"github.com/lxzan/fi"
)

type Template struct {
	Name   string `filter:"cmp=like"`
	Age    int    `filter:"cmp=lte"`
	Height int
}

func main() {
	var v = &Template{
		Name: "caster",
		Age:  18,
	}
	var filter = fi.NewFilter().Like("name", v.Name).Lte("age", v.Age).Eq("height", v.Height)
	fmt.Printf("%s %v\n", filter.GetExpression(), filter.Args)
}

// `name` LIKE ? AND `age` <= ? AND `height` = ? [%caster% 18 0]
```

### 标签

```go
type T struct {
    StartTime   int64  `filter:"column=created_at;cmp=gte"`   // 自定义字段名
    Password    string `filter:"-"`                           // 跳过
}
```

