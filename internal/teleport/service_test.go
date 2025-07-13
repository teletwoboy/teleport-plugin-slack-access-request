package teleport

import (
	"context"
	"teleport-plugin-slack-access-request/internal/teleport/mocks"
	"testing"

	"github.com/gravitational/teleport/api/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetUsersWithoutSecrets(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPI := mocks.NewMockAPI(ctrl)
	mockUser1 := mocks.NewMockUser(ctrl)
	mockUser1.EXPECT().GetName().Return("user1").AnyTimes()
	mockUser1.EXPECT().IsBot().Return(false).AnyTimes()
	mockUser2 := mocks.NewMockUser(ctrl)
	mockUser2.EXPECT().GetName().Return("bot123").AnyTimes()
	mockUser2.EXPECT().IsBot().Return(true).AnyTimes()
	mockUser3 := mocks.NewMockUser(ctrl)
	mockUser3.EXPECT().GetName().Return("user2").AnyTimes()
	mockUser3.EXPECT().IsBot().Return(false).AnyTimes()

	mockAPI.EXPECT().
		GetUsers(gomock.Any(), false).
		Return([]types.User{mockUser1, mockUser2, mockUser3}, nil)
	service := &Service{api: mockAPI}

	users, err := service.GetUsersWithoutSecrets(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)
}
