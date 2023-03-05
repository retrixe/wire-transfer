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
	SupportDirectTransfers  bool `toml:"support_direct_transfers"`
	SupportProxiedTransfers bool `toml:"support_proxied_transfers"`
	AllowEncryption         bool `toml:"allow_encryption"`
	RequireEncryption       bool `toml:"require_encryption"`
	DefaultFileExpiryTime   int  `toml:"default_file_expiry_time"`
	MaxFileExpiryTime       int  `toml:"max_file_expiry_time"`
	Advanced                struct {
		UDPPort            int `toml:"udp_port"`
		UDPTimeoutDuration int `toml:"udp_timeout_duration"`
	} `toml:"advanced"`
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
