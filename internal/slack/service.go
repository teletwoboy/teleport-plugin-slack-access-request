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

package slack

import (
	"context"
	"fmt"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/message"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/types"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/slack-go/slack"
)

var tracer = otel.Tracer(telemetry.SlackService)

type Service interface {
	AddPinContext(ctx context.Context, channel, timestamp string) error
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) (*models.User, error)
	ExistsUserByID(ctx context.Context, id string) (bool, error)
	ExistsUserInChannelByID(ctx context.Context, id string, channelID string) (bool, error)
	FetchAllChannelsContext(ctx context.Context) ([]slack.Channel, error)
	FetchReviewersChannelByRole(ctx context.Context, role string) ([]types.ReviewersChannel, error)
	FetchTeamInfoContext(ctx context.Context) (*types.TeamInfo, error)
	FetchUserInfoContext(ctx context.Context, user string) (*models.User, error)
	FetchUsersContext(ctx context.Context) ([]models.User, error)
	FetchUsersInConversationContext(ctx context.Context, channelID string) ([]string, error)
	GetPermalinkContext(ctx context.Context, channel, timestamp string) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserBySlackUserID(ctx context.Context, id int32) (*models.User, error)
	OpenModalContext(ctx context.Context, triggerID string, builder modal.Builder) error
	PostEphemeralContext(ctx context.Context, channelID, userID string, builder message.Builder) (timestamp string, err error)
	PostMessageContext(ctx context.Context, channelID string, builder message.Builder) (string, string, error)
	PushModalContext(ctx context.Context, triggerID string, builder modal.Builder) error
	RemovePinContext(ctx context.Context, channel string, timestamp string) error
	UpdateMessageContext(ctx context.Context, channel, timestamp string, builder message.Builder) (string, string, string, error)
	UpdateModalContext(ctx context.Context, builder modal.Builder, externalID, hash, viewID string) error
}

type API interface {
	AddPinContext(ctx context.Context, channel string, item slack.ItemRef) error
	GetConversationsContext(ctx context.Context, params *slack.GetConversationsParameters) (channels []slack.Channel, nextCursor string, err error)
	GetPermalinkContext(ctx context.Context, params *slack.PermalinkParameters) (string, error)
	GetTeamInfoContext(ctx context.Context) (*slack.TeamInfo, error)
	GetUserInfoContext(ctx context.Context, user string) (*slack.User, error)
	GetUsersContext(ctx context.Context, options ...slack.GetUsersOption) ([]slack.User, error)
	GetUsersInConversationContext(ctx context.Context, params *slack.GetUsersInConversationParameters) ([]string, string, error)
	OpenViewContext(ctx context.Context, triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error)
	PostEphemeralContext(ctx context.Context, channelID, userID string, options ...slack.MsgOption) (timestamp string, err error)
	PostMessageContext(ctx context.Context, channel string, options ...slack.MsgOption) (string, string, error)
	PushViewContext(ctx context.Context, triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error)
	RemovePinContext(ctx context.Context, channel string, item slack.ItemRef) error
	UpdateMessageContext(ctx context.Context, channelID string, timestamp string, options ...slack.MsgOption) (string, string, string, error)
	UpdateViewContext(ctx context.Context, view slack.ModalViewRequest, externalID string, hash string, viewID string) (*slack.ViewResponse, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) (*models.User, error)
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

func (s *service) AddPinContext(ctx context.Context, channel, timestamp string) error {
	ctx, span := tracer.Start(ctx, "AddPinContext",
		trace.WithAttributes(
			attribute.String("channel", channel),
			attribute.String("timestamp", timestamp),
		),
	)
	defer span.End()

	itemRef := slack.ItemRef{
		Timestamp: timestamp,
	}
	return s.api.AddPinContext(ctx, channel, itemRef)
}

func (s *service) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.Start(ctx, "CreateUser",
		trace.WithAttributes(
			attribute.Int64("user.slackUserID", int64(user.SlackUserID)),
		),
	)
	defer span.End()

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create slack user: %w", err)
	}
	return createdUser, nil
}

