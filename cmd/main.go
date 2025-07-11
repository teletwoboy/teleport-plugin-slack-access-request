package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/log"
	"teleport-plugin-slack-access-request/internal/slack"

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

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		encrypted, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(encrypted))
	})

	log.Println(" Server Port : 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
