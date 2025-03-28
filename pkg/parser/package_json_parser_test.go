package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scagogogo/package-json-parser/pkg/models"
	baseModels "github.com/scagogogo/sca-base-module-components/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageJsonParser_Parse(t *testing.T) {
	// 创建测试目录
	testDir, err := os.MkdirTemp("", "package_json_parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	// 测试用例
	tests := []struct {
		name      string
		content   string
		wantError bool
		checkFunc func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem])
	}{
		{
			name: "基本的package.json文件",
			content: `{
				"name": "basic-package",
				"version": "1.0.0",
				"description": "A basic package for testing",
				"main": "index.js",
				"dependencies": {
					"lodash": "^4.17.21",
					"react": "^17.0.2"
				},
				"devDependencies": {
					"jest": "^27.0.6"
				}
			}`,
			wantError: false,
			checkFunc: func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]) {
				assert.NotNil(t, project)
				assert.Equal(t, "basic-package", project.Name)
				assert.Equal(t, "1.0.0", project.Version)

				// 验证模块
				assert.Len(t, project.Modules, 1)
				// 查找特定名称的模块
				var module *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]
				for _, m := range project.Modules {
					if m.Name == "basic-package" {
						module = m
						break
					}
				}
				require.NotNil(t, module)
				assert.Equal(t, "basic-package", module.Name)
				assert.Equal(t, "1.0.0", module.Version)

				// 验证依赖
				assert.Len(t, module.Dependencies, 3) // 2个dependencies + 1个devDependencies

				// 查找指定依赖
				var lodashDep, reactDep, jestDep *baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]
				for i := range module.Dependencies {
					switch module.Dependencies[i].DependencyName {
					case "lodash":
						lodashDep = module.Dependencies[i]
					case "react":
						reactDep = module.Dependencies[i]
					case "jest":
						jestDep = module.Dependencies[i]
					}
				}

				// 验证lodash依赖
				require.NotNil(t, lodashDep)
				assert.Equal(t, "^4.17.21", lodashDep.DependencyVersion)
				assert.Nil(t, lodashDep.ComponentDependencyEcosystem.Dev) // 非开发依赖

				// 验证react依赖
				require.NotNil(t, reactDep)
				assert.Equal(t, "^17.0.2", reactDep.DependencyVersion)
				assert.Nil(t, reactDep.ComponentDependencyEcosystem.Dev) // 非开发依赖

				// 验证jest依赖
				require.NotNil(t, jestDep)
				assert.Equal(t, "^27.0.6", jestDep.DependencyVersion)
				assert.NotNil(t, jestDep.ComponentDependencyEcosystem.Dev)
				assert.True(t, *jestDep.ComponentDependencyEcosystem.Dev) // 开发依赖
			},
		},
		{
			name: "包含所有类型依赖的package.json",
			content: `{
				"name": "all-deps-package",
				"version": "1.0.0",
				"dependencies": {
					"express": "^4.17.1"
				},
				"devDependencies": {
					"eslint": "^7.32.0"
				},
				"peerDependencies": {
					"react": "^16.0.0 || ^17.0.0"
				},
				"optionalDependencies": {
					"colors": "^1.4.0"
				},
				"bundledDependencies": ["moment"]
			}`,
			wantError: false,
			checkFunc: func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]) {
				assert.NotNil(t, project)

				// 查找特定名称的模块
				var module *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]
				for _, m := range project.Modules {
					if m.Name == "all-deps-package" {
						module = m
						break
					}
				}
				require.NotNil(t, module)
				// 注意: 当前实现仅处理dependencies和devDependencies
				assert.Len(t, module.Dependencies, 2) // 1 regular + 1 dev

				// 检查各类依赖
				var expressDep, eslintDep *baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]
				for i := range module.Dependencies {
					switch module.Dependencies[i].DependencyName {
					case "express":
						expressDep = module.Dependencies[i]
					case "eslint":
						eslintDep = module.Dependencies[i]
					}
				}

				// 验证普通依赖
				require.NotNil(t, expressDep)
				assert.Equal(t, "^4.17.1", expressDep.DependencyVersion)
				assert.Nil(t, expressDep.ComponentDependencyEcosystem.Dev) // 非开发依赖

				// 验证开发依赖
				require.NotNil(t, eslintDep)
				assert.Equal(t, "^7.32.0", eslintDep.DependencyVersion)
				assert.NotNil(t, eslintDep.ComponentDependencyEcosystem.Dev)
				assert.True(t, *eslintDep.ComponentDependencyEcosystem.Dev) // 开发依赖
			},
		},
		{
			name: "workspaces字段的package.json",
			content: `{
				"name": "monorepo-root",
				"version": "1.0.0",
				"private": true,
				"workspaces": [
					"packages/*"
				],
				"dependencies": {
					"shared-lib": "^1.0.0"
				}
			}`,
			wantError: false,
			checkFunc: func(t *testing.T, project *baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]) {
				assert.NotNil(t, project)

				// 验证模块
				assert.Len(t, project.Modules, 1)
				// 查找特定名称的模块
				var module *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]
				for _, m := range project.Modules {
					if m.Name == "monorepo-root" {
						module = m
						break
					}
				}
				require.NotNil(t, module)
				assert.Equal(t, "monorepo-root", module.Name)

				// 验证依赖
				assert.Len(t, module.Dependencies, 1)
				dep := module.Dependencies[0]
				assert.Equal(t, "shared-lib", dep.DependencyName)
				assert.Equal(t, "^1.0.0", dep.DependencyVersion)
				assert.Nil(t, dep.ComponentDependencyEcosystem.Dev) // 非开发依赖
			},
		},
		{
			name: "无效的JSON格式",
			content: `{
				"name": "invalid-json",
				"version": "1.0.0",
				"dependencies": {
					"lodash": "^4.17.21",
				}
			}`,
			wantError: true,
		},
		{
			name:      "空文件",
			content:   ``,
			wantError: true,
		},
		{
			name: "无name字段",
			content: `{
				"version": "1.0.0",
				"dependencies": {
					"lodash": "^4.17.21"
				}
			}`,
			wantError: true, // 当前实现要求name字段
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试文件
			packageJsonPath := filepath.Join(testDir, fmt.Sprintf("%s.json", tt.name))
			err := os.WriteFile(packageJsonPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// 创建解析器并执行
			parser := &PackageJsonParser{}
			err = parser.Init(context.Background())
			require.NoError(t, err)

			input := &PackageJsonParserInput{
				PackageJsonPath: packageJsonPath,
			}

			// 解析并验证结果
			project, err := parser.Parse(context.Background(), input)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, project)

			if tt.checkFunc != nil {
				tt.checkFunc(t, project)
			}
		})
	}
}

