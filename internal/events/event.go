/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package events

import (
	"context"
	"fmt"
	"log/slog"
	v1 "teleport-plugin-slack-access-request/internal/events/v1"
	"teleport-plugin-slack-access-request/internal/util/container"

	"github.com/gravitational/teleport/api/types"
)

type Event struct {
	CreateUserHandler *v1.CreateUserHandler
	DeleteUserHandler *v1.DeleteUserHandler
	Services          *container.Services
}

func NewEvent(c *v1.CreateUserHandler, d *v1.DeleteUserHandler, s *container.Services) *Event {
	return &Event{
		CreateUserHandler: c,
		DeleteUserHandler: d,
		Services:          s,
	}
}

func (e *Event) StartWatcher(ctx context.Context) {
	watcher, err := e.Services.Teleport.NewWatcher(ctx, types.Watch{
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
	slog.Info("Teleport EventWatcher started")

	for {
		select {
		case event := <-watcher.Events():
			switch event.Type {
			case types.OpPut:
				switch resource := event.Resource.(type) {
				case *types.UserV2:
					slog.Info("EventWatcher: Received Teleport User Put event")
					e.CreateUserHandler.Handle(ctx, resource)
				default:
					slog.Warn("EventWatcher: received OpPut for unhandled resource type",
						"kind", resource.GetKind(),
						"name", resource.GetName(),
						"type", fmt.Sprintf("%T", resource))
				}
			case types.OpDelete:
				switch resource := event.Resource.(type) {
				case *types.ResourceHeader:
					slog.Info("EventWatcher: Received Teleport User Delete event")
					e.DeleteUserHandler.Handle(ctx, resource)
				default:
					slog.Warn("EventWatcher: received OpPut for unhandled resource type",
						"kind", resource.GetKind(),
						"name", resource.GetName(),
						"type", fmt.Sprintf("%T", resource))
				}
			case types.OpInit:
				slog.Info("EventWatcher: Initial stream established")
			default:
				slog.Info("EventWatcher: Received unhandled event type", "type", event.Type, "kind", event.Resource.GetKind())
			}
		case <-ctx.Done():
			slog.Info("EventWatcher context canceled, shutting down")
			return
		case <-watcher.Done():
			return
		}
	}
}
