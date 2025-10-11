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

type fourthStepDateBuilder struct {
	payload *accessrequest.AccessDurationOptionSelect
}

func NewFourthStepDateBuilder(p *accessrequest.AccessDurationOptionSelect) modal.Builder {
	return &fourthStepDateBuilder{
		payload: p,
	}
}

func (f *fourthStepDateBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.ARequestTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          nil,
		CallbackID:      util.ARequestCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *fourthStepDateBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, f.BuildStartDateBlock()...)
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildFourthStepSectionBlock())
	blockSet = append(blockSet, f.BuildAccessDurationBlock()...)
	return slack.Blocks{BlockSet: blockSet}
}

func (f *fourthStepDateBuilder) BuildStartDateBlock() []slack.Block {
	if f.payload.SelectedStartDateOptionID == util.ARequestStartDateFirstOption {
		text := BuildStartDateFirstOptInfoText(f.payload.SelectedStartDateOptionName)
		startDateFirstOptInfoBlock := slack.NewSectionBlock(
			slack.NewTextBlockObject(util.Markdown, text, false, false),
			nil,
			nil,
		)
		return []slack.Block{startDateFirstOptInfoBlock}
	}
	text := BuildStartDateSecondOptInfoText(f.payload.SelectedStartDateOptionName, f.payload.TTL)
	startDateSecondOptInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	startDateTimeBlock := slack.NewActionBlock(
		util.ARequestStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.ARequestStartTimeBlockActionID),
	)
	return []slack.Block{startDateSecondOptInfoBlock, startDateTimeBlock}
}

func (f *fourthStepDateBuilder) BuildAccessDurationBlock() []slack.Block {
	text := BuildAccessDurationSecondOptInfoText(f.payload.AccessDurationOptionName, f.payload.TTL)
	accessDurationOptInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	accessDurationDateTimeBlock := slack.NewActionBlock(
		util.ARequestAccessDurationDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestAccessDurationDateBlockActionID),
	)
	return []slack.Block{accessDurationOptInfoBlock, accessDurationDateTimeBlock}
}

func (f *fourthStepDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.AccessDurationDateSelectPrivateMetadataPayload{
		ChannelID:                        f.payload.RequesterChannelID,
		ChannelName:                      f.payload.RequesterChannelName,
		RealName:                         f.payload.RequesterRealName,
		RequireReason:                    f.payload.RequireReason,
		SelectedRole:                     f.payload.SelectedRole,
		SelectedChannelID:                f.payload.SelectedChannelID,
		SelectedChannelName:              f.payload.SelectedChannelName,
		SelectedStartDateOptionID:        f.payload.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      f.payload.SelectedStartDateOptionName,
		TTL:                              f.payload.TTL,
		SelectedStartDate:                f.payload.SelectedStartDate,
		SelectedStartTime:                f.payload.SelectedStartTime,
		SelectedAccessDurationOptionID:   f.payload.AccessDurationOptionID,
		SelectedAccessDurationOptionName: f.payload.AccessDurationOptionName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
