package main

import (
	"context"
	"fmt"
	"log"

	"github.com/scagogogo/package-json-parser/pkg/parser"
)

func main() {
	// 直接硬编码文件路径 - 不再需要命令行参数
	yarnLockPath := "./sample_yarn.lock"

	// 创建解析器实例
	yarnLockParser := parser.NewYarnLockParser()

	// 初始化解析器
	ctx := context.Background()
	err := yarnLockParser.Init(ctx)
	if err != nil {
		log.Fatalf("初始化解析器失败: %v\n", err)
	}
	defer yarnLockParser.Close(ctx)

	// 创建解析输入
	input := &parser.YarnLockParserInput{
		YarnLockPath: yarnLockPath,
	}

	// 解析文件
	project, err := yarnLockParser.Parse(ctx, input)
	if err != nil {
		log.Fatalf("解析yarn.lock失败: %v\n", err)
	}

	// 打印解析结果
	fmt.Printf("项目名称: %s\n", project.Name)
	if project.Version != "" {
		fmt.Printf("项目版本: %s\n", project.Version)
	} else {
		fmt.Println("项目版本: 未知 (yarn.lock不包含版本信息)")
	}

	// 遍历项目模块
	fmt.Println("\n模块信息:")
	for _, module := range project.Modules {
		fmt.Printf("- 模块: %s", module.Name)
		if module.Version != "" {
			fmt.Printf("@%s", module.Version)
		}
		fmt.Println()

		// 打印依赖信息
		fmt.Printf("  依赖数量: %d\n", len(module.Dependencies))

		// 只打印前15个依赖以避免输出过长
		maxDeps := 15
		if len(module.Dependencies) > 0 {
			fmt.Printf("\n  依赖列表 (前%d个):\n", min(maxDeps, len(module.Dependencies)))
			for i, dep := range module.Dependencies {
				if i >= maxDeps {
					break
				}
				fmt.Printf("  - %s: %s\n", dep.DependencyName, dep.DependencyVersion)

				// 打印附加信息
				if dep.ComponentDependencyEcosystem != nil {
					if dep.ComponentDependencyEcosystem.Resolved != "" {
						fmt.Printf("    Resolved: %s\n", dep.ComponentDependencyEcosystem.Resolved)
					}
					if dep.ComponentDependencyEcosystem.Integrity != "" {
						fmt.Printf("    Integrity: %s\n", dep.ComponentDependencyEcosystem.Integrity)
					}
					if dep.ComponentDependencyEcosystem.HasPeerDependencies {
						fmt.Println("    有对等依赖")
					}
					if dep.ComponentDependencyEcosystem.Bundled {
						fmt.Println("    已打包")
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
