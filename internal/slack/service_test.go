package slack

import (
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockAPI struct{}

func (m *mockAPI) GetUsers(_ ...slack.GetUsersOption) ([]slack.User, error) {
	return []slack.User{
		{ID: "U1", Name: "user1", Deleted: false, Profile: slack.UserProfile{Email: "user1@example.com"}},
		{ID: "U2", Name: "user2", Deleted: true, Profile: slack.UserProfile{Email: "user2@example.com"}},
		{ID: "U3", Name: "user3", Deleted: false, Profile: slack.UserProfile{Email: "user3@example.com"}},
		{ID: "U4", Name: "user4", Deleted: false, Profile: slack.UserProfile{Email: "user4@example.com"}},
	}, nil
}

func (m *mockAPI) GetTeamInfo() (*slack.TeamInfo, error) {
	return &slack.TeamInfo{
		ID:   "T123456",
		Name: "Test Team",
	}, nil
}

func TestService_GetUsers(t *testing.T) {
	client := &Client{api: &mockAPI{}}
	service := &Service{client: client}

	users, err := service.GetUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 3)
	assert.Equal(t, users[0].ID, "U1")
}

func TestService_GetTeamInfo(t *testing.T) {
	client := &Client{api: &mockAPI{}}
	service := &Service{client: client}

	teamInfo, err := service.GetTeamInfo()

	assert.NoError(t, err)
	assert.Equal(t, "T123456", teamInfo.ID)
	assert.Equal(t, "Test Team", teamInfo.Name)
}
