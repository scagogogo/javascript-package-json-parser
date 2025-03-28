package parser

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/package-json-parser/pkg/models"
	baseModels "github.com/scagogogo/sca-base-module-components/pkg/models"
	"github.com/scagogogo/sca-base-module-ecosystem-parser/pkg/parser"
)

// YarnLockParser 用于解析yarn.lock文件
type YarnLockParser struct {
}

var _ parser.Parser[*YarnLockParserInput, *models.YarnLockProjectEcosystem, *models.YarnLockModuleEcosystem, *models.YarnLockComponentEcosystem, *models.YarnLockComponentDependencyEcosystem] = &YarnLockParser{}

func NewYarnLockParser() *YarnLockParser {
	return &YarnLockParser{}
}

const YarnLockParserName = "yarn-lock-parser"

func (x *YarnLockParser) GetName() string {
	return YarnLockParserName
}

func (x *YarnLockParser) Init(ctx context.Context) error {
	return nil
}

// Parse 解析yarn.lock文件
func (x *YarnLockParser) Parse(ctx context.Context, input *YarnLockParserInput) (*baseModels.Project[*models.YarnLockProjectEcosystem, *models.YarnLockModuleEcosystem, *models.YarnLockComponentEcosystem, *models.YarnLockComponentDependencyEcosystem], error) {
	// 读取yarn.lock文件
	yarnLockBytes, err := input.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read yarn.lock: %w", err)
	}

	// 解析yarn.lock文件
	yarnLock, moduleName, err := x.parseYarnLock(yarnLockBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yarn.lock: %w", err)
	}

	// 创建项目对象
	project := &baseModels.Project[*models.YarnLockProjectEcosystem, *models.YarnLockModuleEcosystem, *models.YarnLockComponentEcosystem, *models.YarnLockComponentDependencyEcosystem]{}
	project.Name = moduleName
	project.Version = "" // yarn.lock不包含版本信息，需要从package.json获取

	// 设置项目生态系统信息
	project.ProjectEcosystem = &models.YarnLockProjectEcosystem{}

	// 创建模块
	module := x.createModule(yarnLock, moduleName)
	project.SetModule(moduleName, module)

	return project, nil
}

// parseYarnLock 解析yarn.lock文件内容
func (x *YarnLockParser) parseYarnLock(data []byte) (*models.YarnLock, string, error) {
	// yarn.lock是一个非标准格式的文件，需要自定义解析
	yarnLock := &models.YarnLock{
		Dependencies: make(map[string]*models.YarnLockDependency),
	}

	// 模块名称暂时设为unknown，后续从解析过程中推断
	moduleName := "unknown"

	// 用正则表达式匹配依赖项
	lines := bytes.Split(data, []byte("\n"))

	// 解析内容
	var currentDep *models.YarnLockDependency
	var currentDepKeys []string
	indentLevel := 0

	// 依赖声明正则
	depRegex := regexp.MustCompile(`^"?([^"@]+)@([^"]+)"?:$`)

	// 版本正则
	versionRegex := regexp.MustCompile(`^\s+version "(.+)"$`)

	// 解析地址正则
	resolvedRegex := regexp.MustCompile(`^\s+resolved "(.+)"$`)

	// 完整性校验和正则
	integrityRegex := regexp.MustCompile(`^\s+integrity (.+)$`)

	// 依赖声明正则
	dependenciesRegex := regexp.MustCompile(`^\s+dependencies:$`)

	// 依赖项正则
	dependencyItemRegex := regexp.MustCompile(`^\s+([^@]+)@(.+) "(.+)"$`)

	for _, line := range lines {
		lineStr := string(line)

		// 跳过空行和注释
		if len(strings.TrimSpace(lineStr)) == 0 || strings.HasPrefix(strings.TrimSpace(lineStr), "#") {
			continue
		}

		// 寻找依赖声明
		if matches := depRegex.FindStringSubmatch(lineStr); len(matches) == 3 {
			// 完成前一个依赖的处理并开始新依赖
			if currentDep != nil && len(currentDepKeys) > 0 {
				for _, key := range currentDepKeys {
					yarnLock.Dependencies[key] = currentDep
				}
			}

			// 开始新依赖
			currentDep = &models.YarnLockDependency{
				Dependencies:         make(map[string]string),
				OptionalDependencies: make(map[string]string),
				PeerDependencies:     make(map[string]string),
			}

			// 获取依赖键
			depKey := strings.TrimSuffix(lineStr, ":")
			currentDepKeys = []string{depKey}

			// 尝试从依赖推断模块名称
			pkgName := matches[1]
			if moduleName == "unknown" && (pkgName != "npm" && pkgName != "yarn" && pkgName != "babel" && !strings.HasPrefix(pkgName, "@")) {
				moduleName = pkgName
			}

			indentLevel = 0
			continue
		}

		// 对于当前处理的依赖，解析其属性
		if currentDep != nil {
			// 版本
			if matches := versionRegex.FindStringSubmatch(lineStr); len(matches) == 2 {
				currentDep.Version = matches[1]
				continue
			}

			// 解析地址
			if matches := resolvedRegex.FindStringSubmatch(lineStr); len(matches) == 2 {
				currentDep.Resolved = matches[1]
				continue
			}

			// 完整性校验和
			if matches := integrityRegex.FindStringSubmatch(lineStr); len(matches) == 2 {
				currentDep.Integrity = matches[1]
				continue
			}

			// 处理依赖声明
			if dependenciesRegex.MatchString(lineStr) {
				indentLevel = len(lineStr) - len(strings.TrimLeft(lineStr, " "))
				continue
			}

			// 依赖项
			if matches := dependencyItemRegex.FindStringSubmatch(lineStr); len(matches) == 4 {
				currentIndent := len(lineStr) - len(strings.TrimLeft(lineStr, " "))
				if currentIndent > indentLevel {
					depName := matches[1]
					depVersion := matches[3]
					currentDep.Dependencies[depName] = depVersion
				}
				continue
			}
		}
	}

	// 处理最后一个依赖
	if currentDep != nil && len(currentDepKeys) > 0 {
		for _, key := range currentDepKeys {
			yarnLock.Dependencies[key] = currentDep
		}
	}

	return yarnLock, moduleName, nil
}

