package accessrequest

import (
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/util/container"
)

type Handler struct {
	DB       *database.DB
	Clients  *container.Clients
	Repos    *container.Repositories
	Services *container.Services
}

func NewHandler(db *database.DB, c *container.Clients, r *container.Repositories, s *container.Services) *Handler {
	return &Handler{
		DB:       db,
		Clients:  c,
		Repos:    r,
		Services: s,
	}
}
