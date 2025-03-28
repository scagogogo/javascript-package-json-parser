# JavaScript Package.json Parser

[![Build Status](https://github.com/scagogogo/javascript-package-json-parser/actions/workflows/test.yml/badge.svg)](https://github.com/scagogogo/javascript-package-json-parser/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/javascript-package-json-parser)](https://goreportcard.com/report/github.com/scagogogo/javascript-package-json-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/javascript-package-json-parser?status.svg)](https://godoc.org/github.com/scagogogo/javascript-package-json-parser)
[![codecov](https://codecov.io/gh/scagogogo/javascript-package-json-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/javascript-package-json-parser)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个高效、灵活的 Go 语言库，用于解析 JavaScript 项目中的 `package.json`、`package-lock.json` 和 `yarn.lock` 文件。该库将这些文件解析为结构化的对象，便于进行依赖分析、安全扫描和生态系统研究。

## 目录

- [功能特点](#功能特点)
- [安装](#安装)
- [快速开始](#快速开始)
- [使用示例](#使用示例)
  - [解析 package.json](#解析-packagejson)
  - [解析 package-lock.json](#解析-package-lockjson)
  - [解析 yarn.lock](#解析-yarnlock)
  - [内存中的 JSON 解析](#内存中的-json-解析)
- [API 文档](#api-文档)
- [数据模型](#数据模型)
- [示例代码](#示例代码)
- [版本兼容性](#版本兼容性)
- [持续集成](#持续集成)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## 功能特点

- ✅ **多格式支持** - 解析 `package.json`、`package-lock.json` 和 `yarn.lock` 文件
- ✅ **结构化数据** - 提取项目元数据、依赖信息、版本约束等
- ✅ **生态系统兼容** - 支持 npm 和 yarn 生态系统
- ✅ **强类型模型** - 提供结构化的数据模型，使用 Go 泛型
- ✅ **高性能** - 高效的文件解析和内存管理
- ✅ **内存解析** - 支持解析内存中的 JSON 字符串
- ✅ **完整测试** - 高测试覆盖率保证代码质量
- ✅ **详细文档** - 全面的代码注释和使用示例

## 安装

### 要求

- Go 1.18+ (使用了泛型特性)

### 使用 Go Modules 安装

```bash
go get -u github.com/scagogogo/package-json-parser
```

### 在你的 Go 项目中导入

```go
import "github.com/scagogogo/package-json-parser/pkg/parser"
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/package-json-parser/pkg/parser"
)

func main() {
    // 创建解析器实例
    packageJsonParser := &parser.PackageJsonParser{}
    
    // 初始化并解析
    ctx := context.Background()
    err := packageJsonParser.Init(ctx)
    if err != nil {
        log.Fatalf("初始化失败: %v", err)
    }
    defer packageJsonParser.Close(ctx)
    
    // 解析文件
    input := &parser.PackageJsonParserInput{
        PackageJsonPath: "./package.json",
    }
    project, err := packageJsonParser.Parse(ctx, input)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 使用解析结果
    fmt.Printf("项目: %s@%s\n", project.Name, project.Version)
    fmt.Printf("依赖数量: %d\n", len(project.Modules[0].Dependencies))
}
```

## 使用示例

### 解析 package.json

```go
packageJsonParser := &parser.PackageJsonParser{}
err := packageJsonParser.Init(context.Background())
if err != nil {
    panic(err)
}
defer packageJsonParser.Close(context.Background())

input := &parser.PackageJsonParserInput{
    PackageJsonPath: "./path/to/package.json",
}

project, err := packageJsonParser.Parse(context.Background(), input)
if err != nil {
    panic(err)
}

// 使用解析结果
fmt.Printf("项目名称: %s\n", project.Name)
fmt.Printf("项目版本: %s\n", project.Version)

// 获取模块(默认使用项目名称作为模块名)
module := project.GetModule(project.Name)
if module != nil {
    fmt.Printf("依赖数量: %d\n", len(module.Dependencies))
    
    // 遍历依赖
    for _, dep := range module.Dependencies {
        fmt.Printf("依赖: %s@%s\n", dep.DependencyName, dep.DependencyVersion)
    }
}
```

### 解析 package-lock.json

```go
// 创建package-lock.json解析器
packageLockParser := parser.NewPackageLockParser()

// 初始化解析器
err := packageLockParser.Init(context.Background())
if err != nil {
    panic(err)
}
defer packageLockParser.Close(context.Background())

// 创建解析输入
input := &parser.PackageLockJsonParserInput{
    PackageLockJsonPath: "./path/to/package-lock.json",
}

// 解析文件
project, err := packageLockParser.Parse(context.Background(), input)
if err != nil {
    panic(err)
}

// 使用解析结果
fmt.Printf("项目名称: %s\n", project.Name)
fmt.Printf("项目版本: %s\n", project.Version)

// 获取模块信息
module := project.GetModule(project.Name)
if module != nil {
    fmt.Printf("依赖数量: %d\n", len(module.Dependencies))
    
    // 遍历依赖
    for _, dep := range module.Dependencies {
        fmt.Printf("依赖: %s@%s\n", dep.DependencyName, dep.DependencyVersion)
        if dep.ComponentDependencyEcosystem != nil {
            fmt.Printf("  - Resolved: %s\n", dep.ComponentDependencyEcosystem.Resolved)
            fmt.Printf("  - Integrity: %s\n", dep.ComponentDependencyEcosystem.Integrity)
        }
    }
}
```

### 解析 yarn.lock

```go
// 创建yarn.lock解析器
yarnLockParser := parser.NewYarnLockParser()

// 初始化解析器
err := yarnLockParser.Init(context.Background())
if err != nil {
    panic(err)
}
defer yarnLockParser.Close(context.Background())

// 创建解析输入
input := &parser.YarnLockParserInput{
    YarnLockPath: "./path/to/yarn.lock",
}

// 解析文件
project, err := yarnLockParser.Parse(context.Background(), input)
if err != nil {
    panic(err)
}

// 使用解析结果
fmt.Printf("项目名称: %s\n", project.Name)
if project.Version != "" {
    fmt.Printf("项目版本: %s\n", project.Version)
}

// 获取依赖信息
for _, module := range project.Modules {
    fmt.Printf("模块: %s\n", module.Name)
    fmt.Printf("依赖数量: %d\n", len(module.Dependencies))
    
    // 遍历依赖
    for _, dep := range module.Dependencies {
        fmt.Printf("依赖: %s@%s\n", dep.DependencyName, dep.DependencyVersion)
    }
}
```

### 内存中的 JSON 解析

```go
// 直接从字符串解析 package.json 内容
jsonContent := `{
    "name": "example-project",
    "version": "1.0.0",
    "dependencies": {
        "express": "^4.17.1"
    }
}`

packageJsonParser := &parser.PackageJsonParser{}
err := packageJsonParser.Init(context.Background())
if err != nil {
    panic(err)
}
defer packageJsonParser.Close(context.Background())

// 使用内容直接解析，无需文件路径
input := &parser.PackageJsonParserInput{
    PackageJsonContent: jsonContent,
}

project, err := packageJsonParser.Parse(context.Background(), input)
if err != nil {
    panic(err)
}

fmt.Printf("项目名称: %s\n", project.Name)
```

## API 文档

详细的 API 文档可以在 [GoDoc](https://godoc.org/github.com/scagogogo/package-json-parser) 上找到。

主要接口和类型包括:

- **解析器接口**: 所有解析器都实现了通用接口
- **输入类型**: 每种解析器有对应的输入结构体
- **项目模型**: 表示整个项目结构的通用模型
- **模块和依赖**: 表示模块和依赖关系的模型
- **生态系统特定信息**: npm 和 yarn 生态系统特有的信息

## 数据模型

该项目定义了详细的数据模型来表示各种包管理文件的内容：

- `PackageJson`：表示 package.json 文件的完整结构
- `PackageLock`：表示 package-lock.json 文件的结构
- `YarnLock`: 表示 yarn.lock 文件的结构
- `Dependencies`：表示依赖映射
- 各种生态系统特定的模型，如 `PackageLockComponentEcosystem` 和 `YarnLockComponentDependencyEcosystem` 等

## 示例代码

本项目包含一系列演示如何使用解析器库的示例代码。详细信息请参阅 [examples 目录](./examples/)。

每个示例都已经配置为直接运行，不需要提供命令行参数：

```bash
# 解析 package.json 示例
cd examples/01_parse_package_json
go run main.go

# 解析 package-lock.json 示例
cd examples/02_parse_package_lock
go run main.go

# 解析 yarn.lock 示例
cd examples/03_parse_yarn_lock
go run main.go

# 综合解析示例
cd examples/04_combined_parsing
go run main.go
```

## 版本兼容性

该解析器兼容以下版本的文件格式：

- npm v5-v7 的 package-lock.json 格式
- 标准的 package.json 格式
- yarn.lock v1 格式

## 持续集成

本项目使用 GitHub Actions 实现自动化测试：

- **单元测试**：所有包中的测试都会自动运行
- **代码覆盖率**：测试结果会生成覆盖率报告
- **性能基准测试**：确保性能稳定
- **代码质量**：使用 golangci-lint 进行代码质量检查
- **示例验证**：所有示例代码都会被执行以确保它们能正常工作

每次代码提交或 PR 都会触发这些自动化测试。

更多信息，请参阅：
- [GitHub Actions 工作流程说明](./docs/GITHUB_ACTIONS.md)
- [贡献指南](./docs/CONTRIBUTING.md)

## 贡献指南

欢迎提交问题和拉取请求！

贡献前，请确保：

1. 添加相应的单元测试
2. 更新文档
3. 确保所有测试通过
4. 确保代码符合 Go 的代码规范

详细的贡献流程请参见 [贡献指南](./docs/CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](./LICENSE) 文件了解详情。 