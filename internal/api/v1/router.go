package v1

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	AccessRoleHandler  *OpenAccessRoleModalHandler
	InteractionHandler *InteractionHandler
}

func NewRouter(a *OpenAccessRoleModalHandler, i *InteractionHandler) *Router {
	return &Router{
		AccessRoleHandler:  a,
		InteractionHandler: i,
	}
}

func (r *Router) Route(router chi.Router) {
	router.Route("/v1", func(router chi.Router) {
		router.Post("/access-request", r.AccessRoleHandler.Handle)
		router.Post("/interaction", r.InteractionHandler.Handle)
	})
}
