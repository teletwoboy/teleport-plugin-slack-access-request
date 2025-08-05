package v1

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	AccessPolicyHandler *OpenAccessPolicyModalHandler
	AccessRoleHandler   *OpenAccessRoleModalHandler
	InteractionHandler  *InteractionHandler
}

func NewRouter(ap *OpenAccessPolicyModalHandler, ar *OpenAccessRoleModalHandler, i *InteractionHandler) *Router {
	return &Router{
		AccessPolicyHandler: ap,
		AccessRoleHandler:   ar,
		InteractionHandler:  i,
	}
}

func (r *Router) Route(router chi.Router) {
	router.Route("/v1", func(router chi.Router) {
		router.Post("/access-policy", r.AccessPolicyHandler.Handle)
		router.Post("/access-request", r.AccessRoleHandler.Handle)
		router.Post("/interaction", r.InteractionHandler.Handle)
	})
}
