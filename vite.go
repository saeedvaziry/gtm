package gtm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var vite *Vite

func ViteInit(v *Vite) {
	vite = v
	if vite.BuildPath == "" {
		vite.BuildPath = "static/build"
	}
	manifestFile, _ := os.ReadFile(fmt.Sprintf("%s/manifest.json", vite.BuildPath))
	manifest := map[string]ViteManifest{}
	json.Unmarshal(manifestFile, &manifest)
	vite.Manifest = manifest
	if vite.DevServerUrl == "" {
		vite.DevServerUrl = "http://localhost:5173"
	}
}

func ViteAsset(path string) string {
	devPath := fmt.Sprintf("%s/%s", vite.DevServerUrl, path)
	viteDevServer, err := http.Get(devPath)
	if err == nil && viteDevServer.StatusCode == 200 {
		return devPath
	}
	return fmt.Sprintf("/%s/%s", vite.BuildPath, vite.Manifest[path].File)
}
