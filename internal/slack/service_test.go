package slack

import (
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockSlackAPI struct{}

func (m *mockSlackAPI) GetUsers(_ ...slack.GetUsersOption) ([]slack.User, error) {
	return []slack.User{
		{ID: "U1", Name: "user1", Deleted: false, Profile: slack.UserProfile{Email: "user1@example.com"}},
		{ID: "U2", Name: "user2", Deleted: true, Profile: slack.UserProfile{Email: "user2@example.com"}},
		{ID: "U3", Name: "user3", Deleted: false, Profile: slack.UserProfile{Email: "user3@example.com"}},
		{ID: "U4", Name: "user4", Deleted: false, Profile: slack.UserProfile{Email: "user4@example.com"}},
	}, nil
}

func TestService_GetUsers(t *testing.T) {
	client := &Client{api: &mockSlackAPI{}}
	service := &Service{client: client}

	users, err := service.GetUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 3)
	assert.Equal(t, users[0].ID, "U1")
}
