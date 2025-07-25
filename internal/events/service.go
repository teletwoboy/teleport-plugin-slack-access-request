package events

import (
	"context"
	"log/slog"

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
	defer watcher.Close()

	slog.Info("Teleport EventWatcher started.", "kinds", watcher)

	for {
		select {
		case event := <-watcher.Events():
			s.eventHandler.TeleportEventHandle(ctx, event)
		case <-ctx.Done():
			slog.Info("EventWatcher context canceled, shutting down.")
			return
		case <-watcher.Done():
			return
		}
	}
}
