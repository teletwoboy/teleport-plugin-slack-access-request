package events

import (
	"context"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/database"
	slackmodel "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodel "teleport-plugin-slack-access-request/internal/teleport/models"
	usermodel "teleport-plugin-slack-access-request/internal/user/models"
	"time"
)

func (h *EventHandler) createSlackUser(ctx context.Context, userInfo slackmodel.User) *slackmodel.User {
	info := slackmodel.User{
		ID:         userInfo.ID,
		Name:       userInfo.Name,
		RealName:   userInfo.RealName,
		Email:      userInfo.Email,
		UseYn:      true,
		CreateCode: database.CreateCode,
		CreateDate: time.Now().UTC(),
	}

	user, err := h.SlackSrv.CreateUser(ctx, info)
	if err != nil {
		slog.Error("failed", "err", err)
	}

	return user
}

func (h *EventHandler) createTeleportUser(ctx context.Context, email string) *teleportmodel.User {
	info := teleportmodel.User{
		Username:   email,
		UseYn:      true,
		CreateCode: database.CreateCode,
		CreateDate: time.Now().UTC(),
	}
	user, err := h.TeleportSrv.CreateUser(ctx, info)
	if err != nil {
		slog.Error("failed", "err", err)
	}

	return user
}

func (h *EventHandler) createTotalUser(ctx context.Context, createdTeleportUser *teleportmodel.User, createdSlackUser *slackmodel.User) *usermodel.User {
	info := usermodel.User{
		TeleportUser: createdTeleportUser,
		SlackUser:    createdSlackUser,
		UseYn:        true,
		CreateCode:   database.CreateCode,
		CreateDate:   time.Now().UTC(),
	}
	user, err := h.UserSrv.CreateUser(ctx, info)
	if err != nil {
		slog.Error("failed", "err", err)
	}
	return user
}
