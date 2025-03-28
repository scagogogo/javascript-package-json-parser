package parser

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/scagogogo/package-json-parser/pkg/models"
	baseModels "github.com/scagogogo/sca-base-module-components/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageLockParser_Parse(t *testing.T) {
	parser := NewPackageLockParser()

	err := parser.Init(context.Background())
	assert.Nil(t, err)

	tests := []struct {
		name      string
		inputFile string
		wantError bool
		checkFunc func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem])
	}{
		{
			name:      "基本package-lock.json解析",
			inputFile: "./test_data/package-lock.json/join-dev-design.json",
			wantError: false,
			checkFunc: func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]) {
				assert.NotNil(t, project)
				assert.Equal(t, "join-dev-design", project.Name)

				// 检查模块
				module := findModuleInPackageLock(project, project.Name)
				require.NotNil(t, module)

				// 检查依赖数量
				assert.True(t, len(module.Dependencies) > 0, "应该有依赖")

				// 检查特定依赖
				found := false
				for _, dep := range module.Dependencies {
					if dep.DependencyName == "glob" {
						found = true
						assert.NotEmpty(t, dep.DependencyVersion)
						break
					}
				}
				assert.True(t, found, "应该包含glob依赖")
			},
		},
		{
			name:      "较大的package-lock.json",
			inputFile: "./test_data/package-lock.json/universal-module-tree.json",
			wantError: false,
			checkFunc: func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]) {
				assert.NotNil(t, project)

				// 检查模块
				module := findModuleInPackageLock(project, project.Name)
				require.NotNil(t, module)

				// 较大文件应该有很多依赖
				assert.True(t, len(module.Dependencies) > 20, "大型package-lock.json应该有很多依赖")
			},
		},
		{
			name:      "最小的package-lock.json",
			inputFile: "./test_data/package-lock.json/gitlab-ci-yarn-audit-parser.json",
			wantError: false,
			checkFunc: func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]) {
				assert.NotNil(t, project)

				// 检查模块
				module := findModuleInPackageLock(project, project.Name)
				require.NotNil(t, module)

				// 这个文件的依赖应该很少或没有
				t.Logf("依赖数量: %d", len(module.Dependencies))
			},
		},
		// 测试错误情况
		{
			name:      "空文件",
			inputFile: createEmptyFile(t),
			wantError: true,
		},
		{
			name:      "无效JSON格式",
			inputFile: createInvalidJsonFile(t),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &PackageLockJsonParserInput{
				PackageLockJsonPath: tt.inputFile,
			}

			project, err := parser.Parse(context.Background(), input)

			if tt.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, project)
				}
			}
		})
	}
}

// 测试解析器的各个方法
func TestPackageLockParser_Methods(t *testing.T) {
	// 测试GetName方法
	t.Run("GetName", func(t *testing.T) {
		parser := NewPackageLockParser()
		assert.Equal(t, PackageLockParserName, parser.GetName())
	})

	// 测试Init方法
	t.Run("Init", func(t *testing.T) {
		parser := NewPackageLockParser()
		err := parser.Init(context.Background())
		assert.Nil(t, err)
	})

	// 测试Close方法
	t.Run("Close", func(t *testing.T) {
		parser := NewPackageLockParser()
		err := parser.Close(context.Background())
		assert.Nil(t, err)
	})
}

