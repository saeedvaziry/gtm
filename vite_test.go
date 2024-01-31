package gtm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestViteInit(t *testing.T) {
	v := &Vite{}
	ViteInit(v)

	// Test default build path
	expectedBuildPath := "static/build"
	if v.BuildPath != expectedBuildPath {
		t.Errorf("Expected BuildPath to be %s, but got %s", expectedBuildPath, v.BuildPath)
	}

	// Test reading manifest file
	expectedManifest := map[string]ViteManifest{}
	manifestFile, _ := os.ReadFile(fmt.Sprintf("%s/manifest.json", v.BuildPath))
	json.Unmarshal(manifestFile, &expectedManifest)
	if len(v.Manifest) != len(expectedManifest) {
		t.Errorf("Expected Manifest to have %d entries, but got %d", len(expectedManifest), len(v.Manifest))
	}
}

func TestViteAsset(t *testing.T) {
	v := &Vite{}

	ViteInit(v)

	// Test asset from dev server
	mockDevServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockDevServer.Close()
	v.DevServerUrl = mockDevServer.URL
	expectedDevPath := fmt.Sprintf("%s/%s", mockDevServer.URL, "path/to/asset.js")
	if result := ViteAsset("path/to/asset.js"); result != expectedDevPath {
		t.Errorf("Expected ViteAsset to return %s, but got %s", expectedDevPath, result)
	}

	// Test asset from build path
	v.BuildPath = "static/build"
	v.Manifest = map[string]ViteManifest{
		"path/to/asset.js": {
			File: "asset.123456.js",
		},
	}
	expectedBuildPath := fmt.Sprintf("%s/%s", v.DevServerUrl, "path/to/asset.js")
	if result := ViteAsset("path/to/asset.js"); result != expectedBuildPath {
		t.Errorf("Expected ViteAsset to return %s, but got %s", expectedBuildPath, result)
	}
}
