package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"teleport-plugin-slack-access-request/internal/config"
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
	_, err := slack.Init()
	if err != nil {
		slog.Error("Error initializing slack client", "err", err)
		os.Exit(1)
	}

	_, err = teleport.Init()
	if err != nil {
		slog.Error("Error initializing teleport client", "err", err)
		os.Exit(1)
	}

	http.HandleFunc("/register", func(_ http.ResponseWriter, _ *http.Request) {
		encrypted, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(encrypted))
	})

	log.Println(" Server Port : 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
