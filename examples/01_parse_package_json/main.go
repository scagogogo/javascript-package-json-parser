package main

import (
	"context"
	"fmt"
	"log"

	"github.com/scagogogo/package-json-parser/pkg/parser"
)

func main() {
	// 直接硬编码文件路径 - 不再需要命令行参数
	packageJsonPath := "./sample_package.json"

	// 创建解析器实例
	packageJsonParser := &parser.PackageJsonParser{}

	// 初始化解析器
	ctx := context.Background()
	err := packageJsonParser.Init(ctx)
	if err != nil {
		log.Fatalf("初始化解析器失败: %v\n", err)
	}
	defer packageJsonParser.Close(ctx)

	// 创建解析输入
	input := &parser.PackageJsonParserInput{
		PackageJsonPath: packageJsonPath,
	}

	// 解析文件
	project, err := packageJsonParser.Parse(ctx, input)
	if err != nil {
		log.Fatalf("解析package.json失败: %v\n", err)
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

		// 分类依赖
		var regularDeps, devDeps int
		for _, dep := range module.Dependencies {
			if dep.ComponentDependencyEcosystem != nil &&
				dep.ComponentDependencyEcosystem.Dev != nil &&
				*dep.ComponentDependencyEcosystem.Dev {
				devDeps++
			} else {
				regularDeps++
			}
		}

		fmt.Printf("  - 常规依赖: %d\n", regularDeps)
		fmt.Printf("  - 开发依赖: %d\n", devDeps)

		// 打印详细依赖列表
		if len(module.Dependencies) > 0 {
			fmt.Println("\n  依赖列表:")
			for _, dep := range module.Dependencies {
				depType := "常规"
				if dep.ComponentDependencyEcosystem != nil &&
					dep.ComponentDependencyEcosystem.Dev != nil &&
					*dep.ComponentDependencyEcosystem.Dev {
					depType = "开发"
				}
				fmt.Printf("  - %s (%s依赖): %s\n", dep.DependencyName, depType, dep.DependencyVersion)
			}
		}
	}
}
