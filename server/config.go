package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

//go:embed config.default.toml
var defaultConfig string

var config Config

type Config struct {
	Port  int `toml:"port"`
	HTTPS struct {
		Enabled bool   `toml:"enabled"`
		Cert    string `toml:"cert"`
		Key     string `toml:"key"`
	} `toml:"https"`
}

func LoadConfig() {
	configData, err := os.ReadFile("config.toml")
	if err != nil && os.IsNotExist(err) {
		configData = []byte(defaultConfig)
		err := os.WriteFile("config.toml", []byte(defaultConfig), 0644)
		if err != nil {
			log.Panicln("Failed to create config.toml:", err)
		}
	} else if err != nil {
		log.Panicln("Failed to read config.toml:", err)
	}
	err = toml.Unmarshal(configData, &config)
	if err != nil {
		log.Panicln("Failed to parse config.toml:", err)
	}
}
