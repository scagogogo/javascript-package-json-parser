package models

// PackageLock package-lock.json文件对应的model
type PackageLock struct {
	Name            string                            `json:"name"`
	Version         string                            `json:"version"`
	LockFileVersion uint                              `json:"lockfileVersion"`
	Requires        *bool                             `json:"requires"`
	Dependencies    map[string]*PackageLockDependency `json:"dependencies"`
}

// PackageLockDependency package-lock.json中的依赖关系
type PackageLockDependency struct {
	Version   string       `json:"version"`
	Resolved  string       `json:"resolved"`
	Integrity string       `json:"integrity"`
	Dev       *bool        `json:"dev"`
	Requires  Dependencies `json:"requires"`
}
