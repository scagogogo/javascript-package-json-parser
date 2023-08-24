package models

// PackageJson 跟package.json文件结构对应
type PackageJson struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	DisplayName      string      `json:"displayName"`
	Description      string      `json:"description"`
	Publisher        string      `json:"publisher"`
	Engines          Engines     `json:"engines"`
	Categories       []string    `json:"categories"`
	ActivationEvents []string    `json:"activationEvents"`
	Main             string      `json:"main"`
	Contributes      Contributes `json:"contributes"`

	Scripts Scripts `json:"scripts"`

	DevDependencies Dependencies `json:"devDependencies"`
	Dependencies    Dependencies `json:"dependencies"`

	License string `json:"license"`

	Repository string `json:"repository"`

	Icon string `json:"icon"`

	GalleryBanner GalleryBanner `json:"galleryBanner"`
}

type Engines struct {
	Vscode string `json:"vscode"`
}

type Task struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Args struct {
	Type        string        `json:"type"`
	Default     []interface{} `json:"default"`
	Description string        `json:"description"`
}

type Properties struct {
	Task Task `json:"task"`
	Args Args `json:"args"`
}

type TaskDefinitions struct {
	Type       string     `json:"type"`
	Required   []string   `json:"required"`
	Properties Properties `json:"properties"`
}

type Contributes struct {
	TaskDefinitions []TaskDefinitions `json:"taskDefinitions"`
}

type Scripts struct {
	VscodePrepublish string `json:"vscode:prepublish"`
	Compile          string `json:"compile"`
	Watch            string `json:"watch"`
	Build            string `json:"build"`
	Postinstall      string `json:"postinstall"`
}

type GalleryBanner struct {
	Color string `json:"color"`
	Theme string `json:"theme"`
}
