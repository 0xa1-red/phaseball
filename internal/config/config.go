package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Service  ServiceSettings  `json:"service"`
	GameLog  GameLogSettings  `json:"game_log"`
	Database DatabaseSettings `json:"database"`
}

type ServiceSettings struct {
	Address string
}

type GameLogSettings struct {
	Kind string
	Path string
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

	settings.defaults()

	log.Println(settings.Service.Address)

	return nil
}

func Get() Settings {
	return settings
}

func (s *Settings) defaults() {
	if s.Service.Address == "" {
		s.Service.Address = "0.0.0.0:80"
	}
}
