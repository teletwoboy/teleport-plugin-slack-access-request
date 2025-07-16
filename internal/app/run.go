package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/integration"
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"

	"golang.org/x/crypto/bcrypt"
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

	seedInitRepo := seedinit.NewRepository(db)
	seedInitSrv := seedinit.NewService(seedInitRepo)
	err = seedInitSrv.Create(ctx)
	if err != nil {
		slog.Error("failed to create seedinit row", "err", err)
	}

	status, err := seedInitSrv.GetStatus(ctx)
	if err != nil {
		slog.Error("failed to get seedinit status", "err", err)
	}

	if status != "initialized" {
		err := initializeSeed(ctx, db, slackClt, teleportClt)
		if err != nil {
			slog.Error("failed to initialize seed", "err", err)
			return
		}
		slog.Info("successfully initialized seed")
	} else {
		slog.Info("already initialized seed")
	}

	http.HandleFunc("/register", func(_ http.ResponseWriter, _ *http.Request) {
		encrypted, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(encrypted))
	})

	slog.Info(" Server Port : 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("failed to start server", "err", err)
		return
	}
}

func initializeSeed(ctx context.Context, db *database.DB, slackClt *slack.Client, teleportClt *teleport.Client) error {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", "err", err)
		}
	}(tx)

	qtx := db.Queries.WithTx(tx)
	integrationSrv := integration.NewServiceWithTx(qtx, slackClt, teleportClt)
	mappedUsers, err := integrationSrv.FetchAndMapUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mapped users: %w", err)
	}

	err = integrationSrv.ProvisionUsers(ctx, mappedUsers)
	if err != nil {
		return fmt.Errorf("failed to sync users: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
