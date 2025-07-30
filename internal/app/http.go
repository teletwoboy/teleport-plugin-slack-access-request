package app

import (
	"teleport-plugin-slack-access-request/internal/api"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func NewRouter(db *database.DB, c *container.Clients, r *container.Repositories, s *container.Services) *api.Router {
	v1APHandler := v1.NewOpenAccessPolicyModalHandler(s)
	v1ARHandler := v1.NewOpenAccessRoleModalHandler(s)
	v1IHandler := v1.NewInteractionHandler(db, c, r, s)
	v1Router := v1.NewRouter(v1APHandler, v1ARHandler, v1IHandler)
	return api.NewRouter(v1Router)
}
