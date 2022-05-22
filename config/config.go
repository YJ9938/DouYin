package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

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
	Static struct {
		Path string
	}
}

var C Config

const ConfigPath = "config/config.yaml"

func Init() {
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&C); err != nil {
		log.Fatalf("failed to decode config file: %s", err)
	}
}