// createModule 创建模块对象
func (x *YarnLockParser) createModule(yarnLock *models.YarnLock, moduleName string) *baseModels.Module[*models.YarnLockModuleEcosystem, *models.YarnLockComponentEcosystem, *models.YarnLockComponentDependencyEcosystem] {
	module := &baseModels.Module[*models.YarnLockModuleEcosystem, *models.YarnLockComponentEcosystem, *models.YarnLockComponentDependencyEcosystem]{}
	module.Name = moduleName
	module.ModuleEcosystem = &models.YarnLockModuleEcosystem{}

	// 解析依赖
	dependencies := make([]*baseModels.ComponentDependency[*models.YarnLockComponentDependencyEcosystem], 0, len(yarnLock.Dependencies))

	// 依赖名称去重
	processedDeps := make(map[string]bool)

	for depKey, dep := range yarnLock.Dependencies {
		// 从depKey中提取实际包名
		// 例如 "lodash@^4.17.21" -> "lodash"
		parts := strings.Split(depKey, "@")
		var pkgName string
		if strings.HasPrefix(depKey, "@") && len(parts) >= 3 {
			// 处理作用域包，如 "@babel/core@^7.0.0"
			pkgName = "@" + parts[1] + "/" + parts[2]
		} else if len(parts) >= 2 {
			pkgName = parts[0]
		} else {
			continue
		}

		// 去重
		if processedDeps[pkgName+"@"+dep.Version] {
			continue
		}
		processedDeps[pkgName+"@"+dep.Version] = true

		// 创建依赖对象
		dependency := &baseModels.ComponentDependency[*models.YarnLockComponentDependencyEcosystem]{}
		dependency.DependencyName = pkgName
		dependency.DependencyVersion = dep.Version

		// 设置生态系统特定信息
		ecosystem := &models.YarnLockComponentDependencyEcosystem{}
		ecosystem.Resolved = dep.Resolved
		ecosystem.Integrity = dep.Integrity
		ecosystem.Source = dep.Source
		ecosystem.LanguageName = dep.LanguageName
		ecosystem.Bundled = dep.Bundled
		ecosystem.HasPeerDependencies = len(dep.PeerDependencies) > 0

		dependency.ComponentDependencyEcosystem = ecosystem

		dependencies = append(dependencies, dependency)
	}

	module.Dependencies = dependencies
	return module
}

func (x *YarnLockParser) Close(ctx context.Context) error {
	return nil
}
