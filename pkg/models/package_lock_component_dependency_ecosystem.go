package models

type PackageLockComponentDependencyEcosystem struct {
	Resolved  string       `json:"resolved"`
	Integrity string       `json:"integrity"`
	Dev       *bool        `json:"dev"`
	Requires  Dependencies `json:"requires"`
}
