package parser

import (
	"context"
	"os"
	"path/filepath"
)

const PackageLockJsonFileName = "package-lock.json"

type PackageLockJsonParserInput struct {
	PackageLockJsonPath    string
	PackageLockJsonContent string
	ProjectRootDirectory   string
}

func (x *PackageLockJsonParserInput) Read(ctx context.Context) ([]byte, error) {

	if x.PackageLockJsonContent != "" {
		return []byte(x.PackageLockJsonContent), nil
	}

	if x.PackageLockJsonPath != "" {
		bytes, err := os.ReadFile(x.PackageLockJsonPath)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}

	packageJsonPath := filepath.Join(x.ProjectRootDirectory, PackageLockJsonFileName)
	bytes, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
