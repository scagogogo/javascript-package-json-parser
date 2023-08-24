package models

type PackageLockModuleEcosystem struct {
	LockFileVersion uint  `json:"lockfileVersion"`
	Requires        *bool `json:"requires"`

	// package-lock.json的原始内容
	PackageLockContent string `json:"package_lock_content"`
}
