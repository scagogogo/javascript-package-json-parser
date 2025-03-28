package main

import (
	"context"
	"fmt"
	"log"

	"github.com/scagogogo/package-json-parser/pkg/parser"
)

func main() {
	// 直接硬编码文件路径 - 不再需要命令行参数
	packageLockPath := "./sample_package_lock.json"

	// 创建解析器实例
	packageLockParser := parser.NewPackageLockParser()

	// 初始化解析器
	ctx := context.Background()
	err := packageLockParser.Init(ctx)
	if err != nil {
		log.Fatalf("初始化解析器失败: %v\n", err)
	}
	defer packageLockParser.Close(ctx)

	// 创建解析输入
	input := &parser.PackageLockJsonParserInput{
		PackageLockJsonPath: packageLockPath,
	}

	// 解析文件
	project, err := packageLockParser.Parse(ctx, input)
	if err != nil {
		log.Fatalf("解析package-lock.json失败: %v\n", err)
	}

	// 打印解析结果
	fmt.Printf("项目名称: %s\n", project.Name)
	fmt.Printf("项目版本: %s\n", project.Version)

	// 遍历项目模块
	fmt.Println("\n模块信息:")
	for _, module := range project.Modules {
		fmt.Printf("- 模块: %s@%s\n", module.Name, module.Version)

		// 打印依赖信息
		fmt.Printf("  依赖数量: %d\n", len(module.Dependencies))

		// 只打印前20个依赖以避免输出过长
		maxDeps := 20
		if len(module.Dependencies) > 0 {
			fmt.Printf("\n  依赖列表 (前%d个):\n", min(maxDeps, len(module.Dependencies)))
			for i, dep := range module.Dependencies {
				if i >= maxDeps {
					break
				}
				depType := "常规"
				if dep.ComponentDependencyEcosystem != nil &&
					dep.ComponentDependencyEcosystem.Dev != nil &&
					*dep.ComponentDependencyEcosystem.Dev {
					depType = "开发"
				}
				fmt.Printf("  - %s (%s依赖): %s\n", dep.DependencyName, depType, dep.DependencyVersion)

				// 打印附加信息
				if dep.ComponentDependencyEcosystem != nil {
					if dep.ComponentDependencyEcosystem.Resolved != "" {
						fmt.Printf("    Resolved: %s\n", dep.ComponentDependencyEcosystem.Resolved)
					}
					if dep.ComponentDependencyEcosystem.Integrity != "" {
						fmt.Printf("    Integrity: %s\n", dep.ComponentDependencyEcosystem.Integrity)
					}
				}
				fmt.Println()
			}

			if len(module.Dependencies) > maxDeps {
				fmt.Printf("  ... 还有 %d 个依赖未显示\n", len(module.Dependencies)-maxDeps)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
