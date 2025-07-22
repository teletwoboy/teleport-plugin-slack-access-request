package events

import (
	"context"
	"fmt"

	// "teleport-plugin-slack-access-request/internal/teleport"
	"time"

	"github.com/gravitational/teleport/api/types"
	auditeventtypes "github.com/gravitational/teleport/api/types/events"
)

type TeleportEventProcessor interface {
	SearchEvents(ctx context.Context, startTime, endTime time.Time, eventType string, fields []string, limit int, order types.EventOrder, token string) ([]auditeventtypes.AuditEvent, string, error)
}

type Processor struct {
	TeleportEventProcessor TeleportEventProcessor
}

func NewService(teleportclient TeleportEventProcessor) TeleportEventProcessor {
	return &Processor{TeleportEventProcessor: teleportclient}
}

func (p *Processor) EventPolling(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastEventTime := time.Now().UTC()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, _, err := p.TeleportEventProcessor.SearchEvents(ctx, lastEventTime, time.Now().UTC(), "", nil, 100, types.EventOrderAscending, "")
			if err != nil {
				fmt.Printf("Error searching events: %v\n", err)
			}
			eventTime, _ := HandleEvent(ctx, events)
			if eventTime.After(lastEventTime) {
				lastEventTime = eventTime
			}
		}
	}
}

func (p *Processor) SearchEvents(ctx context.Context, start, end time.Time, eventType string, fields []string, limit int, order types.EventOrder, token string) ([]auditeventtypes.AuditEvent, string, error) {
	return p.TeleportEventProcessor.SearchEvents(ctx, start, end, eventType, fields, limit, order, token)
}
