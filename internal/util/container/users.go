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

package container

import (
	"context"
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	usermodels "teleport-plugin-slack-access-request/internal/user/models"
)

type Users struct {
	Slack    *slackmodels.User
	Teleport *teleportmodels.User
	User     *usermodels.User
}

func NewUsers(ctx context.Context, services *Services, slackUserID string) (*Users, error) {
	slackUser, err := services.Slack.GetUserByID(ctx, slackUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get slack user by ID: %w", err)
	}

	//    2. User
	user, err := services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Slack user ID: %w", err)
	}

	//    3. Teleport User
	teleportUser, err := services.Teleport.GetUserByTeleportUserID(ctx, user.TeleportUser.TeleportUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by Teleport user ID: %w", err)
	}
	return &Users{
		Slack:    slackUser,
		Teleport: teleportUser,
		User:     user,
	}, nil
}
