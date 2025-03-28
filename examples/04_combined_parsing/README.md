# 示例 04: 综合解析

这个示例展示了如何使用JavaScript Package.json Parser库综合解析项目中的`package.json`、`package-lock.json`和`yarn.lock`文件。

## 文件说明

- `main.go` - 示例代码
- `sample_project/` - 包含示例文件的项目目录
  - `package.json` - 示例package.json文件
  - `package-lock.json` - 示例package-lock.json文件
  - `yarn.lock` - 示例yarn.lock文件

## 运行示例

```bash
# 直接运行，无需提供参数，文件路径已在代码中硬编码
go run main.go
```

## 功能演示

这个示例演示了以下功能：

1. 自动检测项目目录中的包管理文件
2. 使用适当的解析器解析每种文件类型
3. 整合并显示解析结果
4. 处理不同文件类型的特定特性

这种综合解析对于以下场景特别有用：
- 分析完整的JavaScript/Node.js项目
- 处理混合使用npm和yarn的项目
- 提取项目的完整依赖信息

## 输出示例

```
=== 解析 package.json ===
项目名称: example-app
项目版本: 1.0.0
模块 'example-app' 包含 8 个依赖
- 常规依赖: 4
- 开发依赖: 4

=== 解析 package-lock.json ===
项目名称: example-app
项目版本: 1.0.0
模块 'example-app' 包含 6 个锁定依赖

=== 解析 yarn.lock ===
项目名称: example-app
项目版本: 未知 (yarn.lock不包含版本信息)
模块 'example-app' 包含 6 个锁定依赖
``` 