# 示例 01: 解析 package.json

这个示例展示了如何使用JavaScript Package.json Parser库解析`package.json`文件。

## 文件说明

- `main.go` - 示例代码
- `sample_package.json` - 用于测试的示例package.json文件

## 运行示例

```bash
# 直接运行，无需提供参数，文件路径已在代码中硬编码
go run main.go
```

## 功能演示

这个示例演示了以下功能：

1. 初始化PackageJsonParser
2. 读取并解析package.json文件
3. 提取项目元数据（名称、版本）
4. 遍历和分类依赖项（常规依赖和开发依赖）
5. 输出依赖详细信息

## 输出示例

```
项目名称: example-app
项目版本: 1.0.0

模块信息:
- 模块: example-app@1.0.0
  依赖数量: 8
  - 常规依赖: 4
  - 开发依赖: 4

  依赖列表:
  - express (常规依赖): ^4.17.1
  - react (常规依赖): ^17.0.2
  - lodash (常规依赖): ^4.17.21
  - axios (常规依赖): ^0.21.1
  - jest (开发依赖): ^27.0.6
  - typescript (开发依赖): ^4.3.5
  - webpack (开发依赖): ^5.47.0
  - eslint (开发依赖): ^7.32.0
``` 