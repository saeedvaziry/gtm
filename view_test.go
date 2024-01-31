package gtm

import (
	"os"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	v := &View{}
	Init(v)

	// Test default values
	if v.BaseDir != "resources/views" {
		t.Errorf("Expected BaseDir to be 'resources/views', got %s", v.BaseDir)
	}
	if v.Extension != "html" {
		t.Errorf("Expected Extension to be 'html', got %s", v.Extension)
	}
	if v.ComponentsDir != "components" {
		t.Errorf("Expected ComponentsDir to be 'components', got %s", v.ComponentsDir)
	}

	// Test empty maps
	if v.Views == nil {
		t.Error("Expected Views to be initialized as an empty map")
	}
	if v.Funcs == nil {
		t.Error("Expected Funcs to be initialized as an empty map")
	}
}

func TestInitWithData(t *testing.T) {
	v := &View{
		BaseDir:       "test",
		Extension:     "test",
		ComponentsDir: "test",
		Views:         make(map[string]string),
		Funcs:         make(map[string]interface{}),
	}
	Init(v)

	// Test default values
	if v.BaseDir != "test" {
		t.Errorf("Expected BaseDir to be 'test', got %s", v.BaseDir)
	}
	if v.Extension != "test" {
		t.Errorf("Expected Extension to be 'test', got %s", v.Extension)
	}
	if v.ComponentsDir != "test" {
		t.Errorf("Expected ComponentsDir to be 'test', got %s", v.ComponentsDir)
	}

	// Test empty maps
	if v.Views == nil {
		t.Error("Expected Views to be initialized as an empty map")
	}
	if v.Funcs == nil {
		t.Error("Expected Funcs to be initialized as an empty map")
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a temporary file with a view content
	viewContent := "Test View Content"
	tempFile, err := os.CreateTemp(tempDir, "*.html")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.WriteString(viewContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Set up the view object
	v := &View{
		BaseDir:   tempDir,
		Extension: "html",
		Views:     make(map[string]string),
	}
	Init(v)

	// Call the Load function
	Load()

	// Verify that the view content is loaded correctly
	expectedKey := strings.ReplaceAll(tempFile.Name(), tempDir+"/", "")
	actualContent, ok := v.Views[expectedKey]
	if !ok {
		t.Errorf("Expected view with key '%s' to be loaded, but it was not found", expectedKey)
	}
	if actualContent != viewContent {
		t.Errorf("Expected view content '%s', but got '%s'", viewContent, actualContent)
	}
}

func TestRender(t *testing.T) {
	// Set up the view object
	v := &View{
		BaseDir:   "views",
		Extension: "html",
		Views: map[string]string{
			"test.html": "Test View Content",
		},
		Loaded: true,
	}
	Init(v)

	// Test rendering a view
	expected := "Test View Content"
	actual := Render("test.html", nil)
	if actual != expected {
		t.Errorf("Expected rendered view '%s', but got '%s'", expected, actual)
	}
}

func TestRenderWithData(t *testing.T) {
	// Set up the view object
	v := &View{
		BaseDir:   "views",
		Extension: "html",
		Views: map[string]string{
			"test.html": "Hello {{ .Name }}!",
		},
		Loaded: true,
	}
	Init(v)

	// Test rendering a view
	expected := "Hello Saeed!"
	actual := Render("test.html", map[string]string{"Name": "Saeed"})
	if actual != expected {
		t.Errorf("Expected rendered view '%s', but got '%s'", expected, actual)
	}
}

func TestRenderLayout(t *testing.T) {
	// Set up the view object
	v := &View{
		BaseDir:   "views",
		Extension: "html",
		Views: map[string]string{
			"layout.html": "This @child('content') is in the layout",
			"test.html":   "@layout('layout.html')\n@section('content') part @end",
		},
		Loaded: true,
	}
	Init(v)

	// Test rendering a view
	expected := "This  part  is in the layout"
	actual := Render("test.html", nil)
	if actual != expected {
		t.Errorf("Expected rendered view '%s', but got '%s'", expected, actual)
	}
}

func TestRenderComponent(t *testing.T) {
	// Set up the view object
	v := &View{
		BaseDir:       "views",
		Extension:     "html",
		ComponentsDir: "components",
		Views: map[string]string{
			"components/component.html": "This is a component",
			"test.html":                 "<x-component>hello</x-component>",
		},
		Loaded: true,
	}
	Init(v)

	// Test rendering a view
	expected := "This is a component"
	actual := Render("test.html", nil)
	if actual != expected {
		t.Errorf("Expected rendered view '%s', but got '%s'", expected, actual)
	}
}

func TestRenderIncludes(t *testing.T) {
	// Set up the view object
	v := &View{
		BaseDir:   "views",
		Extension: "html",
		Views: map[string]string{
			"include.html": "This is an include",
			"test.html":    "@include('include.html')",
		},
		Loaded: true,
	}
	Init(v)

	// Test rendering a view
	expected := "This is an include"
	actual := Render("test.html", nil)
	if actual != expected {
		t.Errorf("Expected rendered view '%s', but got '%s'", expected, actual)
	}
}
