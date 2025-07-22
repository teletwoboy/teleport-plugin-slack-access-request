package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/types/events"
)

type Service interface {
	EventPolling(ctx context.Context)
}

type API interface {
	SearchEvents(ctx context.Context, startTime, endTime time.Time, eventType string, fields []string, limit int, order types.EventOrder, token string) ([]events.AuditEvent, string, error)
}

type service struct {
	api API
}

func NewService(api API) Service {
	return &service{
		api: api,
	}
}

func (s *service) EventPolling(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastEventTime := time.Now().UTC()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, _, err := s.api.SearchEvents(ctx, lastEventTime, time.Now().UTC(), "", nil, 100, types.EventOrderAscending, "")
			if err != nil {
				slog.Error("Error searching events", "error", err.Error())
			}
			for _, event := range events {
				eventTime, _ := HandleEvent(event)
				if eventTime.After(lastEventTime) {
					lastEventTime = eventTime
				}
			}
		}
	}
}
