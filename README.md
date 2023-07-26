# fi

SQL 条件过滤器

[![Build Status][1]][2] [![MIT licensed][3]][4] [![Go Version][5]][6] [![codecov][7]][8] [![Go Report Card][9]][10]

[1]: https://github.com/lxzan/fi/workflows/Go%20Test/badge.svg?branch=main

[2]: https://github.com/lxzan/fi/actions?query=branch%3Amain

[3]: https://img.shields.io/badge/license-MIT-blue.svg

[4]: LICENSE

[5]: https://img.shields.io/badge/go-%3E%3D1.18-30dff3?style=flat-square&logo=go

[6]: https://github.com/lxzan/fi

[7]: https://codecov.io/github/lxzan/fi/branch/main/graph/badge.svg?token=uMNEU3cEsM

[8]: https://app.codecov.io/gh/lxzan/fi

[9]: https://goreportcard.com/badge/github.com/lxzan/fi

[10]: https://goreportcard.com/report/github.com/lxzan/fi

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

| 字段     | 描述                     |
|--------|------------------------|
| column | 自定义的字段名; 默认值是下划线风格的字段名 |
| cmp    | 比较操作符; 默认值是eq          |
| -      | 忽略                     |

### 性能测试

```
go test -benchmem -run=^$ -bench ^Benchmark github.com/lxzan/fi
goos: darwin
goarch: arm64
pkg: github.com/lxzan/fi
BenchmarkGetFilterReflect-8               445311              2671 ns/op            1928 B/op         48 allocs/op
BenchmarkGetFilterNoReflect-8            1735392               692.3 ns/op          1104 B/op         18 allocs/op
PASS
ok      github.com/lxzan/fi     4.330s
```