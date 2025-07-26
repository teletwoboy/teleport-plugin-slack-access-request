package slack

import (
	"context"
	"fmt"
	"strings"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/types"

	"github.com/slack-go/slack"
)

type Service interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	ExistsUserByID(ctx context.Context, id string) (bool, error)
	ExistsUserInChannelByID(id string, channelID string) (bool, error)
	FetchAllChannels() ([]slack.Channel, error)
	FetchReviewersChannels() ([]types.ReviewersChannel, error)
	FetchTeamInfo() (*types.TeamInfo, error)
	FetchUsers() ([]models.User, error)
	FetchUsersInConversation(channelID string) ([]string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserBySlackUserID(ctx context.Context, id int32) (*models.User, error)
	OpenModal(triggerID string, builder modal.Builder) error
	PostMessage(channelID string, builder message.Builder) (string, string, error)
	PushModal(triggerID string, builder modal.Builder) error
}

type API interface {
	GetConversations(params *slack.GetConversationsParameters) (channels []slack.Channel, nextCursor string, err error)
	GetTeamInfo() (*slack.TeamInfo, error)
	GetUsers(options ...slack.GetUsersOption) ([]slack.User, error)
	GetUsersInConversation(params *slack.GetUsersInConversationParameters) ([]string, string, error)
	OpenView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error)
	PostMessage(channel string, options ...slack.MsgOption) (string, string, error)
	PushView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	ExistsUserByID(ctx context.Context, id string) (bool, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserBySlackUserID(ctx context.Context, id int32) (*models.User, error)
}

type service struct {
	api  API
	repo Repository
}

func NewService(api API, repo Repository) Service {
	return &service{api: api, repo: repo}
}

func (s *service) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create slack user: %w", err)
	}
	return createdUser, nil
}

func (s *service) ExistsUserByID(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.ExistsUserByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (s *service) ExistsUserInChannelByID(id, channelID string) (bool, error) {
	ids, err := s.FetchUsersInConversation(channelID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch users in channel from Slack API: %w", err)
	}
	for _, identification := range ids {
		if identification == id {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) FetchAllChannels() ([]slack.Channel, error) {
	var channels []slack.Channel
	params := &slack.GetConversationsParameters{
		ExcludeArchived: true,
		Types:           []string{"public_channel", "private_channel"},
	}

	for {
		rawChannels, nextCursor, err := s.api.GetConversations(params)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversations (cursor=%s): %w", params.Cursor, err)
		}
		channels = append(channels, rawChannels...)
		if nextCursor == "" {
			break
		}
	}
	return channels, nil
}

func (s *service) FetchReviewersChannels() ([]types.ReviewersChannel, error) {
	channels, err := s.FetchAllChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all channels: %w", err)
	}

	reviewersChannels := filterReviewersChannels(channels)
	joinedChannels := filterJoinedChannels(reviewersChannels)
	return convertToReviewersChannels(joinedChannels), nil
}

func (s *service) FetchReviewersChannelByRole(role string) ([]types.ReviewersChannel, error) {
	channels, err := s.FetchAllChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all channels: %w", err)
	}

	reviewersChannels := filterReviewersChannel(channels, role)
	joinedChannels := filterJoinedChannels(reviewersChannels)
	return convertToReviewersChannels(joinedChannels), nil
}

func (s *service) FetchTeamInfo() (*types.TeamInfo, error) {
	rawTeamInfo, err := s.api.GetTeamInfo()
	if err != nil {
		return &types.TeamInfo{}, fmt.Errorf("failed to get team info from Slack API: %w", err)
	}
	return &types.TeamInfo{
		ID:   rawTeamInfo.ID,
		Name: rawTeamInfo.Name,
	}, nil
}

func (s *service) FetchUsers() ([]models.User, error) {
	rawUsers, err := s.api.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users from Slack API: %w", err)
	}

	activeUsers := filterActiveUsers(rawUsers)
	return convertToUsers(activeUsers), nil
}

func (s *service) FetchUsersInConversation(channelID string) ([]string, error) {
	var ids []string
	params := &slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	}

	for {
		rawUsers, nextCursor, err := s.api.GetUsersInConversation(params)
		if err != nil {
			return nil, fmt.Errorf("failed to get users in conversation from Slack API: %w", err)
		}
		ids = append(ids, rawUsers...)
		if nextCursor == "" {
			break
		}
	}
	return ids, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID (%s): %w", id, err)
	}
	return user, nil
}

func (s *service) GetUserBySlackUserID(ctx context.Context, id int32) (*models.User, error) {
	user, err := s.repo.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Slack user ID (%d): %w", id, err)
	}
	return user, nil
}

func (s *service) OpenModal(triggerID string, builder modal.Builder) error {
	builtModal, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build modal: %w", err)
	}

	_, err = s.api.OpenView(triggerID, *builtModal)
	if err != nil {
		return fmt.Errorf("failed to open modal: %w", err)
	}
	return nil
}

func (s *service) PostMessage(channelID string, builder message.Builder) (string, string, error) {
	msgOption := builder.Build()
	return s.api.PostMessage(channelID, msgOption)
}

func (s *service) PushModal(triggerID string, builder modal.Builder) error {
	builtModal, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build modal: %w", err)
	}

	_, err = s.api.PushView(triggerID, *builtModal)
	if err != nil {
		return fmt.Errorf("failed to push modal: %w", err)
	}
	return nil
}

// --- Internal Util Functions related to User ---
func filterActiveUsers(users []slack.User) []slack.User {
	var activeUsers []slack.User
	for _, user := range users {
		copiedUser := user
		if !copiedUser.Deleted {
			activeUsers = append(activeUsers, copiedUser)
		}
	}
	return activeUsers
}

func convertToUsers(users []slack.User) []models.User {
	var result []models.User
	for _, user := range users {
		copiedUser := user
		result = append(result, models.User{
			ID:       copiedUser.ID,
			Name:     copiedUser.Name,
			RealName: copiedUser.RealName,
			Email:    copiedUser.Profile.Email,
		})
	}
	return result
}

// --- Internal Util Functions related to ReviewersChannel ---
func filterReviewersChannel(channels []slack.Channel, role string) []slack.Channel {
	var reviewersChannels []slack.Channel
	for _, channel := range channels {
		copiedChannel := channel
		if copiedChannel.Name == role+"-reviewers" {
			reviewersChannels = append(reviewersChannels, copiedChannel)
		}
	}
	return reviewersChannels
}

func filterReviewersChannels(channels []slack.Channel) []slack.Channel {
	var reviewersChannels []slack.Channel
	for _, channel := range channels {
		copiedChannel := channel
		if strings.HasSuffix(copiedChannel.Name, "-reviewers") {
			reviewersChannels = append(reviewersChannels, copiedChannel)
		}
	}
	return reviewersChannels
}

func filterJoinedChannels(channels []slack.Channel) []slack.Channel {
	var joinedChannels []slack.Channel
	for _, channel := range channels {
		copiedChannel := channel
		if copiedChannel.IsMember {
			joinedChannels = append(joinedChannels, copiedChannel)
		}
	}
	return joinedChannels
}

func convertToReviewersChannels(channels []slack.Channel) []types.ReviewersChannel {
	var result []types.ReviewersChannel
	for _, channel := range channels {
		copiedChannel := channel
		result = append(result, types.ReviewersChannel{
			ID:       copiedChannel.ID,
			Name:     copiedChannel.Name,
			IsMember: copiedChannel.IsMember,
		})
	}
	return result
}
