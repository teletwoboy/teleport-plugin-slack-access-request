package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/logging"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	logging.Init()
	config.Init()
}

func main() {
	if err := run(); err != nil {
		slog.Error("Error occurred", "err", err)
	}
}

// run is temporary function
func run() error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing database connection", "err", err)
		}
	}(db.Conn)

	ctx := context.Background()

	_, err = slack.Init()
	if err != nil {
		return fmt.Errorf("error initializing slack client: %w", err)
	}

	_, err = teleport.Init(ctx)
	if err != nil {
		return fmt.Errorf("error initializing teleport client: %w", err)
	}

	http.HandleFunc("/register", func(_ http.ResponseWriter, _ *http.Request) {
		encrypted, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(encrypted))
	})

	log.Println(" Server Port : 8080")
	return http.ListenAndServe(":8080", nil)
}
