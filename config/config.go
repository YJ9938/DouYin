package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Gobal configuration.
var C Config

type Config struct {
	MySQL struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}
	Redis struct {
		Host string
		Port string
		DB   int
	}
	JWT struct {
		ExpireMinutes int    `yaml:"expire_minutes"`
		SecretKey     string `yaml:"secret_key"`
	}
}

// The path of the configuration file.
const ConfigPath = "config/config.yaml"

// init parses the configuration file.
func init() {
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&C); err != nil {
		log.Fatalf("failed to decode config file: %s", err)
	}
}