func TestPackageJsonParser_ParseErrors(t *testing.T) {
	parser := &PackageJsonParser{}
	err := parser.Init(context.Background())
	require.NoError(t, err)

	// 测试不存在的文件
	input := &PackageJsonParserInput{
		PackageJsonPath: "/path/to/nonexistent/package.json",
	}
	project, err := parser.Parse(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, project)
}

// 测试createDependency函数
func TestPackageJsonParser_CreateDependency(t *testing.T) {
	parser := &PackageJsonParser{}

	tests := []struct {
		name         string
		depName      string
		depVersion   string
		isDev        bool
		expectedName string
		expectedVer  string
		devShouldBe  bool
	}{
		{
			name:         "普通依赖",
			depName:      "lodash",
			depVersion:   "^4.17.21",
			isDev:        false,
			expectedName: "lodash",
			expectedVer:  "^4.17.21",
			devShouldBe:  false,
		},
		{
			name:         "开发依赖",
			depName:      "jest",
			depVersion:   "^27.0.6",
			isDev:        true,
			expectedName: "jest",
			expectedVer:  "^27.0.6",
			devShouldBe:  true,
		},
		{
			name:         "作用域包",
			depName:      "@types/node",
			depVersion:   "^16.11.7",
			isDev:        true,
			expectedName: "@types/node",
			expectedVer:  "^16.11.7",
			devShouldBe:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dependency := parser.createDependency(tt.depName, tt.depVersion, tt.isDev)

			assert.Equal(t, tt.expectedName, dependency.DependencyName)
			assert.Equal(t, tt.expectedVer, dependency.DependencyVersion)

			if tt.isDev {
				assert.NotNil(t, dependency.ComponentDependencyEcosystem.Dev)
				assert.True(t, *dependency.ComponentDependencyEcosystem.Dev)
			} else {
				// 非开发依赖的Dev应为nil
				assert.Nil(t, dependency.ComponentDependencyEcosystem.Dev)
			}
		})
	}
}

// 测试parseComponent函数
func TestPackageJsonParser_ParseComponent(t *testing.T) {
	parser := &PackageJsonParser{}

	component := parser.parseComponent("test-component", "1.0.0")

	assert.Equal(t, "test-component", component.Name)
	assert.Equal(t, "1.0.0", component.Version)
	assert.NotNil(t, component.ComponentEcosystem)
}

func TestPackageJsonParser_Close(t *testing.T) {
	parser := &PackageJsonParser{}
	ctx := context.Background()

	// 初始化解析器
	err := parser.Init(ctx)
	assert.NoError(t, err)

	// 关闭解析器
	err = parser.Close(ctx)
	assert.NoError(t, err)
}

func TestPackageJsonParser_GetName(t *testing.T) {
	parser := &PackageJsonParser{}
	name := parser.GetName()
	assert.Equal(t, PackageJsonParserName, name)
}
