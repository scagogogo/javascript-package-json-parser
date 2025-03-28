package models

// PackageLock package-lock.json文件对应的model
type PackageLock struct {
	Name            string                            `json:"name"`
	Version         string                            `json:"version"`
	LockFileVersion uint                              `json:"lockfileVersion"`
	Requires        *bool                             `json:"requires"`
	Dependencies    map[string]*PackageLockDependency `json:"dependencies"`

	// npm v7+ 新增字段
	Packages map[string]*PackageLockPackage `json:"packages"`

	// npm v7+ 新增字段，用于标识依赖树的根
	Workspaces map[string]interface{} `json:"workspaces"`
}

// PackageLockDependency package-lock.json中的依赖关系
type PackageLockDependency struct {
	Version   string       `json:"version"`
	Resolved  string       `json:"resolved"`
	Integrity string       `json:"integrity"`
	Dev       *bool        `json:"dev"`
	Requires  Dependencies `json:"requires"`

	// npm v5-v6 可能包含的字段
	Bundled  *bool `json:"bundled"`
	Optional *bool `json:"optional"`

	// npm v6+ 可能包含的嵌套依赖
	Dependencies map[string]*PackageLockDependency `json:"dependencies"`
}

// PackageLockPackage npm v7+ 引入的packages字段中的包定义
type PackageLockPackage struct {
	Version      string       `json:"version"`
	Resolved     string       `json:"resolved"`
	Integrity    string       `json:"integrity"`
	Dev          *bool        `json:"dev"`
	Requires     Dependencies `json:"requires"`
	Dependencies Dependencies `json:"dependencies"`

	// npm v7+ 特有字段
	Link             *bool             `json:"link"`
	Engines          map[string]string `json:"engines"`
	Os               []string          `json:"os"`
	Cpu              []string          `json:"cpu"`
	License          string            `json:"license"`
	Bin              map[string]string `json:"bin"`
	Funding          interface{}       `json:"funding"`
	DevOptional      *bool             `json:"devOptional"`
	InBundle         *bool             `json:"inBundle"`
	HasInstallScript *bool             `json:"hasInstallScript"`
}
