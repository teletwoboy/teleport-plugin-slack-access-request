package v1

import "github.com/go-chi/chi/v5"

type Handlers struct {
	AccessRequest *AccessRequestHandler
}

func Route(r chi.Router, h *Handlers) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/access-request", h.AccessRequest.HandleRequestModal)
	})
}
