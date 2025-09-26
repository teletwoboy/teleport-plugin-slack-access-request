/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	Otel     OtelConfig
}

type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:"8080"`
}

type SlackConfig struct {
	Token                 string `envconfig:"SLACK_TOKEN" required:"true"`
	SigningSecret         string `envconfig:"SLACK_SIGNING_SECRET" required:"true"`
	DefaultNotifChannelID string `envconfig:"SLACK_DEFAULT_NOTIF_CHANNEL_ID" required:"true"`
}

type TeleportConfig struct {
	Addr         string `envconfig:"TELEPORT_ADDRESS" required:"true"`
	IdentityPath string `envconfig:"TELEPORT_IDENTITY_PATH" required:"true"`
}

type DatabaseConfig struct {
	Host        string `envconfig:"DATABASE_HOST" required:"true"`
	Port        string `envconfig:"DATABASE_PORT" required:"true"`
	Database    string `envconfig:"DATABASE_NAME" required:"true"`
	Username    string `envconfig:"DATABASE_USERNAME" required:"true"`
	Password    string `envconfig:"DATABASE_PASSWORD" required:"true"`
	SslMode     string `envconfig:"DATABASE_SSL_MODE" default:"disable"`
	SslRootCert string `envconfig:"DATABASE_SSL_ROOT_CERT" required:"false"`
}

type OtelConfig struct {
	Enable        bool    `envconfig:"OTEL_ENABLE" default:"false"`
	EndPoint      string  `envconfig:"OTEL_ENDPOINT" required:"false"`
	ServiceName   string  `envconfig:"OTEL_SERVICE_NAME" required:"false"`
	SamplingRatio float64 `envconfig:"OTEL_SAMPLING_RATIO" default:"1.0"`
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
