package events

import (
	"context"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	teleportmodel "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/user"
	"teleport-plugin-slack-access-request/internal/util/verifier"
	"time"

	"github.com/gravitational/teleport/api/types"
)

type EventHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
	UserSrv     user.Service
}

func NewEventHandler(s slack.Service, t teleport.Service, u user.Service) *EventHandler {
	return &EventHandler{
		SlackSrv:    s,
		TeleportSrv: t,
		UserSrv:     u,
	}
}

func (h *EventHandler) TeleportEventHandle(ctx context.Context, event types.Event) {
	teleportVerifier := verifier.NewTeleport(h.TeleportSrv)
	switch event.Type {
	case types.OpPut:
		switch resource := event.Resource.(type) {
		case *types.UserV2:
			dupCheck, err := teleportVerifier.VerifyUserExistsByID(ctx, resource.GetName())
			if err != nil {
				slog.Error("failed to verify existing user", "err", err)
				return
			}
			if dupCheck {
				slog.Error("already exist user", "username", resource.GetName())
				return
			}
			h.createTeleportUser(ctx, resource.GetName())
		default:
			slog.Warn("EventWatcher: received OpPut for unhandled resource type",
				"kind", resource.GetKind(),
				"name", resource.GetName(),
				"type", fmt.Sprintf("%T", resource))
		}
	case types.OpDelete:
		switch resource := event.Resource.(type) {
		case *types.ResourceHeader:
			fmt.Printf("🗑️ 사용자 삭제: %s\n", resource.GetName())
		default:
			slog.Warn("EventWatcher: received OpDelete for unhandled resource type",
				"kind", resource.GetKind(),
				"name", resource.GetName(),
				"type", fmt.Sprintf("%T", resource))
		}
	case types.OpInit:
		slog.Info("EventWatcher: Initial stream established.")
	default:
		slog.Info("EventWatcher: Received unhandled event type", "type", event.Type, "kind", event.Resource.GetKind())
	}
}

func (h *EventHandler) createTeleportUser(ctx context.Context, email string) {
	info := teleportmodel.User{
		Username:   email,
		UseYn:      true,
		CreateCode: database.CreateCode,
		CreateDate: time.Now().UTC(),
	}
	_, err := h.TeleportSrv.CreateUser(ctx, info)
	if err != nil {
		slog.Error("failed", "err", err)
	}
}
