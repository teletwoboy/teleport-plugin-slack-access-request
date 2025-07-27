package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func Run() {
	ctx := context.Background()

	db, err := database.Connect()
	if err != nil {
		slog.Error("failed to connect to Database", "err", err)
		return
	}
	slog.Info("successfully connected to database")
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing database connection", "err", err)
		}
	}(db.Conn)

	clients, err := container.NewClients(ctx)
	if err != nil {
		slog.Error("failed to initialize clients", "err", err)
		return
	}
	slog.Info("successfully initialized clients")

	repos := container.NewRepositories(db.Queries)
	services := container.NewServices(clients, repos)

	if err := services.SeedInit.Init(ctx, db, clients.Slack, clients.Teleport); err != nil {
		slog.Error("failed to initialize seed", "err", err)
	}

	slog.Info("starting event watching")
	event := NewEvent(db, clients, services)
	go event.StartWatcher(ctx)

	router := NewRouter(db, clients, repos, services)
	serve := router.Setup()

	slog.Info("starting server", "port", config.Cfg.Server.Port)
	if err := http.ListenAndServe(":"+config.Cfg.Server.Port, serve); err != nil {
		slog.Error("failed to start server", "err", err)
		return
	}
}
