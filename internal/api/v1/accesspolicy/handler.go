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

package accesspolicy

import (
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util/container"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer(telemetry.APolicy)

type Handler struct {
	DB       *database.DB
	Clients  *container.Clients
	Repos    *container.Repositories
	Services *container.Services
}

func NewHandler(db *database.DB, c *container.Clients, r *container.Repositories, s *container.Services) *Handler {
	return &Handler{
		DB:       db,
		Clients:  c,
		Repos:    r,
		Services: s,
	}
}
