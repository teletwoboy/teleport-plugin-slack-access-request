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
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type secondStepBuilder struct {
	payload *accesspolicy.ChannelSelect
	roles   map[string]struct{}
}

func NewSecondStepBuilder(p *accesspolicy.ChannelSelect, r map[string]struct{}) modal.Builder {
	return &secondStepBuilder{
		payload: p,
		roles:   r,
	}
}

func (s *secondStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := s.BuildBlocks()
	privateMetadata, err := s.BuildPrivateMetadata()
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

func (s *secondStepBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildFirstStepSectionBlock())
	blockSet = append(blockSet, s.BuildChannelBlock())
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildSecondStepSectionBlock())
	blockSet = append(blockSet, s.BuildRoleBlock())
	return slack.Blocks{BlockSet: blockSet}
}

func (s *secondStepBuilder) BuildChannelBlock() *slack.SectionBlock {
	text := "```\n" + s.payload.ChannelName + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (s *secondStepBuilder) BuildRoleBlock() *slack.ActionBlock {
	roleOpts := s.BuildRoleOpts()
	return slack.NewActionBlock(
		util.APolicyRoleActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.APolicyRoleOptionBlockActionID,
			roleOpts...,
		),
	)
}

func (s *secondStepBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	roleOpts = append(roleOpts, slack.NewOptionBlockObject(
		util.APolicyAllOptionValue,
		slack.NewTextBlockObject(util.PlainText, util.APolicyAllOption, false, false),
		nil,
	))
	for r := range s.roles {
		copiedRole := r
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			copiedRole,
			slack.NewTextBlockObject(util.PlainText, copiedRole, false, false),
			nil,
		))
	}
	return roleOpts
}

func (s *secondStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.RoleSelectPrivateMetadataPayload{
		ChannelID:           s.payload.RequesterChannelID,
		ChannelName:         s.payload.RequesterChannelName,
		RealName:            s.payload.RequesterRealName,
		TimeZone:            s.payload.RequesterTimeZone,
		SelectedChannelID:   s.payload.ChannelID,
		SelectedChannelName: s.payload.ChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