// 测试parseModule方法
func TestPackageLockParser_ParseModule(t *testing.T) {
	parser := NewPackageLockParser()

	// 创建一个简单的PackageLock对象
	packageLock := &models.PackageLock{
		Name:    "test-package",
		Version: "1.0.0",
		Dependencies: map[string]*models.PackageLockDependency{
			"dep1": {
				Version:   "1.2.3",
				Resolved:  "https://registry.npmjs.org/dep1/-/dep1-1.2.3.tgz",
				Integrity: "sha512-abc123",
			},
			"dep2": {
				Version:   "4.5.6",
				Resolved:  "https://registry.npmjs.org/dep2/-/dep2-4.5.6.tgz",
				Integrity: "sha512-def456",
				Dependencies: map[string]*models.PackageLockDependency{
					"nested-dep": {
						Version:   "7.8.9",
						Resolved:  "https://registry.npmjs.org/nested-dep/-/nested-dep-7.8.9.tgz",
						Integrity: "sha512-ghi789",
					},
				},
			},
		},
	}

	// 测试parseModule方法
	module := parser.parseModule(packageLock)

	// 验证结果
	assert.Equal(t, "test-package", module.Name)
	assert.Equal(t, "1.0.0", module.Version)
	assert.Len(t, module.Dependencies, 3) // 2个顶级依赖 + 1个嵌套依赖

	// 检查依赖
	var dep1, dep2, nestedDep *baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]
	for _, dep := range module.Dependencies {
		if dep.DependencyName == "dep1" {
			dep1 = dep
		} else if dep.DependencyName == "dep2" {
			dep2 = dep
		} else if dep.DependencyName == "nested-dep" {
			nestedDep = dep
		}
	}

	// 验证dep1
	require.NotNil(t, dep1)
	assert.Equal(t, "1.2.3", dep1.DependencyVersion)
	assert.Equal(t, "https://registry.npmjs.org/dep1/-/dep1-1.2.3.tgz", dep1.ComponentDependencyEcosystem.Resolved)
	assert.Equal(t, "sha512-abc123", dep1.ComponentDependencyEcosystem.Integrity)

	// 验证dep2
	require.NotNil(t, dep2)
	assert.Equal(t, "4.5.6", dep2.DependencyVersion)
	assert.Equal(t, "https://registry.npmjs.org/dep2/-/dep2-4.5.6.tgz", dep2.ComponentDependencyEcosystem.Resolved)
	assert.Equal(t, "sha512-def456", dep2.ComponentDependencyEcosystem.Integrity)

	// 验证嵌套依赖
	require.NotNil(t, nestedDep)
	assert.Equal(t, "7.8.9", nestedDep.DependencyVersion)
	assert.Equal(t, "https://registry.npmjs.org/nested-dep/-/nested-dep-7.8.9.tgz", nestedDep.ComponentDependencyEcosystem.Resolved)
	assert.Equal(t, "sha512-ghi789", nestedDep.ComponentDependencyEcosystem.Integrity)
}

// 测试parseAndAddDependency方法
func TestPackageLockParser_ParseAndAddDependency(t *testing.T) {
	parser := NewPackageLockParser()

	// 创建测试数据
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)
	processed := make(map[string]bool)

	// 创建测试依赖
	dep := &models.PackageLockDependency{
		Version:   "1.0.0",
		Resolved:  "https://registry.npmjs.org/test/-/test-1.0.0.tgz",
		Integrity: "sha512-test123",
		Dependencies: map[string]*models.PackageLockDependency{
			"nested": {
				Version:   "2.0.0",
				Resolved:  "https://registry.npmjs.org/nested/-/nested-2.0.0.tgz",
				Integrity: "sha512-nested456",
			},
		},
	}

	// 执行方法
	parser.parseAndAddDependency("test-package", dep, &dependencies, processed, 0)

	// 验证结果
	assert.Len(t, dependencies, 2) // 一个主依赖和一个嵌套依赖
	assert.True(t, processed["test-package@1.0.0"])

	// 检查主依赖
	assert.Equal(t, "test-package", dependencies[0].DependencyName)
	assert.Equal(t, "1.0.0", dependencies[0].DependencyVersion)
	assert.Equal(t, "https://registry.npmjs.org/test/-/test-1.0.0.tgz", dependencies[0].ComponentDependencyEcosystem.Resolved)

	// 检查嵌套依赖
	assert.Equal(t, "nested", dependencies[1].DependencyName)
	assert.Equal(t, "2.0.0", dependencies[1].DependencyVersion)
}

// 测试parseDependency方法
func TestPackageLockParser_ParseDependency(t *testing.T) {
	parser := NewPackageLockParser()

	// 创建测试依赖
	packageLockDependency := &models.PackageLockDependency{
		Version:   "1.2.3",
		Resolved:  "https://registry.npmjs.org/test/-/test-1.2.3.tgz",
		Integrity: "sha512-test123",
	}

	// 设置Dev字段
	dev := true
	packageLockDependency.Dev = &dev

	// 执行方法
	result := parser.parseDependency("test-package", packageLockDependency)

	// 验证结果
	assert.Equal(t, "test-package", result.DependencyName)
	assert.Equal(t, "1.2.3", result.DependencyVersion)
	assert.Equal(t, "https://registry.npmjs.org/test/-/test-1.2.3.tgz", result.ComponentDependencyEcosystem.Resolved)
	assert.Equal(t, "sha512-test123", result.ComponentDependencyEcosystem.Integrity)
	assert.NotNil(t, result.ComponentDependencyEcosystem.Dev)
	assert.True(t, *result.ComponentDependencyEcosystem.Dev)
}

