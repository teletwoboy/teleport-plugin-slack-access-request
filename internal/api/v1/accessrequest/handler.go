package accessrequest

import "teleport-plugin-slack-access-request/internal/util/container"

type Handler struct {
	Services *container.Services
}

func NewHandler(services *container.Services) *Handler {
	return &Handler{
		Services: services,
	}
}
