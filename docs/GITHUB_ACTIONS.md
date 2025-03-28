# GitHub Actions 工作流程

本项目使用GitHub Actions实现自动化测试和持续集成。以下是有关工作流的详细信息。

## 工作流程概述

项目的主要工作流程定义在 `.github/workflows/test.yml` 文件中，包含两个主要任务：

1. **测试（test）**: 运行单元测试、覆盖率报告、基准测试，并验证所有示例代码
2. **代码规范检查（lint）**: 使用golangci-lint检查代码质量

## 触发条件

工作流在以下情况下会被触发：

- 向主分支（main）推送代码
- 创建针对主分支的Pull Request

```yaml
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
```

## 测试任务详情

测试任务包含以下步骤：

### 1. 设置Go环境

```yaml
- name: Set up Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.21'
```

### 2. 检出代码

```yaml
- name: Check out code
  uses: actions/checkout@v4
```

### 3. 获取依赖

```yaml
- name: Get dependencies
  run: go mod download
```

### 4. 运行单元测试并生成覆盖率报告

```yaml
- name: Run unit tests with coverage
  run: go test ./pkg/... -v -race -coverprofile=coverage.txt -covermode=atomic
```

这个命令：
- 运行`pkg`目录及其子目录中的所有测试
- 启用竞态条件检测（-race）
- 生成覆盖率报告（-coverprofile=coverage.txt）
- 采用原子模式计算覆盖率（-covermode=atomic）

### 5. 上传覆盖率报告

```yaml
- name: Upload coverage report
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage.txt
    flags: unittests
    fail_ci_if_error: false
```

### 6. 运行基准测试

```yaml
- name: Run benchmark tests
  run: go test -bench=. ./pkg/...
```

### 7. 运行各个示例代码

```yaml
- name: Test example 01 - Parse package.json
  run: |
    cd examples/01_parse_package_json
    go run main.go
    
# ... 其他示例
```

## 代码规范检查任务

代码规范检查任务使用golangci-lint工具：

```yaml
lint:
  name: Lint
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
        
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
```

## 查看结果

每次工作流运行后，你可以在GitHub仓库的Actions选项卡中查看详细结果：

1. 进入GitHub仓库
2. 点击"Actions"选项卡
3. 选择最近的工作流运行
4. 展开各个步骤查看详细日志

## 常见问题解决

### 测试失败

如果测试失败，请检查：

1. 错误日志以了解具体的失败原因
2. 确保所有依赖都已正确安装
3. 确保本地测试通过后再提交代码

### 代码规范检查失败

如果代码规范检查失败，请：

1. 查看详细的lint错误报告
2. 在本地安装golangci-lint并运行检查
   ```bash
   # 安装golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # 运行检查
   golangci-lint run
   ```
3. 修复所有报告的问题后再次提交 