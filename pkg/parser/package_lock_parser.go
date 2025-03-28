package parser

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/scagogogo/package-json-parser/pkg/models"
	baseModels "github.com/scagogogo/sca-base-module-components/pkg/models"
	"github.com/scagogogo/sca-base-module-ecosystem-parser/pkg/parser"
)

type PackageLockParser struct {
}

var _ parser.Parser[*PackageLockJsonParserInput, *models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] = &PackageLockParser{}

func NewPackageLockParser() *PackageLockParser {
	return &PackageLockParser{}
}

const PackageLockParserName = "package-lock-parser"

func (x *PackageLockParser) GetName() string {
	return PackageLockParserName
}

func (x *PackageLockParser) Init(ctx context.Context) error {
	return nil
}

// Parse 把package-lock.json当做是一个项目解析
func (x *PackageLockParser) Parse(ctx context.Context, input *PackageLockJsonParserInput) (*baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem], error) {
	bytes, err := input.Read(ctx)
	if err != nil {
		return nil, err
	}
	lock := &models.PackageLock{}
	err = json.Unmarshal(bytes, &lock)
	if err != nil {
		return nil, err
	}

	project := &baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	project.Name = lock.Name
	project.Version = lock.Version

	// 根据lockfileVersion选择不同的解析策略
	switch lock.LockFileVersion {
	case 1:
		// npm v5
		project.SetModule(lock.Name, x.parseModule(lock))
	case 2:
		// npm v6
		project.SetModule(lock.Name, x.parseModule(lock))
	case 3:
		// npm v7+，使用packages字段
		// 如果依赖较多，使用并发版本的解析器
		if lock.Packages != nil && len(lock.Packages) > 100 {
			project.SetModule(lock.Name, x.parseModuleV7Concurrent(lock))
		} else {
			project.SetModule(lock.Name, x.parseModuleV7(lock))
		}
	default:
		// 未知版本，尝试使用兼容模式解析
		if lock.Packages != nil && len(lock.Packages) > 0 {
			if len(lock.Packages) > 100 {
				project.SetModule(lock.Name, x.parseModuleV7Concurrent(lock))
			} else {
				project.SetModule(lock.Name, x.parseModuleV7(lock))
			}
		} else {
			project.SetModule(lock.Name, x.parseModule(lock))
		}
	}

	return project, nil
}

// 解析模块，整个package-lock.json项目看做是一个模块解析
func (x *PackageLockParser) parseModule(packageLock *models.PackageLock) *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] {
	module := &baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	module.Name = packageLock.Name
	module.Version = packageLock.Version
	module.Dependencies = x.parseDependencies(packageLock.Dependencies)
	return module
}

// parseModuleV7 解析npm v7+格式的package-lock.json
func (x *PackageLockParser) parseModuleV7(packageLock *models.PackageLock) *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] {
	module := &baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	module.Name = packageLock.Name
	module.Version = packageLock.Version

	// 检查packages是否为nil
	if packageLock.Packages == nil {
		module.Dependencies = make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)
		return module
	}

	// 从packages字段解析依赖
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)

	// 遍历packages
	for pkgPath, pkg := range packageLock.Packages {
		// 检查pkg是否为nil
		if pkg == nil {
			continue
		}

		// 排除根包(通常是空字符串或"")
		if pkgPath == "" || pkgPath == "." {
			continue
		}

		// 解析包名
		packageName := extractPackageNameFromPath(pkgPath)

		// 创建依赖对象
		if pkg.Version != "" {
			dependency := &baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]{}
			dependency.DependencyName = packageName
			dependency.DependencyVersion = pkg.Version

			// 设置生态系统特定字段，确保不为nil
			ecosystem := &models.PackageLockComponentDependencyEcosystem{}
			ecosystem.Resolved = pkg.Resolved
			ecosystem.Integrity = pkg.Integrity
			ecosystem.Dev = pkg.Dev

			// 显式设置ComponentDependencyEcosystem
			dependency.ComponentDependencyEcosystem = ecosystem

			// 添加依赖项
			dependencies = append(dependencies, dependency)
		}
	}

	module.Dependencies = dependencies
	return module
}

