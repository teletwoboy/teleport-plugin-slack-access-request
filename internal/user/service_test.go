package user

import (
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_MapUsersByUsername(t *testing.T) {
	service := &Service{}

	teleportUsers := []teleport.User{
		{Username: "alice@example.com"},
		{Username: "bob@example.com"},
		{Username: "carol@example.com"},
	}

	slackUsers := []slack.User{
		{
			ID:       "U123",
			Name:     "bob",
			RealName: "Bob Smith",
			Email:    "bob@example.com",
			Deleted:  false,
		},
		{
			ID:       "U124",
			Name:     "dave",
			RealName: "Dave Lee",
			Email:    "dave@example.com",
			Deleted:  false,
		},
	}

	users := service.MapUsersByUsername(slackUsers, teleportUsers)

	assert.Len(t, users, 1)
	assert.Equal(t, "bob@example.com", users[0].TeleportUser.Username)
	assert.Equal(t, "bob@example.com", users[0].SlackUser.Email)
}
