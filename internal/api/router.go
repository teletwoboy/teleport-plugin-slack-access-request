package api

import (
	"net/http"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	v1 *v1.Router
}

func NewRouter(v1 *v1.Router) *Router {
	return &Router{
		v1: v1,
	}
}

func (r *Router) Setup(router *chi.Mux) http.Handler {
	router.With(VerifySlackRequest()).
		Route("/api", func(router chi.Router) {
			r.v1.Route(router)
		})
	return router
}