func (s *service) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.Start(ctx, "DeleteUser",
		trace.WithAttributes(
			attribute.Int64("user.slackUserID", int64(user.SlackUserID)),
		),
	)
	defer span.End()

	DeletedUser, err := s.repo.DeleteUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create slack user: %w", err)
	}
	return DeletedUser, nil
}

func (s *service) ExistsUserByID(ctx context.Context, id string) (bool, error) {
	ctx, span := tracer.Start(ctx, "ExistsUserByID",
		trace.WithAttributes(
			attribute.String("id", id),
		),
	)
	defer span.End()

	exists, err := s.repo.ExistsUserByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (s *service) ExistsUserInChannelByID(ctx context.Context, id, channelID string) (bool, error) {
	ctx, span := tracer.Start(ctx, "ExistsUserInChannelByID",
		trace.WithAttributes(
			attribute.String("id", id),
			attribute.String("channel", channelID),
		),
	)
	defer span.End()

	ids, err := s.FetchUsersInConversationContext(ctx, channelID)
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

func (s *service) FetchAllChannelsContext(ctx context.Context) ([]slack.Channel, error) {
	ctx, span := tracer.Start(ctx, "FetchAllChannelsContext")
	defer span.End()

	var channels []slack.Channel
	params := &slack.GetConversationsParameters{
		ExcludeArchived: true,
		Types:           []string{"public_channel", "private_channel"},
	}

	for {
		rawChannels, nextCursor, err := s.api.GetConversationsContext(ctx, params)
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

func (s *service) FetchReviewersChannelByRole(ctx context.Context, role string) ([]types.ReviewersChannel, error) {
	ctx, span := tracer.Start(ctx, "FetchReviewersChannelByRole",
		trace.WithAttributes(
			attribute.String("role", role),
		),
	)
	defer span.End()

	channels, err := s.FetchAllChannelsContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all channels: %w", err)
	}

	reviewersChannels := filterReviewersChannel(channels, role)
	return convertToReviewersChannels(reviewersChannels), nil
}

func (s *service) FetchTeamInfoContext(ctx context.Context) (*types.TeamInfo, error) {
	ctx, span := tracer.Start(ctx, "FetchTeamInfoContext")
	defer span.End()

	rawTeamInfo, err := s.api.GetTeamInfoContext(ctx)
	if err != nil {
		return &types.TeamInfo{}, fmt.Errorf("failed to get team info from Slack API: %w", err)
	}
	return &types.TeamInfo{
		ID:   rawTeamInfo.ID,
		Name: rawTeamInfo.Name,
	}, nil
}

func (s *service) FetchUserInfoContext(ctx context.Context, user string) (*models.User, error) {
	ctx, span := tracer.Start(ctx, "FetchUserInfoContext",
		trace.WithAttributes(
			attribute.String("user", user),
		),
	)
	defer span.End()

	rawUser, err := s.api.GetUserInfoContext(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info from Slack API: %w", err)
	}
	return convertToUser(rawUser), nil
}

func (s *service) FetchUsersContext(ctx context.Context) ([]models.User, error) {
	ctx, span := tracer.Start(ctx, "FetchUsersContext")
	defer span.End()

	rawUsers, err := s.api.GetUsersContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users from Slack API: %w", err)
	}

	activeUsers := filterActiveUsers(rawUsers)
	return convertToUsers(activeUsers), nil
}

func (s *service) FetchUsersInConversationContext(ctx context.Context, channelID string) ([]string, error) {
	ctx, span := tracer.Start(ctx, "FetchUsersInConversationContext",
		trace.WithAttributes(
			attribute.String("channelID", channelID),
		),
	)
	defer span.End()

	var ids []string
	params := &slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	}

	for {
		rawUsers, nextCursor, err := s.api.GetUsersInConversationContext(ctx, params)
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

func (s *service) GetPermalinkContext(ctx context.Context, channel, timestamp string) (string, error) {
	ctx, span := tracer.Start(ctx, "GetPermalinkContext",
		trace.WithAttributes(
			attribute.String("channel", channel),
		),
	)
	defer span.End()

	params := &slack.PermalinkParameters{
		Channel: channel,
		Ts:      timestamp,
	}
	return s.api.GetPermalinkContext(ctx, params)
}

func (s *service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	ctx, span := tracer.Start(ctx, "GetUserByID",
		trace.WithAttributes(
			attribute.String("id", id),
		),
	)
	defer span.End()

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID (%s): %w", id, err)
	}
	return user, nil
}

func (s *service) GetUserBySlackUserID(ctx context.Context, id int32) (*models.User, error) {
	ctx, span := tracer.Start(ctx, "GetUserBySlackUserID",
		trace.WithAttributes(
			attribute.Int64("slackUserId", int64(id)),
		),
	)
	defer span.End()

	user, err := s.repo.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Slack user ID (%d): %w", id, err)
	}
	return user, nil
}

func (s *service) OpenModalContext(ctx context.Context, triggerID string, builder modal.Builder) error {
	ctx, span := tracer.Start(ctx, "OpenModalContext",
		trace.WithAttributes(
			attribute.String("triggerID", triggerID),
		),
	)
	defer span.End()

	builtModal, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build modal: %w", err)
	}

	_, err = s.api.OpenViewContext(ctx, triggerID, *builtModal)
	if err != nil {
		return fmt.Errorf("failed to open modal: %w", err)
	}
	return nil
}

