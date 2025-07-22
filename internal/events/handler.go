package events

import (
	"context"
	"fmt"
	"time"

	"github.com/gravitational/teleport/api/types/events"
)

func HandleEvent(_ context.Context, event events.AuditEvent) (time.Time, error) {
	switch e := event.(type) {
	case *events.MFADeviceAdd:
		fmt.Print("Handling MFADeviceAdd event for user: ", e.User)
		return event.GetTime(), nil
	case *events.UserDelete:
		fmt.Printf("User %s deleted\n", e.Name)
		return event.GetTime(), nil
	default:
		return event.GetTime(), nil
	}
}
