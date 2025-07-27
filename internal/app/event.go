package app

import (
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/events"
	v1 "teleport-plugin-slack-access-request/internal/events/v1"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func NewEvent(db *database.DB, c *container.Clients, s *container.Services) *events.Event {
	v1CUHandler := v1.NewCreateUserHandler(db, c, s)
	v1DUHandler := v1.NewDeleteUserHandler(db, c, s)
	event := events.NewEvent(v1CUHandler, v1DUHandler, s)
	return event
}
