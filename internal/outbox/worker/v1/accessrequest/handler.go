package accessrequest

import (
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/util/container"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer(telemetry.WorkerAccessRequest)

type Handler struct {
	DB       *database.DB
	Clients  *container.Clients
	Services *container.Services
}

func NewHandler(db *database.DB, c *container.Clients, s *container.Services) *Handler {
	return &Handler{
		DB:       db,
		Clients:  c,
		Services: s,
	}
}
