package gtm

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var view *View

func Init(v *View) {
	view = v
	if view.BaseDir == "" {
		view.BaseDir = "resources/views"
	}
	if view.Extension == "" {
		view.Extension = "html"
	}
	if view.ComponentsDir == "" {
		view.ComponentsDir = "components"
	}
	if view.Views == nil {
		view.Views = make(map[string]string)
	}
	if view.Funcs == nil {
		view.Funcs = make(template.FuncMap)
	}
}

func Load() {
	err := filepath.Walk(view.BaseDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name()[len(info.Name())-len(view.Extension):] == view.Extension {
			content, err := os.ReadFile(p)
			if err != nil {
				panic(err)
			}
			view.Views[strings.ReplaceAll(p, view.BaseDir+"/", "")] = string(content)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	view.Loaded = true
}

func Render(name string, data interface{}) string {
	return parseView(getView(name), data)
}

func AddFunc(name string, fn interface{}) {
	view.Funcs[name] = fn
}

func parseView(v string, data interface{}) string {
	v = parseLayout(v)
	v = parseIncludes(v)
	v = parseComponents(v)
	return executeTemplate(v, data)
}

func parseLayout(v string) string {
	layoutPattern := `@layout\(["'](.+?)["']\)`
	re := regexp.MustCompile(layoutPattern)
	matches := re.FindStringSubmatch(v)
	if len(matches) > 1 {
		layout := getView(matches[1])
		children := extractChildren(layout)
		for _, child := range children {
			section := extractSection(v, child)
			layout = strings.ReplaceAll(layout, fmt.Sprintf(`@child("%s")`, child), section)
			layout = strings.ReplaceAll(layout, fmt.Sprintf(`@child('%s')`, child), section)
		}
		return layout
	}
	return v
}

func parseIncludes(v string) string {
	includePattern := `@include\(["'](.+?)["']\)`
	re := regexp.MustCompile(includePattern)
	matches := re.FindAllStringSubmatch(v, -1)
	for _, match := range matches {
		include := getView(match[1])
		v = strings.ReplaceAll(v, match[0], include)
	}
	return v
}

func parseComponents(v string) string {
	re := regexp.MustCompile(`<x-([^>\s]+)[^>]*>([\s\S]*?)<\/x-(.*?)>`)
	matches := re.FindAllStringSubmatch(v, -1)
	for _, match := range matches {
		name := match[1]
		content := match[2]
		component := getView(fmt.Sprintf("%s/%s.%s", view.ComponentsDir, name, view.Extension))
		component = strings.ReplaceAll(component, "@slot", content)
		v = strings.ReplaceAll(v, match[0], component)
	}
	return v
}

func extractChildren(v string) []string {
	childrenPattern := `@child\(["'](.+?)["']\)`
	re := regexp.MustCompile(childrenPattern)
	matches := re.FindAllStringSubmatch(v, -1)
	children := make([]string, len(matches))
	for i, match := range matches {
		children[i] = match[1]
	}
	return children
}

func extractSection(v string, section string) string {
	sectionPattern := fmt.Sprintf(`@section\(['"]%s['"]\)([\s\S]*?)@end`, section)
	re := regexp.MustCompile(sectionPattern)
	matches := re.FindStringSubmatch(v)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func executeTemplate(v string, data interface{}) string {
	pattern := `\{\{(.+?)\}\}`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(v, -1)
	expressions := make([]string, len(matches))
	for i, match := range matches {
		expressions[i] = match[1]
	}
	tmpl, err := template.New("template").Funcs(view.Funcs).Funcs(getFunctionsFromData(data)).Parse(v)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func getFunctionsFromData(data interface{}) template.FuncMap {
	functions := make(template.FuncMap)
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return functions
	}
	for k, v := range dataMap {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			functions[k] = v
		}
	}
	return functions
}

func getView(name string) string {
	if view.Loaded {
		return view.Views[name]
	}
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", view.BaseDir, name))
	if err != nil {
		panic(err)
	}
	return string(file)
}
