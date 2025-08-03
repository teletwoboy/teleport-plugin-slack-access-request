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
	Server   ServerConfig
	Slack    SlackConfig
	Teleport TeleportConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:"8080"`
}

type SlackConfig struct {
	Token         string `envconfig:"SLACK_TOKEN" required:"true"`
	SigningSecret string `envconfig:"SLACK_SIGNING_SECRET" required:"true"`
}

type TeleportConfig struct {
	Addr         string `envconfig:"TELEPORT_ADDRESS" required:"true"`
	IdentityPath string `envconfig:"TELEPORT_IDENTITY_PATH" required:"true"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DATABASE_HOST" required:"true"`
	Port     string `envconfig:"DATABASE_PORT" required:"true"`
	Database string `envconfig:"DATABASE_NAME" required:"true"`
	Username string `envconfig:"DATABASE_USERNAME" required:"true"`
	Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	SslMode  string `envconfig:"DATABASE_SSL_MODE" default:"disable"`
}

var Cfg Config

func Init() {
	_ = godotenv.Load(envFileName)
	if err := envconfig.Process(envPrefix, &Cfg); err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	slog.Info("successfully loaded configs")
}
