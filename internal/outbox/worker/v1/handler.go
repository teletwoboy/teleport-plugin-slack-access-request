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

package v1

import (
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/worker/v1/accesspolicy"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/worker/v1/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/worker/v1/accessreview"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util/container"
)

type Handler struct {
	ARequest *accessrequest.Handler
	AReview  *accessreview.Handler
	APolicy  *accesspolicy.Handler
}

func NewHandler(db *database.DB, clients *container.Clients, srv *container.Services) *Handler {
	v1ARequest := accessrequest.NewHandler(db, clients, srv)
	v1AReview := accessreview.NewHandler(srv)
	v1APolicy := accesspolicy.NewHandler(db, clients, srv)
	return &Handler{
		ARequest: v1ARequest,
		AReview:  v1AReview,
		APolicy:  v1APolicy,
	}
}
