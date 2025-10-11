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

package accessrequest

import (
	"encoding/json"
	"fmt"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/types"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type firstStepBuilder struct {
	accessInfo *types.UserAccessInfo
	payload    *slashcommands.AccessRole
	slackUser  *models.User
}

func NewFirstStepBuilder(a *types.UserAccessInfo, p *slashcommands.AccessRole, s *models.User) modal.Builder {
	return &firstStepBuilder{
		accessInfo: a,
		payload:    p,
		slackUser:  s,
	}
}

func (f *firstStepBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(f.accessInfo.Roles) == 0 {
		return nil, fmt.Errorf("<%s> does not have any Role to request. Please contact the administrator", f.slackUser.RealName)
	}
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.ARequestTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Close, false, false),
		Submit:          nil,
		CallbackID:      util.ARequestCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *firstStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock(f.slackUser.RealName)
	roleBlock := f.BuildRoleBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			roleBlock,
		},
	}
	return blocks
}

func (f *firstStepBuilder) BuildRoleBlock() *slack.ActionBlock {
	roleOpts := f.BuildRoleOpts()
	return slack.NewActionBlock(
		util.ARequestRoleActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARequestRoleOptionBlockActionID,
			roleOpts...,
		),
	)
}

func (f *firstStepBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	for _, role := range f.accessInfo.Roles {
		r := role
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			r,
			slack.NewTextBlockObject("plain_text", role, false, false),
			nil,
		))
	}
	return roleOpts
}

func (f *firstStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.RoleSelectPrivateMetadataPayload{
		ChannelID:     f.payload.ChannelID,
		ChannelName:   f.payload.ChannelName,
		RealName:      f.slackUser.RealName,
		RequireReason: f.accessInfo.RequireReason,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
