package api

import (
	"net/http"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(h *v1.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)    // Logs all incoming HTTP requests
	r.Use(middleware.Recoverer) // Recovers from panics and prevents a server crash

	r.Route("/api", func(r chi.Router) {
		v1.Route(r, h)
	})

	return r
}
