package verifier

import (
	"context"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/util"
)

type Slack struct {
	Srv slack.Service
}

func NewSlack(srv slack.Service) *Slack {
	return &Slack{
		Srv: srv,
	}
}

func (s *Slack) VerifyUserExistsBySlackID(ctx context.Context, id int32) (*models.User, error) {
	slackUser, err := s.Srv.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	if slackUser == nil {
		return nil, fmt.Errorf("user id %d not found in database", id)
	}
	return slackUser, nil
}

func (s *Slack) VerifyUserExistsByUsernameFromClient(username string) (*models.User, error) {
	slackUsers, err := s.Srv.FetchUsers()
	if err != nil {
		slog.Error("failed to fetch users", "err", err)
		return nil, err
	}

	for _, s := range slackUsers {
		copiedUser := s
		email := copiedUser.Email
		if util.MatchesIdentifier(username, email) {
			return &copiedUser, nil
		}
	}
	return nil, fmt.Errorf("username %s not found in Slack", username)
}

func (s *Slack) VerifyUserExistsByID(ctx context.Context, id, name string) error {
	exists, err := s.Srv.ExistsUserByID(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user %s not found in database", name)
	}
	return nil
}

func (s *Slack) VerifyUserInChannelExistsByID(id, channelID string) error {
	exists, err := s.Srv.ExistsUserInChannelByID(id, channelID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user %s not found in channel %s", id, channelID)
	}
	return nil
}
