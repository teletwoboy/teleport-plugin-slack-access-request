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
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type firstStepBuilder struct {
	channels  []slack.Channel
	payload   *slashcommands.AccessPolicy
	slackUser *models.User
}

func NewFirstStepBuilder(c []slack.Channel, p *slashcommands.AccessPolicy, s *models.User) modal.Builder {
	return &firstStepBuilder{
		channels:  c,
		payload:   p,
		slackUser: s,
	}
}

func (f *firstStepBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(f.channels) == 0 {
		return nil, fmt.Errorf("no channels found. Please contact the administrator")
	}
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.APolicyTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, "Close", false, false),
		Submit:          nil,
		CallbackID:      util.APolicyCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *firstStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock()
	channelBlock := f.BuildChannelBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			channelBlock,
		},
	}
	return blocks
}

func (f *firstStepBuilder) BuildChannelBlock() *slack.ActionBlock {
	channelOpts := f.BuildChannelOpts()
	return slack.NewActionBlock(
		util.APolicyChanActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.APolicyChanOptionBlockActionID,
			channelOpts...,
		),
	)
}

func (f *firstStepBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOpts []*slack.OptionBlockObject
	channelOpts = append(channelOpts, slack.NewOptionBlockObject(
		util.APolicyAllOptionValue,
		slack.NewTextBlockObject(util.PlainText, util.APolicyAllOption, false, false),
		nil,
	))
	for _, channel := range f.channels {
		copiedChannel := channel
		channelOpts = append(channelOpts, slack.NewOptionBlockObject(
			copiedChannel.ID,
			slack.NewTextBlockObject(util.PlainText, copiedChannel.Name, false, false),
			nil,
		))
	}
	return channelOpts
}

func (f *firstStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.ChannelSelectPrivateMetadataPayload{
		ChannelID:   f.payload.ChannelID,
		ChannelName: f.payload.ChannelName,
		RealName:    f.slackUser.RealName,
		TimeZone:    f.slackUser.TimeZone,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
