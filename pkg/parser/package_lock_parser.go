package parser

import (
	"context"
	"encoding/json"
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

	project.SetModule(lock.Name, x.parseModule(lock))

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

// 解析所有的依赖
func (x *PackageLockParser) parseDependencies(packageJsonDependencies map[string]*models.PackageLockDependency) []*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem] {
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)
	for dependencyPackageName, packageLockDependency := range packageJsonDependencies {
		dependencies = append(dependencies, x.parseDependency(dependencyPackageName, packageLockDependency))
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

func (x *PackageLockParser) Close(ctx context.Context) error {
	return nil
}
