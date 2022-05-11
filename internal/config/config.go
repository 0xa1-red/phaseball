package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Database DatabaseSettings
}

type DatabaseSettings struct {
	Host     string
	Port     string
	User     string
	Password string
	Db       string
}

var settings Settings

func Init(path string) error {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()

	decoder := yaml.NewDecoder(fp)
	if err := decoder.Decode(&settings); err != nil {
		return err
	}
	return nil
}

func Get() Settings {
	return settings
}
