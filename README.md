# fi

[![Build Status][1]][2] [![MIT licensed][3]][4] [![Go Version][5]][6] [![codecov][7]][8] [![Go Report Card][9]][10]

[1]: https://github.com/lxzan/fi/actions/workflows/go.yml/badge.svg
[2]: https://github.com/lxzan/fi/actions/workflows/go.yml
[3]: https://img.shields.io/badge/license-MIT-blue.svg
[4]: LICENSE
[5]: https://img.shields.io/badge/go-%3E%3D1.18-30dff3?style=flat-square&logo=go
[6]: https://github.com/lxzan/fi
[7]: https://codecov.io/github/lxzan/fi/branch/main/graph/badge.svg?token=uMNEU3cEsM
[8]: https://app.codecov.io/gh/lxzan/fi
[9]: https://goreportcard.com/badge/github.com/lxzan/fi
[10]: https://goreportcard.com/report/github.com/lxzan/fi

## 简介

`fi` 是一个高性能的 Go 语言动态 SQL 查询条件构造器，支持通过反射自动生成查询条件，也支持手动构建复杂的 SQL WHERE 子句。它提供了简洁的 API 和丰富的功能，让数据库查询变得更加简单和高效。

### ✨ 特性

- 🚀 **高性能**: 优化的反射机制，支持预分配内存
- 🔧 **灵活配置**: 支持多种数据库驱动（MySQL、PostgreSQL）
- 🏷️ **智能标签**: 通过结构体标签自动生成查询条件
- 🔗 **链式调用**: 流畅的 API 设计，支持方法链
- 🛡️ **类型安全**: 完整的类型检查和零值处理
- 📦 **轻量级**: 无外部依赖，仅依赖标准库
- 🧪 **完整测试**: 100% 测试覆盖率，包含边界情况测试

## 快速开始

### 安装

```bash
go get github.com/lxzan/fi
```

### 基本用法

#### 1. 反射生成过滤器

通过结构体标签自动生成查询条件，让代码更加简洁：

```go
package main

import (
    "fmt"
    "github.com/lxzan/fi"
)

type User struct {
    Name   string `fi:"cmp=like"`           // 模糊查询
    Age    int    `fi:"cmp=lte"`            // 小于等于
    Status string `fi:"column=user_status"` // 自定义字段名
    Height int                              // 无标签，使用默认配置
}

func main() {
    user := &User{
        Name:   "张三",
        Age:    25,
        Status: "active",
        Height: 175,
    }
    
    filter := fi.GetFilter(user)
    fmt.Printf("SQL: %s\n", filter.GetExpression())
    fmt.Printf("参数: %v\n", filter.Args)
}

// 输出:
// SQL: `name` LIKE ? AND `age` <= ? AND `user_status` = ? AND `height` = ?
// 参数: [%张三% 25 active 175]
```

#### 2. 手动构造过滤器

对于复杂查询，可以手动构建查询条件：

```go
package main

import (
    "fmt"
    "github.com/lxzan/fi"
)

func main() {
    filter := fi.NewFilter().
        Like("name", "张三").
        Lte("age", 25).
        Eq("status", "active").
        Gte("created_at", "2023-01-01").
        In("role", []string{"admin", "user"}).
        NotNull("deleted_at")
    
    fmt.Printf("SQL: %s\n", filter.GetExpression())
    fmt.Printf("参数: %v\n", filter.Args)
}

// 输出:
// SQL: `name` LIKE ? AND `age` <= ? AND `status` = ? AND `created_at` >= ? AND `role` IN (?) AND `deleted_at` IS NOT NULL
// 参数: [%张三% 25 active 2023-01-01 [admin user]]
```

## 配置选项

### 结构体标签

| 标签    | 描述                    | 默认值        |
|---------|------------------------|---------------|
| `column`| 自定义字段名             | 下划线风格转换 |
| `cmp`   | 比较操作符               | `eq`          |
| `-`     | 忽略该字段               | -             |

#### 标签语法

标签支持分号分隔的多个选项：

```go
type User struct {
    // 基本用法
    Name string `fi:"cmp=like"`
    
    // 自定义字段名
    Email string `fi:"column=user_email;cmp=eq"`
    
    // 多个选项用分号分隔
    Age int `fi:"column=user_age;cmp=gte"`
    
    // 忽略字段
    Password string `fi:"-"`
    
    // 空标签使用默认配置
    Status string `fi:""`
}
```

### 比较操作符

| 操作符      | 描述         | SQL 等价        |
|-------------|-------------|-----------------|
| `eq`        | 等于         | `=`             |
| `not_eq`    | 不等于       | `!=`            |
| `gt`        | 大于         | `>`             |
| `lt`        | 小于         | `<`             |
| `gte`       | 大于等于     | `>=`            |
| `lte`       | 小于等于     | `<=`            |
| `like`      | 模糊查询     | `LIKE`          |
| `not_like`  | 非模糊查询   | `NOT LIKE`      |
| `in`        | 包含         | `IN`            |
| `not_in`    | 不包含       | `NOT IN`        |
| `is_null`   | 为空         | `IS NULL`       |
| `not_null`  | 不为空       | `IS NOT NULL`   |

## 高级用法

### 自定义配置

```go
// 创建自定义配置的过滤器
filter := fi.NewFilter(
    fi.WithSkipZeroValue(true),  // 跳过零值
    fi.WithSize(20),             // 预分配参数容量
    fi.WithQuote(false),         // 不添加反引号
)

// 设置数据库驱动
fi.SetDriver(fi.DriverPostgreSQL) // 或 fi.DriverMySQL
```

### 时间范围查询

```go
// 使用时间选择器
filter := fi.NewFilter().
    WithTimeSelector("created_at", 1640995200000, 1672531200000)

// 生成: created_at >= ? AND created_at < ?
```

### 自定义 SQL

```go
// 添加自定义 SQL 片段
filter := fi.NewFilter().
    Customize("EXISTS (SELECT 1 FROM orders WHERE user_id = ?)", 123)

// 生成: EXISTS (SELECT 1 FROM orders WHERE user_id = ?)
```

### 过滤器合并

```go
filter1 := fi.NewFilter().Eq("status", "active")
filter2 := fi.NewFilter().Gte("age", 18)

// 合并过滤器
result := filter1.Merge(filter2)
// 生成: status = ? AND age >= ?
```

## 性能测试

在标准测试环境下的性能表现：

```bash
go test -benchmem -run=^$ -bench ^Benchmark github.com/lxzan/fi
```

```
goos: windows
goarch: amd64
pkg: github.com/lxzan/fi
cpu: 13th Gen Intel(R) Core(TM) i5-1340P
BenchmarkGetFilterReflect-16              487071              2249 ns/op             768 B/op         29 allocs/op
BenchmarkGetFilterNoReflect-16           2214118               569.0 ns/op           558 B/op         11 allocs/op
PASS
ok      github.com/lxzan/fi     3.318s
```

## 最佳实践

### 1. 合理使用反射

- 对于简单查询，推荐使用反射方式，代码更简洁
- 对于复杂查询，建议使用手动构建方式，性能更好

### 2. 预分配容量

```go
// 预估查询条件数量，预分配内存
filter := fi.NewFilter(fi.WithSize(10))
```

### 3. 零值处理

```go
// 根据业务需求决定是否跳过零值
filter := fi.NewFilter(fi.WithSkipZeroValue(true))
```

### 4. 数据库驱动配置

```go
// 根据使用的数据库设置驱动类型
fi.SetDriver(fi.DriverMySQL)     // MySQL
fi.SetDriver(fi.DriverPostgreSQL) // PostgreSQL
```

## 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目！
