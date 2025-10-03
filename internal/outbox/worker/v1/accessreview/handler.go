package accessreview

import (
	"go.opentelemetry.io/otel"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/util/container"
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
