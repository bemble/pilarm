package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

func GetAppDir() (string, error) {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentFileDir := path.Dir(currentFilePath)
	return filepath.Abs(fmt.Sprintf(`%s/../..`, currentFileDir))
}

func GetRessourcesDir() string {
	dir, _ := GetAppDir()
	return filepath.Clean(fmt.Sprintf(`%s/ressources`, dir))
}

func GetRessourcePath(path string) string {
	dir := GetRessourcesDir()
	return fmt.Sprintf(`%s/%s`, dir, path)
}
