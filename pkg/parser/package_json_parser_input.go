package parser

import (
	"context"
	"os"
	"path/filepath"
)

const PackageJsonFileName = "package.json"

type PackageJsonParserInput struct {
	PackageJsonPath      string
	PackageJsonContent   string
	ProjectRootDirectory string
}

func (x *PackageJsonParserInput) Read(ctx context.Context) ([]byte, error) {

	if x.PackageJsonContent != "" {
		return []byte(x.PackageJsonContent), nil
	}

	if x.PackageJsonPath != "" {
		bytes, err := os.ReadFile(x.PackageJsonPath)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}

	packageJsonPath := filepath.Join(x.ProjectRootDirectory, PackageJsonFileName)
	bytes, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
