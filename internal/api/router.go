package api

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/v1"
	"teleport-plugin-slack-access-request/internal/slack"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	v1 *v1.Router
}

func NewRouter(v1 *v1.Router) *Router {
	return &Router{
		v1: v1,
	}
}

func (r *Router) Setup(slackSrv slack.Service) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)    // Logs all incoming HTTP requests
	router.Use(middleware.Recoverer) // Recovers from panics and prevents a server crash

	router.With(VerifySlackRequest()).
		Route("/api", func(router chi.Router) {
			r.v1.Route(router)
		})

	return router
}
