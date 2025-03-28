package parser

import (
	"context"
	"encoding/json"
	"runtime"
	"strings"
	"sync"

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

	// 从packages字段解析依赖
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)

	// 遍历packages
	for pkgPath, pkg := range packageLock.Packages {
		// 排除根包(通常是空字符串或"")
		if pkgPath == "" || pkgPath == "." {
			continue
		}

		// 解析包路径，格式通常为node_modules/package或node_modules/scope/package
		// 我们只关心包名而不关心嵌套路径
		parts := strings.Split(pkgPath, "/")
		var packageName string
		if len(parts) > 1 {
			// 处理作用域包 (@scope/package)
			if len(parts) > 2 && strings.HasPrefix(parts[1], "@") {
				packageName = parts[1] + "/" + parts[2]
			} else {
				packageName = parts[len(parts)-1]
			}
		} else {
			packageName = pkgPath
		}

		// 创建依赖对象
		if pkg.Version != "" {
			dependency := &baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]{}
			dependency.DependencyName = packageName
			dependency.DependencyVersion = pkg.Version

			// 设置生态系统特定字段
			ecosystem := &models.PackageLockComponentDependencyEcosystem{}
			ecosystem.Resolved = pkg.Resolved
			ecosystem.Integrity = pkg.Integrity
			ecosystem.Dev = pkg.Dev

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

	// 预估依赖数量以预分配内存
	estimatedDeps := len(packageLock.Packages)
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0, estimatedDeps)

	// 用于保护dependencies的并发访问
	var mu sync.Mutex

	// 并发处理依赖项
	numWorkers := runtime.NumCPU()
	wg := sync.WaitGroup{}

	// 将包按协程数量分块
	packagePaths := make([]string, 0, len(packageLock.Packages))
	for pkgPath := range packageLock.Packages {
		// 排除根包
		if pkgPath == "" || pkgPath == "." {
			continue
		}
		packagePaths = append(packagePaths, pkgPath)
	}

	// 按协程数分割任务
	chunkSize := (len(packagePaths) + numWorkers - 1) / numWorkers
	if chunkSize < 1 {
		chunkSize = 1
	}

	// 启动协程处理每个块
	for i := 0; i < numWorkers && i*chunkSize < len(packagePaths); i++ {
		wg.Add(1)

		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(packagePaths) {
			end = len(packagePaths)
		}

		go func(pathsChunk []string) {
			defer wg.Done()

			// 本地收集依赖项，减少锁竞争
			localDeps := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0, len(pathsChunk))

			for _, pkgPath := range pathsChunk {
				pkg := packageLock.Packages[pkgPath]

				// 创建依赖对象
				if pkg.Version != "" {
					// 解析包路径，获取真正的包名
					packageName := extractPackageNameFromPath(pkgPath)

					dependency := &baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]{}
					dependency.DependencyName = packageName
					dependency.DependencyVersion = pkg.Version

					// 设置生态系统特定字段
					ecosystem := &models.PackageLockComponentDependencyEcosystem{}
					ecosystem.Resolved = pkg.Resolved
					ecosystem.Integrity = pkg.Integrity
					ecosystem.Dev = pkg.Dev

					// 添加到本地依赖列表
					localDeps = append(localDeps, dependency)
				}
			}

			// 一次性添加本地收集的依赖项到全局列表
			if len(localDeps) > 0 {
				mu.Lock()
				dependencies = append(dependencies, localDeps...)
				mu.Unlock()
			}
		}(packagePaths[start:end])
	}

	// 等待所有协程完成
	wg.Wait()

	module.Dependencies = dependencies
	return module
}

// 从包路径中提取包名
func extractPackageNameFromPath(pkgPath string) string {
	parts := strings.Split(pkgPath, "/")
	var packageName string

	if len(parts) > 1 {
		// 处理作用域包 (@scope/package)
		if len(parts) > 2 && strings.HasPrefix(parts[1], "@") {
			packageName = parts[1] + "/" + parts[2]
		} else {
			packageName = parts[len(parts)-1]
		}
	} else {
		packageName = pkgPath
	}

	return packageName
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
