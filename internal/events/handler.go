package events

import (
	"context"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/database"

	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	teleportmodel "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/user"
	"teleport-plugin-slack-access-request/internal/util/verifier"
	"time"

	"github.com/gravitational/teleport/api/types/events"
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

func (h *EventHandler) TeleportEventHandle(ctx context.Context, event events.AuditEvent) time.Time {
	teleportVerifier := verifier.NewTeleport(h.TeleportSrv)
	switch e := event.(type) {
	case *events.MFADeviceAdd:
		dupCheck, err := teleportVerifier.VerifyUserExistsByID(ctx, e.User)
		if err != nil {
			slog.Error("failed to verify existing user", "err", err)
		}

		if dupCheck {
			slog.Error("already exist user", "username", e.User)
			return time.Time{}
		}
		h.createTeleportUser(ctx, e.User)
		return event.GetTime()
	case *events.UserDelete:
		return event.GetTime()
	default:
		return event.GetTime()
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
