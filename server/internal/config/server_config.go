package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type AppConfig struct {
	PlanetScaleConfig `koanf:"planetscale"`
	WeaviateConfig    `koanf:"weaviate"`
	OpenAIConfig      `koanf:"openai"`
	ServerConfig      `koanf:"server"`
}

type PlanetScaleConfig struct {
	Server       string `koanf:"server"`
	Protocol     string `koanf:"protocol"`
	Username     string `koanf:"username"`
	Password     string `koanf:"password"`
	DatabaseName string `koanf:"db"`
}

type WeaviateConfig struct {
	ApiKey string `koanf:"api_key"`
	Host   string `koanf:"host"`
	Scheme string `koanf:"scheme"`
}

type OpenAIConfig struct {
	ApiKey string `koanf:"api_key"`
}

type ServerConfig struct {
	Port int `koanf:"port"`
}

const LepoPrefix string = "LEPO_"

var config *AppConfig

func init() {
	k := koanf.New(".")

	if err := k.Load(file.Provider("app.toml"), toml.Parser()); err != nil {
		log.Println("info: app.toml not found. Looking from environment variables")
	}

	if err := k.Load(env.Provider(LepoPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, LepoPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return
	}

	var appConfig AppConfig
	err := k.Unmarshal("", &appConfig)
	if err != nil {
		return
	}

	config = &appConfig
}

func GetAppConfig() (*AppConfig, error) {
	if config != nil {
		return config, nil
	} else {
		err := fmt.Errorf("error: AppConfig not found")
		return nil, err
	}
}