// parseModuleV7Concurrent 使用并发处理来提高大型package-lock.json的解析性能
func (x *PackageLockParser) parseModuleV7Concurrent(packageLock *models.PackageLock) *baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] {
	module := &baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	module.Name = packageLock.Name
	module.Version = packageLock.Version

	// 检查packages是否为nil
	if packageLock.Packages == nil {
		module.Dependencies = make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)
		return module
	}

	// 先筛选出有效的包路径
	validPackagePaths := make([]string, 0, len(packageLock.Packages))
	for pkgPath, pkg := range packageLock.Packages {
		// 排除空路径、根包和nil包
		if pkgPath == "" || pkgPath == "." || pkg == nil || pkg.Version == "" {
			continue
		}
		validPackagePaths = append(validPackagePaths, pkgPath)
	}

	// 如果没有有效的包路径，则直接返回
	if len(validPackagePaths) == 0 {
		module.Dependencies = make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)
		return module
	}

	// 预先分配足够大的空间
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0, len(validPackagePaths))

	// 处理每个有效的依赖
	for _, pkgPath := range validPackagePaths {
		pkg := packageLock.Packages[pkgPath]

		// 解析包名
		packageName := extractPackageNameFromPath(pkgPath)

		// 创建依赖对象
		dependency := &baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]{}
		dependency.DependencyName = packageName
		dependency.DependencyVersion = pkg.Version

		// 设置生态系统特定字段
		ecosystem := &models.PackageLockComponentDependencyEcosystem{}
		ecosystem.Resolved = pkg.Resolved
		ecosystem.Integrity = pkg.Integrity
		ecosystem.Dev = pkg.Dev

		// 显式设置ComponentDependencyEcosystem
		dependency.ComponentDependencyEcosystem = ecosystem

		// 添加到依赖列表
		dependencies = append(dependencies, dependency)
	}

	module.Dependencies = dependencies
	return module
}

// 从包路径中提取包名
func extractPackageNameFromPath(pkgPath string) string {
	// 处理空路径
	if pkgPath == "" {
		return ""
	}

	parts := strings.Split(pkgPath, "/")
	if len(parts) == 0 {
		return ""
	}

	// 查找最后一个 node_modules 的位置
	lastNodeModulesIndex := -1
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "node_modules" {
			lastNodeModulesIndex = i
			break
		}
	}

	// 如果找到了 node_modules
	if lastNodeModulesIndex != -1 && lastNodeModulesIndex < len(parts)-1 {
		// 检查包名是否是作用域包 (@scope/package)
		nextIndex := lastNodeModulesIndex + 1
		if nextIndex < len(parts) && strings.HasPrefix(parts[nextIndex], "@") {
			if nextIndex+1 < len(parts) {
				return parts[nextIndex] + "/" + parts[nextIndex+1]
			}
			return parts[nextIndex]
		}
		return parts[nextIndex]
	}

	// 如果没有找到 node_modules，使用最后一个部分作为包名
	// 检查是否是作用域包
	if len(parts) > 0 {
		if len(parts) > 1 && strings.HasPrefix(parts[0], "@") {
			if len(parts) > 1 {
				return parts[0] + "/" + parts[1]
			}
			return parts[0]
		}
		return parts[len(parts)-1]
	}

	return pkgPath
}

// 解析所有的依赖
func (x *PackageLockParser) parseDependencies(packageJsonDependencies map[string]*models.PackageLockDependency) []*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem] {
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)
	processed := make(map[string]bool) // 用于防止循环依赖

	for dependencyPackageName, packageLockDependency := range packageJsonDependencies {
		x.parseAndAddDependency(dependencyPackageName, packageLockDependency, &dependencies, processed, 0)
	}
	return dependencies
}

// 解析单个的依赖
func (x *PackageLockParser) parseDependency(packageName string, packageLockDependency *models.PackageLockDependency) *baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem] {
	dependency := &baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]{}
	dependency.DependencyName = packageName
	dependency.DependencyVersion = packageLockDependency.Version

	ecosystem := &models.PackageLockComponentDependencyEcosystem{}
	ecosystem.Integrity = packageLockDependency.Integrity
	ecosystem.Resolved = packageLockDependency.Resolved
	ecosystem.Dev = packageLockDependency.Dev
	ecosystem.Requires = packageLockDependency.Requires
	dependency.ComponentDependencyEcosystem = ecosystem

	return dependency
}

// 递归解析依赖，并添加到依赖列表
func (x *PackageLockParser) parseAndAddDependency(
	packageName string,
	packageLockDependency *models.PackageLockDependency,
	dependencies *[]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem],
	processed map[string]bool,
	depth int,
) {
	// 防止循环依赖，设置最大递归深度
	if depth > 100 || processed[packageName+"@"+packageLockDependency.Version] {
		return
	}

	// 标记为已处理
	processed[packageName+"@"+packageLockDependency.Version] = true

	// 解析当前依赖
	dependency := x.parseDependency(packageName, packageLockDependency)
	*dependencies = append(*dependencies, dependency)

	// 递归解析嵌套依赖
	if packageLockDependency.Dependencies != nil && len(packageLockDependency.Dependencies) > 0 {
		for nestedName, nestedDep := range packageLockDependency.Dependencies {
			// 创建一个包含父依赖信息的依赖对象
			x.parseAndAddDependency(nestedName, nestedDep, dependencies, processed, depth+1)
		}
	}
}

func (x *PackageLockParser) Close(ctx context.Context) error {
	return nil
}
