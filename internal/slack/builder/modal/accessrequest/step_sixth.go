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
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type sixthStepBuilder struct {
	payload *accessrequest.RequestTTLTimeSelect
}

func NewSixthStepBuilder(p *accessrequest.RequestTTLTimeSelect) modal.Builder {
	return &sixthStepBuilder{
		payload: p,
	}
}

func (s *sixthStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := s.BuildBlocks()
	privateMetadata, err := s.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.ARequestTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          slack.NewTextBlockObject(util.PlainText, util.Submit, false, false),
		CallbackID:      util.ARequestCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (s *sixthStepBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildSixthStepSectionBlock())
	blockSet = append(blockSet, s.BuildSummaryBlock())
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, s.RequestReasonBlock())
	return slack.Blocks{BlockSet: blockSet}
}

func (s *sixthStepBuilder) BuildSummaryBlock() slack.Block {
	text := BuildSummaryInfoText(s.payload)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (s *sixthStepBuilder) RequestReasonBlock() slack.Block {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(util.PlainText, util.ARequestReasonElemBlockTest, false, false),
		util.ARequestReasonElemBlockActionID,
	)
	reasonBlock := slack.NewInputBlock(
		util.ARequestReasonBlockID,
		slack.NewTextBlockObject(util.PlainText, util.ARequestReasonBlockText, false, false),
		nil,
		reasonElement,
	)

	if !s.payload.RequireReason {
		reasonBlock.Optional = true
	}
	return reasonBlock
}

func (s *sixthStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.SummaryPrivateMetadataPayload{
		ChannelID:                        s.payload.RequesterChannelID,
		ChannelName:                      s.payload.RequesterChannelName,
		RealName:                         s.payload.RequesterRealName,
		RequireReason:                    s.payload.RequireReason,
		SelectedRole:                     s.payload.SelectedRole,
		SelectedChannelID:                s.payload.SelectedChannelID,
		SelectedChannelName:              s.payload.SelectedChannelName,
		SelectedStartDateOptionID:        s.payload.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      s.payload.SelectedStartDateOptionName,
		TTL:                              s.payload.TTL,
		SelectedStartDate:                s.payload.SelectedStartDate,
		SelectedStartTime:                s.payload.SelectedStartTime,
		SelectedAccessDurationOptionID:   s.payload.SelectedAccessDurationOptionID,
		SelectedAccessDurationOptionName: s.payload.SelectedAccessDurationOptionName,
		SelectedAccessDurationDate:       s.payload.SelectedAccessDurationDate,
		SelectedAccessDurationTime:       s.payload.SelectedAccessDurationTime,
		RequestTTL:                       s.payload.RequestTTL,
		SelectedRequestTTLOptionID:       s.payload.SelectedRequestTTLOptionID,
		SelectedRequestTTLOptionName:     s.payload.SelectedRequestTTLOptionName,
		SelectedRequestTTLDate:           s.payload.SelectedRequestTTLDate,
		SelectedRequestTTLTime:           s.payload.RequestTTLTime,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
