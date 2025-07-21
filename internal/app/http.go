package app

import (
	"teleport-plugin-slack-access-request/internal/api"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"
)

func NewRouter(s *Services) *api.Router {
	v1ARHandler := v1.NewRequestAccessModalHandler(s.Slack, s.Teleport, s.User)
	v1IHandler := v1.NewInteractionHandler(s.Slack, s.Teleport, s.User)
	v1Router := v1.NewRouter(v1ARHandler, v1IHandler)
	return api.NewRouter(v1Router)
}
