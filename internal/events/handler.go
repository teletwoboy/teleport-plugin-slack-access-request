package events

import (
	"time"

	"github.com/gravitational/teleport/api/types/events"
)

func HandleEvent(event events.AuditEvent) time.Time {
	switch event.(type) {
	case *events.MFADeviceAdd:
		return event.GetTime()
	case *events.UserDelete:
		return event.GetTime()
	default:
		return event.GetTime()
	}
}
