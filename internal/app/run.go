package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
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

	slackClt, err := slack.Init()
	if err != nil {
		slog.Error("failed to initialize slack client", "err", err)
		return
	}
	slog.Info("successfully initialized slack client")

	teleportClt, err := teleport.Init(ctx)
	if err != nil {
		slog.Error("failed to initialize teleport client", "err", err)
		return
	}
	slog.Info("successfully initialized teleport client")

	slackRepo := slack.NewRepository(db.Queries)
	teleportRepo := teleport.NewRepository(db.Queries)
	userRepo := user.NewRepository(db.Queries)
	seedInitRepo := seedinit.NewRepository(db.Queries)

	slackSrv := slack.NewService(slackClt, slackRepo)
	teleportSrv := teleport.NewService(teleportClt, teleportRepo)
	userSrv := user.NewService(userRepo, slackSrv, teleportSrv)
	seedInitSrv := seedinit.NewService(seedInitRepo, slackSrv, teleportSrv, userSrv)

	if err := seedInitSrv.Init(ctx, db, slackClt, teleportClt); err != nil {
		slog.Error("failed to initialize seed", "err", err)
	}

	v1ARHandler := v1.NewAccessRequestHandler(slackSrv, teleportSrv, userSrv)
	v1IHandler := v1.NewInteractionHandler(slackSrv, teleportSrv)
	v1Router := v1.NewRouter(v1ARHandler, v1IHandler)

	router := api.NewRouter(v1Router)
	serve := router.Setup()

	slog.Info("starting server", "port", config.Cfg.Server.Port)
	err = http.ListenAndServe(":"+config.Cfg.Server.Port, serve)
	if err != nil {
		slog.Error("failed to start server", "err", err)
		return
	}
}
