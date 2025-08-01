package app

import (
	"teleport-plugin-slack-access-request/internal/api"
	v1 "teleport-plugin-slack-access-request/internal/api/v1"
	"teleport-plugin-slack-access-request/internal/api/v1/accesspolicy"
	"teleport-plugin-slack-access-request/internal/api/v1/accessrequest"
	"teleport-plugin-slack-access-request/internal/api/v1/accessreview"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func NewRouter(db *database.DB, c *container.Clients, r *container.Repositories, s *container.Services) *api.Router {
	v1OpenAPolicy := v1.NewOpenAccessPolicyModalHandler(s)
	v1OpenARequest := v1.NewOpenAccessRoleModalHandler(s)
	v1APolicy := accesspolicy.NewHandler(db, c, r, s)
	v1ARequest := accessrequest.NewHandler(db, c, r, s)
	v1AReview := accessreview.NewHandler(db, c, r, s)
	v1IHandler := v1.NewInteractionHandler(v1APolicy, v1ARequest, v1AReview, s)
	v1Router := v1.NewRouter(v1OpenAPolicy, v1OpenARequest, v1IHandler)
	return api.NewRouter(v1Router)
}
