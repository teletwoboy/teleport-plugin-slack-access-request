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
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	slacktypes "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/types"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type secondStepBuilder struct {
	channels []slacktypes.ReviewersChannel
	payload  *accessrequest.RoleSelect
}

func NewSecondStepBuilder(c []slacktypes.ReviewersChannel, p *accessrequest.RoleSelect) modal.Builder {
	return &secondStepBuilder{
		channels: c,
		payload:  p,
	}
}

func (s *secondStepBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(s.channels) == 0 {
		return nil, fmt.Errorf("role <%s> does not have any Reviewers Channel. Please contact the administrator", s.payload.Role)
	}
	blocks := s.BuildBlocks()
	privateMetadata, err := s.BuildPrivateMetadata()
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

func (s *secondStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock(s.payload.RequesterRealName)
	roleBlock := s.BuildRoleBlock()
	secondStep := BuildSecondStepSectionBlock()
	channelBlock := s.BuildChannelBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			roleBlock,
			slack.NewDividerBlock(),
			secondStep,
			channelBlock,
		},
	}
	return blocks
}

func (s *secondStepBuilder) BuildRoleBlock() *slack.SectionBlock {
	text := "```\n" + s.payload.Role + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (s *secondStepBuilder) BuildChannelBlock() *slack.ActionBlock {
	channelOpts := s.BuildChannelOpts()
	return slack.NewActionBlock(
		util.ARequestChannelActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARequestChannelOptionBlockActionID,
			channelOpts...,
		),
	)
}

func (s *secondStepBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOptions []*slack.OptionBlockObject
	for _, ch := range s.channels {
		id := ch.ID
		label := ch.Name
		channelOptions = append(channelOptions, slack.NewOptionBlockObject(
			id,
			slack.NewTextBlockObject(util.PlainText, label, false, false),
			nil,
		))
	}
	return channelOptions
}

func (s *secondStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.ChannelSelectPrivateMetadataPayload{
		ChannelID:     s.payload.RequesterChannelID,
		ChannelName:   s.payload.RequesterChannelName,
		RealName:      s.payload.RequesterRealName,
		RequireReason: s.payload.RequireReason,
		SelectedRole:  s.payload.Role,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
