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

package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type thirdStepBuilder struct {
	payload    *accesspolicy.RoleSelect
	slackUsers []models.User
}

func NewThirdStepBuilder(p *accesspolicy.RoleSelect, s []models.User) modal.Builder {
	return &thirdStepBuilder{
		payload:    p,
		slackUsers: s,
	}
}

func (t *thirdStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := t.BuildBlocks()
	privateMetadata, err := t.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.APolicyTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Close, false, false),
		Submit:          nil,
		CallbackID:      util.APolicyCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (t *thirdStepBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildFirstStepSectionBlock())
	blockSet = append(blockSet, t.BuildChannelBlock())
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildSecondStepSectionBlock())
	blockSet = append(blockSet, t.BuildRoleBlock())
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, t.BuildUserBlock())
	return slack.Blocks{BlockSet: blockSet}
}

func (t *thirdStepBuilder) BuildChannelBlock() *slack.SectionBlock {
	text := "```\n" + t.payload.SelectedChannelName + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (t *thirdStepBuilder) BuildRoleBlock() *slack.SectionBlock {
	text := "```\n" + t.payload.RoleName + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (t *thirdStepBuilder) BuildUserBlock() *slack.ActionBlock {
	userOpts := t.BuildUserOpts()
	return slack.NewActionBlock(
		util.APolicyUserActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.APolicyUserOptionBlockActionID,
			userOpts...,
		),
	)
}

func (t *thirdStepBuilder) BuildUserOpts() []*slack.OptionBlockObject {
	var userOpts []*slack.OptionBlockObject
	userOpts = append(userOpts, slack.NewOptionBlockObject(
		util.APolicyAllOptionValue,
		slack.NewTextBlockObject(util.PlainText, util.APolicyAllOption, false, false),
		nil,
	))
	for _, u := range t.slackUsers {
		copiedUser := u
		userOpts = append(userOpts, slack.NewOptionBlockObject(
			copiedUser.ID,
			slack.NewTextBlockObject(util.PlainText, copiedUser.RealName, false, false),
			nil,
		))
	}
	return userOpts
}

func (t *thirdStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.UserSelectPrivateMetadataPayload{
		ChannelID:           t.payload.RequesterChannelID,
		ChannelName:         t.payload.RequesterChannelName,
		RealName:            t.payload.RequesterRealName,
		TimeZone:            t.payload.RequesterTimeZone,
		SelectedChannelID:   t.payload.SelectedChannelID,
		SelectedChannelName: t.payload.SelectedChannelName,
		SelectedRole:        t.payload.Role,
		SelectedRoleName:    t.payload.RoleName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
