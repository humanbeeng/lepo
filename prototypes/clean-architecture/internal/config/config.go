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

const LepoPrefix string = "LEPO_"

func GetAppConfig() (*AppConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider("app.toml"), toml.Parser()); err != nil {
		log.Println("info: app.toml not found. Looking from environment variables")
	}

	if err := k.Load(env.Provider(LepoPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, LepoPrefix)), "_", ".", -1)
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