// 测试extractPackageNameFromPath函数
func TestExtractPackageNameFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			path:     "node_modules/lodash",
			expected: "lodash",
		},
		{
			path:     "node_modules/@babel/core",
			expected: "@babel/core",
		},
		{
			path:     "node_modules/nested/node_modules/react",
			expected: "react",
		},
		{
			path:     "node_modules/@scope/nested/node_modules/@another/package",
			expected: "@another/package",
		},
		{
			path:     "single-package",
			expected: "single-package",
		},
		{
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := extractPackageNameFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// 测试parseModuleV7方法
func TestPackageLockParser_ParseModuleV7(t *testing.T) {
	parser := NewPackageLockParser()

	// 创建一个简单的npm v7+ PackageLock对象
	packageLock := &models.PackageLock{
		Name:    "test-v7-package",
		Version: "1.0.0",
		Packages: map[string]*models.PackageLockPackage{
			"": { // 根包
				Version: "1.0.0",
			},
			"node_modules/lodash": {
				Version:   "4.17.21",
				Resolved:  "https://registry.npmjs.org/lodash/-/lodash-4.17.21.tgz",
				Integrity: "sha512-lodash-hash",
			},
			"node_modules/@babel/core": {
				Version:   "7.15.0",
				Resolved:  "https://registry.npmjs.org/@babel/core/-/core-7.15.0.tgz",
				Integrity: "sha512-babel-hash",
			},
		},
	}

	// 执行方法
	module := parser.parseModuleV7(packageLock)

	// 验证结果
	assert.Equal(t, "test-v7-package", module.Name)
	assert.Equal(t, "1.0.0", module.Version)
	assert.Len(t, module.Dependencies, 2) // 应该有两个依赖（不包括根包）

	// 验证依赖
	var lodashDep, babelDep *baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]
	for _, dep := range module.Dependencies {
		if dep.DependencyName == "lodash" {
			lodashDep = dep
		} else if dep.DependencyName == "@babel/core" {
			babelDep = dep
		}
	}

	require.NotNil(t, lodashDep)
	assert.Equal(t, "4.17.21", lodashDep.DependencyVersion)
	assert.Equal(t, "https://registry.npmjs.org/lodash/-/lodash-4.17.21.tgz", lodashDep.ComponentDependencyEcosystem.Resolved)

	require.NotNil(t, babelDep)
	assert.Equal(t, "7.15.0", babelDep.DependencyVersion)
	require.NotNil(t, babelDep.ComponentDependencyEcosystem)
	assert.Equal(t, "https://registry.npmjs.org/@babel/core/-/core-7.15.0.tgz", babelDep.ComponentDependencyEcosystem.Resolved)
}

// 测试parseModuleV7Concurrent方法
func TestPackageLockParser_ParseModuleV7Concurrent(t *testing.T) {
	parser := NewPackageLockParser()

	// 创建一个简单的npm v7+ PackageLock对象，包含多个依赖
	packageLock := &models.PackageLock{
		Name:     "test-concurrent",
		Version:  "1.0.0",
		Packages: map[string]*models.PackageLockPackage{},
	}

	// 添加正好150个依赖，确保路径不重复
	packagesAdded := 0
	for i := 0; packagesAdded < 150; i++ {
		// 确保生成完全不同的路径字符串
		packageName := fmt.Sprintf("node_modules/uniquepkg-%d", i)

		packageLock.Packages[packageName] = &models.PackageLockPackage{
			Version:   fmt.Sprintf("1.0.%d", i%10),
			Resolved:  "https://registry.npmjs.org/example/-/example-1.0.0.tgz",
			Integrity: "sha512-example",
		}

		packagesAdded++
	}

	// 确认我们正好添加了150个包
	if len(packageLock.Packages) != 150 {
		t.Fatalf("Expected to add 150 packages but got %d", len(packageLock.Packages))
	}

	// 执行方法
	module := parser.parseModuleV7Concurrent(packageLock)

	// 验证结果
	assert.Equal(t, "test-concurrent", module.Name)
	assert.Equal(t, "1.0.0", module.Version)
	assert.Len(t, module.Dependencies, 150) // 应该有150个依赖
}

// 辅助函数：查找指定名称的模块
func findModuleInPackageLock(project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem], name string) *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] {
	for _, module := range project.Modules {
		if module.Name == name {
			return module
		}
	}
	return nil
}

// 辅助函数：创建临时的空文件
func createEmptyFile(t *testing.T) string {
	file, err := os.CreateTemp("", "package-lock-empty-*.json")
	require.NoError(t, err)
	defer file.Close()
	return file.Name()
}

// 辅助函数：创建临时的无效JSON文件
func createInvalidJsonFile(t *testing.T) string {
	file, err := os.CreateTemp("", "package-lock-invalid-*.json")
	require.NoError(t, err)
	defer file.Close()

	_, err = file.Write([]byte("{this is not valid json}"))
	require.NoError(t, err)

	return file.Name()
}
