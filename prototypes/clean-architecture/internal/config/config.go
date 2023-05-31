package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type AppConfig struct {
	DBConfig     `koanf:"database"`
	ServerConfig struct {
		Port int `koanf:"port"`
	} `koanf:"server"`
}

type DBConfig struct {
	Server       string `koanf:"server"`
	Protocol     string `koanf:"protocol"`
	Username     string `koanf:"username"`
	Password     string `koanf:"password"`
	DatabaseName string `koanf:"db"`
}

func GetAppConfig() (*AppConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider("app.toml"), toml.Parser()); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider("LEPO_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "LEPO_")), "_", ".", -1)
	}), nil); err != nil {
		log.Printf("Error loading config from environment %v \n", err)
		return nil, err
	}

	var appConfig AppConfig
	err := k.Unmarshal("", &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
