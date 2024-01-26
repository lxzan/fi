# fi

动态 SQL 查询条件构造器

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

- 反射生成过滤器

```go
package main

import (
	"fmt"
	"github.com/lxzan/fi"
)

type Template struct {
	Name   string `fi:"cmp=like"`
	Age    int    `fi:"cmp=lte"`
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
	Name   string `fi:"cmp=like"`
	Age    int    `fi:"cmp=lte"`
	Height int
}

func main() {
	var v = &Template{
		Name: "caster",
		Age:  18,
	}
	var filter = fi.
		NewFilter().
		Like("name", v.Name).
		Lte("age", v.Age).
		Eq("height", v.Height)
	fmt.Printf("%s %v\n", filter.GetExpression(), filter.Args)
}

// `name` LIKE ? AND `age` <= ? [%caster% 18]
```

### 标签

| 字段     | 描述                |
|--------|-------------------|
| column | 自定义字段名; 默认值是下划线风格 |
| cmp    | 比较操作符; 默认值是eq     |
| -      | 忽略                |

### 操作符

| 操作符      | 描述         |
|----------|------------|
| eq       | `=`        |
| not_eq   | `!=`       |
| gt       | `>`        |
| lt       | `<`        |
| gte      | `>=`       |
| lte      | `<=`       |
| like     | `LIKE`     |
| not_like | `NOT LIKE` |
| in       | `IN`       |
| not_in   | `NOT IN`   |

### 性能测试

```
go test -benchmem -run=^$ -bench ^Benchmark github.com/lxzan/fi
goos: windows
goarch: amd64
pkg: github.com/lxzan/fi
cpu: 13th Gen Intel(R) Core(TM) i5-1340P
BenchmarkGetFilterReflect-16              487071              2249 ns/op             768 B/op         29 allocs/op
BenchmarkGetFilterNoReflect-16           2214118               569.0 ns/op           558 B/op         11 allocs/op
PASS
ok      github.com/lxzan/fi     3.318s
```