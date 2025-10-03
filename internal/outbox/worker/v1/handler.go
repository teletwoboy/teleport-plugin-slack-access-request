package v1

import (
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/outbox/worker/v1/accesspolicy"
	"teleport-plugin-slack-access-request/internal/outbox/worker/v1/accessrequest"
	"teleport-plugin-slack-access-request/internal/outbox/worker/v1/accessreview"
	"teleport-plugin-slack-access-request/internal/util/container"
)

type Handler struct {
	ARequest *accessrequest.Handler
	AReview  *accessreview.Handler
	APolicy  *accesspolicy.Handler
}

func NewHandler(db *database.DB, clients *container.Clients, srv *container.Services) *Handler {
	v1ARequest := accessrequest.NewHandler(db, clients, srv)
	v1AReview := accessreview.NewHandler(srv)
	v1APolicy := accesspolicy.NewHandler(db, clients, srv)
	return &Handler{
		ARequest: v1ARequest,
		AReview:  v1AReview,
		APolicy:  v1APolicy,
	}
}
