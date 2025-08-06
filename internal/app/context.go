package app

import (
	"context"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/util/container"
)

type Context struct {
	DB      *database.DB
	Clients *container.Clients
	Server  *http.Server
}

func NewContext(db *database.DB, c *container.Clients, s *http.Server) *Context {
	return &Context{
		DB:      db,
		Clients: c,
		Server:  s,
	}
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
}
