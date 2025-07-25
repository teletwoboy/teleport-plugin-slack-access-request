package events

import (
	"context"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
	"teleport-plugin-slack-access-request/internal/util/verifier"

	"github.com/gravitational/teleport/api/types"
)

type Service interface {
	EventWatcher(ctx context.Context)
}

type API interface {
	NewWatcher(ctx context.Context, watch types.Watch) (types.Watcher, error)
}

type service struct {
	api          API
	eventHandler *EventHandler
}

type EventHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
	UserSrv     user.Service
}

type Router struct {
	EventService Service
}

func NewEventHandler(s slack.Service, t teleport.Service, u user.Service) *EventHandler {
	return &EventHandler{
		SlackSrv:    s,
		TeleportSrv: t,
		UserSrv:     u,
	}
}

func NewService(api API, h *EventHandler) Service {
	return &service{
		api:          api,
		eventHandler: h,
	}
}

func (s *service) EventWatcher(ctx context.Context) {
	watcher, err := s.api.NewWatcher(ctx, types.Watch{
		Kinds: []types.WatchKind{
			{Kind: types.KindUser},
		},
	})
	if err != nil {
		slog.Error("failed to create watcher", "err", err)
		return
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			slog.Error("failed to close watcher", "err", err)
		}
	}()

	slog.Info("Teleport EventWatcher started.", "kinds", watcher)

	for {
		select {
		case event := <-watcher.Events():
			teleportVerifier := verifier.NewTeleport(s.eventHandler.TeleportSrv)
			switch event.Type {
			case types.OpPut:
				switch resource := event.Resource.(type) {
				case *types.UserV2:
					var username string = resource.GetName()
					dupCheck, err := teleportVerifier.VerifyUserExistsByID(ctx, username)
					if err != nil {
						slog.Error("failed to verify existing user", "err", err)
						continue
					}
					if dupCheck {
						slog.Error("already exist user", "username", username)
						continue
					}
					createdTeleportUser := s.eventHandler.createTeleportUser(ctx, username)

					fetchedUsers, err := s.eventHandler.SlackSrv.FetchUsers()
					if err != nil {
						slog.Error("fetch failed", "err", err)
					}
					userInfo, checkExist := findByName(fetchedUsers, username)
					fmt.Println(userInfo)
					if !checkExist {
						slog.Error("not exist slack user", "err", checkExist)
						continue
					}
					createdSlackUser := s.eventHandler.createSlackUser(ctx, userInfo)

					s.eventHandler.createTotalUser(ctx, createdTeleportUser, createdSlackUser)
				default:
					slog.Warn("EventWatcher: received OpPut for unhandled resource type",
						"kind", resource.GetKind(),
						"name", resource.GetName(),
						"type", fmt.Sprintf("%T", resource))
				}
			case types.OpInit:
				slog.Info("EventWatcher: Initial stream established.")
			default:
				slog.Info("EventWatcher: Received unhandled event type", "type", event.Type, "kind", event.Resource.GetKind())
			}
		case <-ctx.Done():
			slog.Info("EventWatcher context canceled, shutting down.")
			return
		case <-watcher.Done():
			return
		}
	}
}

func findByName(users []models.User, targetName string) (models.User, bool) {
	for _, user := range users {
		if user.Name == targetName {
			return user, true
		}
	}
	return models.User{}, false
}
