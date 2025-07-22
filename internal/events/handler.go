package events

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gravitational/teleport/api/types/events"
	types "github.com/gravitational/teleport/api/types/events"
)

func HandleEvent(ctx context.Context, events []events.AuditEvent) (time.Time, error) {
	for _, event := range events {
		switch e := event.(type) {
		case *types.MFADeviceAdd:
			fmt.Print("Handling MFADeviceAdd event for user: ", e.User)
			return event.GetTime(), nil
		case *types.UserDelete:
			fmt.Printf("User %s deleted\n", e.Name)
			return event.GetTime(), nil
		default: // 이외의 이벤트는 처리하지 않음 -> 다음 이벤트로 넘어감
			slog.Warn("Unhandled event type", "event_type", event.GetType())
			return event.GetTime(), nil
		}
	}
	return time.Time{}, nil
}

// func routeEvent(ctx context.Context, event events.AuditEvent) (bool, error) {
// 	switch e := event.(type) {
// 	case *events.MFADeviceAdd:
// 		// MFA 추가 이벤트 처리
// 		fmt.Print("Handling MFADeviceAdd event for user: ", e.User)
// 		return true, nil
// 	case *events.UserDelete:
// 		// 사용자 삭제 이벤트 처리
// 		fmt.Printf("User %s deleted\n", e.Name)
// 		return true, nil
// 	default:
// 		return false, nil
// 	}
// }
