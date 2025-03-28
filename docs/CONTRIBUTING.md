# 贡献指南

感谢你考虑为JavaScript Package.json Parser项目做出贡献。以下是一些指导方针，可以帮助你顺利地参与到项目中来。

## 提交流程

1. Fork本项目到你的GitHub账户下
2. 克隆你的fork到本地：`git clone https://github.com/YOUR-USERNAME/javascript-package-json-parser.git`
3. 创建一个新分支：`git checkout -b feature/your-feature-name`
4. 在这个分支上进行修改
5. 提交你的更改：`git commit -am '添加新功能：简短描述'`
6. 推送到你的fork: `git push origin feature/your-feature-name`
7. 创建一个Pull Request

## 代码规范

- 遵循Go的标准代码风格
- 使用有意义的变量名和函数名
- 添加适当的注释，特别是对于公共API
- 使用gofmt格式化你的代码

## 测试

所有代码修改都应该有对应的测试。我们使用标准的Go测试框架。

- 单元测试应放在与被测代码相同的包中，文件名以`_test.go`结尾
- 确保测试覆盖率足够高，特别是对于核心功能
- 使用`go test ./...`命令运行所有测试

## GitHub Actions

本项目使用GitHub Actions进行持续集成。每次代码提交和Pull Request都会触发以下工作流：

### 测试工作流 (test.yml)

这个工作流执行以下操作：

1. **单元测试**：运行项目中的所有测试
   ```go
   go test ./pkg/... -v -race -coverprofile=coverage.txt -covermode=atomic
   ```

2. **代码覆盖率报告**：生成测试覆盖率报告并上传到Codecov
   ```yaml
   - name: Upload coverage report
     uses: codecov/codecov-action@v3
     with:
       file: ./coverage.txt
       flags: unittests
   ```

3. **性能基准测试**：运行基准测试以评估性能
   ```go
   go test -bench=. ./pkg/...
   ```

4. **示例代码验证**：运行所有示例代码以确保它们能正常工作

### 代码质量检查 (lint)

使用golangci-lint进行代码质量检查，确保代码符合Go的最佳实践。

## 本地运行GitHub Actions

你可以在提交代码前在本地运行类似的测试：

```bash
# 运行单元测试
go test ./pkg/... -v

# 运行基准测试
go test -bench=. ./pkg/...

# 运行示例代码
cd examples/01_parse_package_json && go run main.go
# ... 运行其他示例
```

## 文档

对于新功能或修改，请确保：

1. 更新对应的README.md文件
2. 添加或更新代码文档注释
3. 如果修改了API，更新对应的示例代码

## 版本控制

本项目遵循[语义化版本](https://semver.org/)规范：

- MAJOR：不兼容的API变更
- MINOR：向后兼容的功能性新增
- PATCH：向后兼容的bug修复

## 问题和功能请求

- 使用GitHub Issues报告bug或请求新功能
- 在提交bug报告时，请尽可能提供详细的复现步骤和环境信息 