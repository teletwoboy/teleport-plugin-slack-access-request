package v1

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	AccessRequestHandler *AccessRequestHandler
	InteractionHandler   *InteractionHandler
}

func NewRouter(a *AccessRequestHandler, i *InteractionHandler) *Router {
	return &Router{
		AccessRequestHandler: a,
		InteractionHandler:   i,
	}
}

func (r *Router) Route(router chi.Router) {
	router.Route("/v1", func(router chi.Router) {
		router.Post("/access-request", r.AccessRequestHandler.HandleRequestModal)
		router.Post("/interaction", r.InteractionHandler.Handle)
	})
}
