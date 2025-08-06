package check

import (
	"log/slog"
	"net/http"
)

// Healthz: Health handler
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	slog.Info("w", w, "OK")
}
