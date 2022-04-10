package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/spf13/viper"
)

var c Config

func Get() *Config {
	if !isInited() {
		log.Println("Reading configuration")

		confFile := fmt.Sprintf(`%s/config.json`, GetDataDir())
		_, err := os.Stat(confFile)
		if os.IsNotExist(err) {
			log.Println("Configuration file does not exists, creating an empty one")
			file, err := os.Create(confFile)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString("{}")
			defer file.Close()
		}

		viper.SetConfigFile(confFile)

		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		unmarshalErr := viper.Unmarshal(&c)
		if unmarshalErr != nil {
			panic(unmarshalErr)
		}
	}
	return &c
}

func isInited() bool {
	return !reflect.DeepEqual(c, Config{})
}

func GetAppDir() (string, error) {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentFileDir := path.Dir(currentFilePath)
	return filepath.Abs(fmt.Sprintf(`%s/../..`, currentFileDir))
}

func GetDataDir() string {
	dir, _ := GetAppDir()
	return filepath.Clean(fmt.Sprintf(`%s/data`, dir))
}

func GetRessourcesDir() string {
	dir, _ := GetAppDir()
	return filepath.Clean(fmt.Sprintf(`%s/ressources`, dir))
}

func GetRessourcePath(path string) string {
	dir := GetRessourcesDir()
	return fmt.Sprintf(`%s/%s`, dir, path)
}
