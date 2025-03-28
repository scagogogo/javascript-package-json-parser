package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/package-json-parser/pkg/parser"
)

func main() {
	// 直接硬编码项目目录路径 - 不再需要命令行参数
	projectDir := "./sample_project"

	// 检查目录是否存在
	info, err := os.Stat(projectDir)
	if os.IsNotExist(err) || !info.IsDir() {
		log.Fatalf("错误: %s 不是有效的目录\n", projectDir)
	}

	// 初始化上下文
	ctx := context.Background()

	// 解析 package.json
	packageJsonPath := filepath.Join(projectDir, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		fmt.Println("=== 解析 package.json ===")
		parsePackageJson(ctx, packageJsonPath)
		fmt.Println()
	}

	// 解析 package-lock.json
	packageLockPath := filepath.Join(projectDir, "package-lock.json")
	if _, err := os.Stat(packageLockPath); err == nil {
		fmt.Println("=== 解析 package-lock.json ===")
		parsePackageLock(ctx, packageLockPath)
		fmt.Println()
	}

	// 解析 yarn.lock
	yarnLockPath := filepath.Join(projectDir, "yarn.lock")
	if _, err := os.Stat(yarnLockPath); err == nil {
		fmt.Println("=== 解析 yarn.lock ===")
		parseYarnLock(ctx, yarnLockPath)
		fmt.Println()
	}

	// 如果没有找到任何文件
	_, err1 := os.Stat(packageJsonPath)
	_, err2 := os.Stat(packageLockPath)
	_, err3 := os.Stat(yarnLockPath)
	if os.IsNotExist(err1) && os.IsNotExist(err2) && os.IsNotExist(err3) {
		log.Fatalf("在指定目录中未找到任何可解析的文件 (package.json, package-lock.json, yarn.lock)")
	}
}

// 解析 package.json 文件
func parsePackageJson(ctx context.Context, path string) {
	packageJsonParser := &parser.PackageJsonParser{}
	err := packageJsonParser.Init(ctx)
	if err != nil {
		fmt.Printf("初始化 package.json 解析器失败: %v\n", err)
		return
	}
	defer packageJsonParser.Close(ctx)

	input := &parser.PackageJsonParserInput{
		PackageJsonPath: path,
	}

	project, err := packageJsonParser.Parse(ctx, input)
	if err != nil {
		fmt.Printf("解析 package.json 失败: %v\n", err)
		return
	}

	// 打印基本信息
	fmt.Printf("项目名称: %s\n", project.Name)
	fmt.Printf("项目版本: %s\n", project.Version)

	// 打印依赖数量
	for _, module := range project.Modules {
		fmt.Printf("模块 '%s' 包含 %d 个依赖\n", module.Name, len(module.Dependencies))

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

		fmt.Printf("- 常规依赖: %d\n", regularDeps)
		fmt.Printf("- 开发依赖: %d\n", devDeps)
	}
}

// 解析 package-lock.json 文件
func parsePackageLock(ctx context.Context, path string) {
	packageLockParser := parser.NewPackageLockParser()
	err := packageLockParser.Init(ctx)
	if err != nil {
		fmt.Printf("初始化 package-lock.json 解析器失败: %v\n", err)
		return
	}
	defer packageLockParser.Close(ctx)

	input := &parser.PackageLockJsonParserInput{
		PackageLockJsonPath: path,
	}

	project, err := packageLockParser.Parse(ctx, input)
	if err != nil {
		fmt.Printf("解析 package-lock.json 失败: %v\n", err)
		return
	}

	// 打印基本信息
	fmt.Printf("项目名称: %s\n", project.Name)
	fmt.Printf("项目版本: %s\n", project.Version)

	// 打印依赖数量
	for _, module := range project.Modules {
		fmt.Printf("模块 '%s' 包含 %d 个锁定依赖\n", module.Name, len(module.Dependencies))
	}
}

// 解析 yarn.lock 文件
func parseYarnLock(ctx context.Context, path string) {
	yarnLockParser := parser.NewYarnLockParser()
	err := yarnLockParser.Init(ctx)
	if err != nil {
		fmt.Printf("初始化 yarn.lock 解析器失败: %v\n", err)
		return
	}
	defer yarnLockParser.Close(ctx)

	input := &parser.YarnLockParserInput{
		YarnLockPath: path,
	}

	project, err := yarnLockParser.Parse(ctx, input)
	if err != nil {
		fmt.Printf("解析 yarn.lock 失败: %v\n", err)
		return
	}

	// 打印基本信息
	fmt.Printf("项目名称: %s\n", project.Name)
	if project.Version != "" {
		fmt.Printf("项目版本: %s\n", project.Version)
	} else {
		fmt.Println("项目版本: 未知 (yarn.lock不包含版本信息)")
	}

	// 打印依赖数量
	for _, module := range project.Modules {
		fmt.Printf("模块 '%s' 包含 %d 个锁定依赖\n", module.Name, len(module.Dependencies))
	}
}
