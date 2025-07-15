package app

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/integration"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
)

func Run() {
	ctx := context.Background()

	db, err := connectDB()
	if err != nil {
		slog.Error("failed to connect to Database", "err", err)
		os.Exit(1)
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
		os.Exit(1)
	}
	slog.Info("successfully initialized slack client")

	teleportClt, err := teleport.Init(ctx)
	if err != nil {
		slog.Error("failed to initialize teleport client", "err", err)
		os.Exit(1)
	}
	slog.Info("successfully initialized teleport client")

	integrationSrv := integration.NewService(slackClt, teleportClt)

	err = integrationSrv.SyncUsers(ctx, db)
	if err != nil {
		slog.Error("failed to sync users", "err", err)
		os.Exit(1)
	}
	slog.Info("successfully synced users")

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
		os.Exit(1)
	}
}

func connectDB() (*database.DB, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return db, nil
}
