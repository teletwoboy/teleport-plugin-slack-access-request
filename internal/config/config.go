package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	envFileName = ".env"
	envPrefix   = ""
)

type Config struct {
	Slack  SlackConfig
	Server ServerConfig
}

type SlackConfig struct {
	Token string `envconfig:"SLACK_TOKEN" required:"true"`
}

type ServerConfig struct {
	Port int `envconfig:"SERVER_PORT" default:"8080"`
}

var Cfg Config

func Init() {
	_ = godotenv.Load(envFileName)
	if err := envconfig.Process(envPrefix, &Cfg); err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	slog.Info("succeeded config loaded")
}
