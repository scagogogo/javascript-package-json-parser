package models

// YarnLock 表示yarn.lock文件的结构
type YarnLock struct {
	// yarn.lock文件被解析为依赖名称到依赖版本的映射
	// 注意：yarn.lock中的依赖键格式通常为 "package-name@^1.0.0"
	Dependencies map[string]*YarnLockDependency
}

// YarnLockDependency 表示yarn.lock中的单个依赖项
type YarnLockDependency struct {
	// 依赖版本，如 "1.2.3"
	Version string

	// 依赖解析地址
	Resolved string

	// 完整性校验和
	Integrity string

	// 依赖声明信息
	Dependencies map[string]string

	// 可选依赖
	OptionalDependencies map[string]string

	// 对等依赖
	PeerDependencies map[string]string

	// 依赖来源
	Source string

	// 语言类型
	LanguageName string

	// 是否捆绑依赖
	Bundled bool
}

// YarnLockProjectEcosystem 项目生态系统特定信息
type YarnLockProjectEcosystem struct {
	// 项目特定的yarn元数据
}

// YarnLockModuleEcosystem 模块生态系统特定信息
type YarnLockModuleEcosystem struct {
	// 模块特定的yarn元数据
}

// YarnLockComponentEcosystem 组件生态系统特定信息
type YarnLockComponentEcosystem struct {
	// 组件特定的yarn元数据
}

// YarnLockComponentDependencyEcosystem 组件依赖生态系统特定信息
type YarnLockComponentDependencyEcosystem struct {
	// 依赖特定元数据
	Resolved            string
	Integrity           string
	Source              string
	LanguageName        string
	Bundled             bool
	HasPeerDependencies bool
}
