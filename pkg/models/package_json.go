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

	DevDependencies      Dependencies `json:"devDependencies"`
	Dependencies         Dependencies `json:"dependencies"`
	PeerDependencies     Dependencies `json:"peerDependencies"`
	OptionalDependencies Dependencies `json:"optionalDependencies"`
	BundledDependencies  []string     `json:"bundledDependencies"`
	BundleDependencies   []string     `json:"bundleDependencies"`

	License string `json:"license"`

	// 新增字段
	Keywords     []string `json:"keywords"`
	Author       Author   `json:"author"`
	Contributors []Author `json:"contributors"`
	Bugs         Bugs     `json:"bugs"`
	Homepage     string   `json:"homepage"`

	// 扩展repository字段为结构体
	Repository Repository `json:"repository"`

	Icon string `json:"icon"`

	GalleryBanner GalleryBanner `json:"galleryBanner"`

	// 添加其他常见字段
	Private       bool              `json:"private"`
	Bin           map[string]string `json:"bin"`
	Files         []string          `json:"files"`
	Man           []string          `json:"man"`
	Os            []string          `json:"os"`
	Cpu           []string          `json:"cpu"`
	Funding       Funding           `json:"funding"`
	Type          string            `json:"type"` // "module" 或 "commonjs"
	Workspaces    []string          `json:"workspaces"`
	Exports       interface{}       `json:"exports"` // 可以是字符串或复杂对象
	Imports       map[string]string `json:"imports"`
	EngineStrict  bool              `json:"engineStrict"`
	PreferGlobal  bool              `json:"preferGlobal"`
	PublishConfig PublishConfig     `json:"publishConfig"`
	Config        Config            `json:"config"`
}

// Author 表示作者或贡献者信息
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

// Bugs 表示问题追踪信息
type Bugs struct {
	Url   string `json:"url"`
	Email string `json:"email"`
}

// Repository 表示代码仓库信息
type Repository struct {
	Type      string `json:"type"`
	Url       string `json:"url"`
	Directory string `json:"directory"`
}

// Funding 表示资金支持信息
type Funding struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

// PublishConfig 发布配置
type PublishConfig struct {
	Registry  string `json:"registry"`
	Access    string `json:"access"`
	Tag       string `json:"tag"`
	Directory string `json:"directory"`
}

// Config 配置信息
type Config map[string]interface{}

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