func (s *service) PostEphemeralContext(ctx context.Context, channelID, userID string, builder message.Builder) (timestamp string, err error) {
	ctx, span := tracer.Start(ctx, "PostEphemeralContext",
		trace.WithAttributes(
			attribute.String("channelID", channelID),
			attribute.String("userID", userID),
		),
	)
	defer span.End()

	msgOption := builder.Build()
	return s.api.PostEphemeralContext(ctx, channelID, userID, msgOption)
}

func (s *service) PostMessageContext(ctx context.Context, channelID string, builder message.Builder) (string, string, error) {
	ctx, span := tracer.Start(ctx, "PostMessageContext",
		trace.WithAttributes(
			attribute.String("channelID", channelID),
		),
	)
	defer span.End()

	msgOption := builder.Build()
	return s.api.PostMessageContext(ctx, channelID, msgOption)
}

func (s *service) PushModalContext(ctx context.Context, triggerID string, builder modal.Builder) error {
	ctx, span := tracer.Start(ctx, "PushModalContext",
		trace.WithAttributes(
			attribute.String("triggerID", triggerID),
		),
	)
	defer span.End()

	builtModal, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build modal: %w", err)
	}

	_, err = s.api.PushViewContext(ctx, triggerID, *builtModal)
	if err != nil {
		return fmt.Errorf("failed to push modal: %w", err)
	}
	return nil
}

func (s *service) RemovePinContext(ctx context.Context, channel, timestamp string) error {
	ctx, span := tracer.Start(ctx, "RemovePinContext",
		trace.WithAttributes(
			attribute.String("channel", channel),
			attribute.String("timestamp", timestamp),
		),
	)
	defer span.End()

	itemRef := slack.ItemRef{
		Timestamp: timestamp,
	}
	return s.api.RemovePinContext(ctx, channel, itemRef)
}

func (s *service) UpdateMessageContext(ctx context.Context, channel, timestamp string, builder message.Builder) (string, string, string, error) {
	ctx, span := tracer.Start(ctx, "UpdateMessageContext",
		trace.WithAttributes(
			attribute.String("channel", channel),
			attribute.String("timestamp", timestamp),
		),
	)
	defer span.End()

	msgOption := builder.Build()
	return s.api.UpdateMessageContext(ctx, channel, timestamp, msgOption)
}

func (s *service) UpdateModalContext(ctx context.Context, builder modal.Builder, externalID, hash, viewID string) error {
	ctx, span := tracer.Start(ctx, "UpdateModalContext",
		trace.WithAttributes(
			attribute.String("externalID", externalID),
			attribute.String("viewID", viewID),
		),
	)
	defer span.End()

	builtModal, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build modal: %w", err)
	}

	_, err = s.api.UpdateViewContext(ctx, *builtModal, externalID, hash, viewID)
	if err != nil {
		return fmt.Errorf("failed to update modal: %w", err)
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

func convertToUser(user *slack.User) *models.User {
	return &models.User{
		ID:       user.ID,
		Name:     user.Name,
		RealName: user.RealName,
		Email:    user.Profile.Email,
		TimeZone: user.TZ,
	}
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
			TimeZone: copiedUser.TZ,
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
