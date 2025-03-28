# 示例 02: 解析 package-lock.json

这个示例展示了如何使用JavaScript Package.json Parser库解析`package-lock.json`文件。

## 文件说明

- `main.go` - 示例代码
- `sample_package_lock.json` - 用于测试的示例package-lock.json文件

## 运行示例

```bash
# 直接运行，无需提供参数，文件路径已在代码中硬编码
go run main.go
```

## 功能演示

这个示例演示了以下功能：

1. 初始化PackageLockParser
2. 读取并解析package-lock.json文件
3. 提取项目元数据（名称、版本）
4. 遍历锁定的依赖项
5. 输出依赖的详细信息，包括：
   - 依赖名称和版本
   - 解析地址（resolved）
   - 完整性校验和（integrity）
   - 是否为开发依赖

## 输出示例

```
项目名称: example-app
项目版本: 1.0.0

模块信息:
- 模块: example-app@1.0.0
  依赖数量: 6

  依赖列表 (前6个):
  - axios (常规依赖): 0.21.1
    Resolved: https://registry.npmjs.org/axios/-/axios-0.21.1.tgz
    Integrity: sha512-dKQiRHxGD9PPRIUNIWvZhPTPpl1rf/OxTYKsqKUDjBwYylTvV7SjSHJb9ratfyzM6wCdLCOYLzs73qpg5c4iGA==

  - express (常规依赖): 4.17.1
    Resolved: https://registry.npmjs.org/express/-/express-4.17.1.tgz
    Integrity: sha512-mHJ9O79RqluphRrcw2X/GTh3k9tVv8YcoyY4Kkh4WDMUYKRZUq0h1o0w2rrrxBqM7VoeUVqgb27xlEMXTnYt4g==

  ... [更多依赖] ... 