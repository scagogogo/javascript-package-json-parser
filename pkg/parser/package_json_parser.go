package parser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scagogogo/package-json-parser/pkg/models"
	baseModels "github.com/scagogogo/sca-base-module-components/pkg/models"
	"github.com/scagogogo/sca-base-module-ecosystem-parser/pkg/parser"
)

// 自定义错误类型
type PackageJsonParserError struct {
	Stage string // 错误发生的阶段
	Msg   string // 错误消息
	Err   error  // 原始错误
}

func (e *PackageJsonParserError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("PackageJsonParser error in %s stage: %s - %v", e.Stage, e.Msg, e.Err)
	}
	return fmt.Sprintf("PackageJsonParser error in %s stage: %s", e.Stage, e.Msg)
}

func (e *PackageJsonParserError) Unwrap() error {
	return e.Err
}

// 定义一个包装错误的辅助函数
func wrapError(stage, msg string, err error) *PackageJsonParserError {
	return &PackageJsonParserError{
		Stage: stage,
		Msg:   msg,
		Err:   err,
	}
}

// PackageJsonParser 是一个用于解析package.json文件的解析器
// 它实现了Parser接口，可以将package.json文件解析为结构化的项目、模块和组件对象
type PackageJsonParser struct {
}

var _ parser.Parser[*PackageJsonParserInput, *models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] = &PackageJsonParser{}

// PackageJsonParserName 是此解析器的唯一标识名称
const PackageJsonParserName = "package-json-parser"

// GetName 返回解析器的唯一名称
// 实现了Parser接口的GetName方法
func (x *PackageJsonParser) GetName() string {
	return PackageJsonParserName
}

// Init 初始化解析器
// 实现了Parser接口的Init方法
// 当前实现不需要特殊的初始化步骤
func (x *PackageJsonParser) Init(ctx context.Context) error {
	return nil
}

// Parse 解析package.json文件并转换为项目对象
// 参数:
//   - ctx: 上下文，用于控制解析过程
//   - input: 解析器输入，包含package.json文件路径
//
// 返回:
//   - 解析后的项目对象
//   - 如果解析过程中出现错误，则返回错误
func (x *PackageJsonParser) Parse(ctx context.Context, input *PackageJsonParserInput) (*baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem], error) {

	// 验证输入
	if input == nil {
		return nil, wrapError("input validation", "input cannot be nil", nil)
	}

	// 如果是通过内容传入的，不需要检查路径
	if input.PackageJsonContent == "" && input.PackageJsonPath == "" {
		return nil, wrapError("input validation", "package.json path cannot be empty", nil)
	}

	packageJsonBytes, err := input.Read(ctx)
	if err != nil {
		return nil, wrapError("file reading", fmt.Sprintf("failed to read package.json from %s", input.PackageJsonPath), err)
	}

	// 检查文件大小
	if len(packageJsonBytes) == 0 {
		return nil, wrapError("file validation", "package.json is empty", nil)
	}

	packageJson := &models.PackageJson{}
	err = json.Unmarshal(packageJsonBytes, &packageJson)
	if err != nil {
		return nil, wrapError("json parsing", "failed to parse package.json content", err)
	}

	// 验证解析结果
	if packageJson.Name == "" {
		return nil, wrapError("content validation", "package.json must have a name field", nil)
	}

	project := &baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	project.Name = packageJson.Name
	project.Version = packageJson.Version

	// 创建模块
	module := &baseModels.Module[*models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	module.Name = packageJson.Name
	module.Version = packageJson.Version

	// 设置模块生态系统信息
	moduleEcosystem := &models.PackageLockModuleEcosystem{}
	// 可以添加更多模块生态系统信息
	module.ModuleEcosystem = moduleEcosystem

	// 处理依赖项
	dependencies := make([]*baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem], 0)

	// 处理常规依赖
	for depName, depVersion := range packageJson.Dependencies {
		dependency := x.createDependency(depName, depVersion, false)
		dependencies = append(dependencies, dependency)
	}

	// 处理开发依赖
	for depName, depVersion := range packageJson.DevDependencies {
		dependency := x.createDependency(depName, depVersion, true)
		dependencies = append(dependencies, dependency)
	}

	module.Dependencies = dependencies

	// 将模块添加到项目中
	project.SetModule(packageJson.Name, module)

	// 设置项目生态系统信息
	projectEcosystem := &models.PackageLockProjectEcosystem{}
	// 可以添加更多项目特定的信息
	project.ProjectEcosystem = projectEcosystem

	return project, nil
}

// createDependency 创建一个表示依赖关系的对象
// 参数:
//   - name: 依赖的名称
//   - version: 依赖的版本
//   - isDev: 标识是否为开发依赖
//
// 返回:
//   - 依赖关系对象
func (x *PackageJsonParser) createDependency(name string, version string, isDev bool) *baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem] {
	dependency := &baseModels.ComponentDependency[*models.PackageLockComponentDependencyEcosystem]{}
	dependency.DependencyName = name
	dependency.DependencyVersion = version

	// 设置依赖生态系统信息
	ecosystem := &models.PackageLockComponentDependencyEcosystem{}
	if isDev {
		dev := true
		ecosystem.Dev = &dev
	}
	dependency.ComponentDependencyEcosystem = ecosystem

	return dependency
}

// parseComponent 解析组件信息
// 参数:
//   - packageName: 包名
//   - version: 版本号
//
// 返回:
//   - 组件对象
func (x *PackageJsonParser) parseComponent(packageName string, version string) *baseModels.Component[*models.PackageLockComponentEcosystem] {
	component := &baseModels.Component[*models.PackageLockComponentEcosystem]{}
	component.Name = packageName
	component.Version = version

	// 设置生态系统特定信息
	ecosystem := &models.PackageLockComponentEcosystem{}
	// 这里可以设置更多生态系统特定的信息
	component.ComponentEcosystem = ecosystem

	return component
}

// Close 关闭解析器并释放资源
// 实现了Parser接口的Close方法
// 当前实现不需要特殊的关闭步骤
func (x *PackageJsonParser) Close(ctx context.Context) error {
	// 对于当前的实现，没有需要关闭的资源
	return nil
}
