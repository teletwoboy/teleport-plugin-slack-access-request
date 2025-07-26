package app

import (
	"teleport-plugin-slack-access-request/internal/api"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func NewRouter(s *container.Services) *api.Router {
	v1AccessRoleHandler := v1.NewOpenAccessRoleModalHandler(s.Slack, s.Teleport, s.User)
	v1IHandler := v1.NewInteractionHandler(s.Slack, s.Teleport, s.User)
	v1Router := v1.NewRouter(v1AccessRoleHandler, v1IHandler)
	return api.NewRouter(v1Router)
}
