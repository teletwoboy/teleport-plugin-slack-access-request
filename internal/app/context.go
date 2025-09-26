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

package app

import (
	"context"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/util/container"
)

type Context struct {
	DB           *database.DB
	Clients      *container.Clients
	Server       *http.Server
	OtelShutdown func(context.Context) error
}

func NewContext() *Context {
	return &Context{}
}

func (app *Context) Cleanup(ctx context.Context) {
	if app.Server != nil {
		if err := app.Server.Shutdown(ctx); err != nil {
			slog.Error("Error closing HTTP server", "err", err)
		} else {
			slog.Info("successfully closed HTTP server")
		}
	}

	if app.DB != nil {
		if err := app.DB.Conn.Close(); err != nil {
			slog.Error("Error closing database connection", "err", err)
		} else {
			slog.Info("successfully closed database connection")
		}
	}

	if app.Clients != nil {
		if err := app.Clients.Teleport.Close(); err != nil {
			slog.Error("Error closing teleport client", "err", err)
		} else {
			slog.Info("successfully closed teleport client")
		}
	}

	if app.OtelShutdown != nil {
		err := app.OtelShutdown(ctx)
		if err != nil {
			slog.Error("Error closing Otel client", "err", err)
		} else {
			slog.Info("successfully closing Otel client")
		}
	}
}
