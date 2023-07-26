# fi

SQL 条件过滤器

### 快速开始

- 结构体标签反射生成过滤器

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

- 手动构造过滤器

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

// `name` LIKE ? AND `age` <= ? [%caster% 18]
```

### 标签

| 字段   | 描述                                       |
| ------ | ------------------------------------------ |
| column | 自定义的字段名; 默认值是下划线风格的字段名 |
| cmp    | 比较操作符; 默认值是eq                     |
| -      | 忽略                                       |

