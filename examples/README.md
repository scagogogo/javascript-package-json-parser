# JavaScript Package.json Parser - Examples

这个目录包含了几个如何使用JavaScript Package.json Parser的示例。

## 示例概述

1. **[01_parse_package_json](./01_parse_package_json)** - 解析package.json文件的基本示例
2. **[02_parse_package_lock](./02_parse_package_lock)** - 解析package-lock.json文件的示例
3. **[03_parse_yarn_lock](./03_parse_yarn_lock)** - 解析yarn.lock文件的示例
4. **[04_combined_parsing](./04_combined_parsing)** - 综合示例：解析一个项目中的package.json、package-lock.json和yarn.lock

## 运行示例

每个示例都已经配置为直接运行，不需要提供命令行参数。只需进入对应的目录并执行：

### 1. 解析package.json

```bash
cd 01_parse_package_json
go run main.go
```

### 2. 解析package-lock.json

```bash
cd 02_parse_package_lock
go run main.go
```

### 3. 解析yarn.lock

```bash
cd 03_parse_yarn_lock
go run main.go
```

### 4. 综合解析

```bash
cd 04_combined_parsing
go run main.go
```

## 示例输出

运行这些示例将会输出所解析文件的结构信息，包括：

- 项目名称和版本
- 模块信息
- 依赖项列表和详细信息
- 依赖类型（常规依赖或开发依赖）
- 其他元数据（如有） 