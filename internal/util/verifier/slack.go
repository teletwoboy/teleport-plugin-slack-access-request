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

package verifier

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"
)

type Slack struct {
	Srv slack.Service
}

func NewSlack(srv slack.Service) *Slack {
	return &Slack{
		Srv: srv,
	}
}

func (s *Slack) VerifyChanIsReviewersChan(channelName string) error {
	if strings.HasSuffix(channelName, "-reviewers") {
		return nil
	}
	return fmt.Errorf("channel %s is not reviewers channel", channelName)
}

func (s *Slack) VerifyUserExistsBySlackUserID(ctx context.Context, id int32) (*models.User, error) {
	slackUser, err := s.Srv.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	if slackUser == nil {
		return nil, fmt.Errorf("user not found in database")
	}
	return slackUser, nil
}

func (s *Slack) VerifyUserExistsByUsernameFromClient(ctx context.Context, username string) (*models.User, error) {
	slackUsers, err := s.Srv.FetchUsersContext(ctx)
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
	return nil, fmt.Errorf("user <%s> not found in Slack", username)
}

func (s *Slack) VerifyUserExistsByID(ctx context.Context, id, name string) error {
	exists, err := s.Srv.ExistsUserByID(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user <%s> not found in database", name)
	}
	return nil
}

func (s *Slack) VerifyUserExistsInChannelByID(ctx context.Context, id, channelID string) error {
	exists, err := s.Srv.ExistsUserInChannelByID(ctx, id, channelID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user not found in channel")
	}
	return nil
}

func (s *Slack) VerifyMessageAlreadyPosted() {
}
