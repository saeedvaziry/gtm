package gtm

import (
	"html/template"
)

type ViteManifest struct {
	File string `json:"file"`
	Src  string `json:"src"`
}

type Vite struct {
	Manifest     map[string]ViteManifest
	BuildPath    string
	DevServerUrl string
}

type View struct {
	Loaded        bool
	Views         map[string]string
	Funcs         template.FuncMap
	BaseDir       string
	ComponentsDir string
	Extension     string
}
