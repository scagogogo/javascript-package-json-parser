# 示例 03: 解析 yarn.lock

这个示例展示了如何使用JavaScript Package.json Parser库解析`yarn.lock`文件。

## 文件说明

- `main.go` - 示例代码
- `sample_yarn.lock` - 用于测试的示例yarn.lock文件

## 运行示例

```bash
# 直接运行，无需提供参数，文件路径已在代码中硬编码
go run main.go
```

## 功能演示

这个示例演示了以下功能：

1. 初始化YarnLockParser
2. 读取并解析yarn.lock文件
3. 提取项目信息（名称，yarn.lock通常不包含版本信息）
4. 遍历所有依赖项
5. 输出依赖的详细信息，包括：
   - 依赖名称和版本
   - 解析地址（resolved）
   - 完整性校验和（integrity）
   - 是否有对等依赖（peer dependencies）
   - 是否已打包（bundled）

## 输出示例

```
项目名称: example-app
项目版本: 未知 (yarn.lock不包含版本信息)

模块信息:
- 模块: example-app
  依赖数量: 6

  依赖列表 (前6个):
  - @babel/code-frame: 7.12.13
    Resolved: https://registry.yarnpkg.com/@babel/code-frame/-/code-frame-7.12.13.tgz
    Integrity: sha512-HV1Cm0Q3ZrpCR93tkWOYiuYIgLxZXZFVG2VgK+MBWjUqZTundupbfx2aXarXuw5Ko5aMcjtJgbSs4vUGBS5v6g==

  - @babel/core: 7.13.8
    Resolved: https://registry.yarnpkg.com/@babel/core/-/core-7.13.8.tgz
    Integrity: sha512-oYapIySGw1zGhEFRd6lzWNLWFX2s5dA/jm+Pw/+59ZdXtjyIuwlXbrId22Md0rgZVop+aVoqow2riXhBLNyuQQ==

  - lodash: 4.17.21
    Resolved: https://registry.yarnpkg.com/lodash/-/lodash-4.17.21.tgz
    Integrity: sha512-v2kDEe57lecTulaDIuNTPy3Ry4gLGJ6Z1O3vE1krgXZNrsQ+LFTGHVxVjcXPs17LhbZVGedAJv8XZ1tvj5FvSg==

  ... [更多依赖] ... 