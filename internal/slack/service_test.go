package slack

import (
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
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

func (m *mockAPI) GetConversations(_ *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	conversation1 := slack.Conversation{ID: "C1234"}
	conversation2 := slack.Conversation{ID: "C2345"}
	conversation3 := slack.Conversation{ID: "C3456"}

	groupConversation1 := slack.GroupConversation{Conversation: conversation1, Name: "backend-reviewers"}
	groupConversation2 := slack.GroupConversation{Conversation: conversation2, Name: "frontend-reviewers"}
	groupConversation3 := slack.GroupConversation{Conversation: conversation3, Name: "random"}

	return []slack.Channel{
		{GroupConversation: groupConversation1, IsMember: true},
		{GroupConversation: groupConversation2, IsMember: true},
		{GroupConversation: groupConversation3, IsMember: true},
	}, "", nil
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

func TestService_GetReviewersChannels(t *testing.T) {
	client := &Client{api: &mockAPI{}}
	service := &Service{client: client}

	reviewersChannels, err := service.GetReviewersChannels()

	assert.NoError(t, err)
	assert.Len(t, reviewersChannels, 2)
	assert.Equal(t, reviewersChannels[0].ID, "C1234")
	assert.Equal(t, reviewersChannels[1].ID, "C2345")
}
