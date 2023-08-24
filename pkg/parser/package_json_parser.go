package parser

import (
	"context"
	"encoding/json"
	"github.com/scagogogo/package-json-parser/pkg/models"
	baseModels "github.com/scagogogo/sca-base-module-components/pkg/models"
	"github.com/scagogogo/sca-base-module-ecosystem-parser/pkg/parser"
)

type PackageJsonParser struct {
}

var _ parser.Parser[*PackageJsonParserInput, *models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem] = &PackageJsonParser{}

const PackageJsonParserName = "package-json-parser"

func (x *PackageJsonParser) GetName() string {
	return PackageJsonParserName
}

func (x *PackageJsonParser) Init(ctx context.Context) error {
	return nil
}

func (x *PackageJsonParser) Parse(ctx context.Context, input *PackageJsonParserInput) (*baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem], error) {

	packageJsonBytes, err := input.Read(ctx)
	if err != nil {
		return nil, err
	}

	packageJson := &models.PackageJson{}
	err = json.Unmarshal(packageJsonBytes, &packageJson)
	if err != nil {
		return nil, err
	}

	project := &baseModels.Project[*models.PackageLockProjectEcosystem, *models.PackageLockModuleEcosystem, *models.PackageLockComponentEcosystem, *models.PackageLockComponentDependencyEcosystem]{}
	project.Name = packageJson.Name
	project.Version = packageJson.Version

	return project, nil
}

func (x *PackageJsonParser) parseComponent(packageName string) {
	//registry := registry.NewRegistry()
	//information, err := registry.GetPackageInformation(context.Background(), packageName)
	//if err != nil {
	//	return err
	//}
	//for versionString, version := range information.Versions {
	//	version.Dependencies
	//}
}

func (x *PackageJsonParser) Close(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
