package accessreview

import (
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/util/container"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer(telemetry.WorkerAccessReview)

type Handler struct {
	Services *container.Services
}

func NewHandler(s *container.Services) *Handler {
	return &Handler{
		Services: s,
	}
}
