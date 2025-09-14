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

type fifthStepBuilder struct {
	payload *accesspolicy.EndTimeSelect
}

func NewFifthStepBuilder(p *accesspolicy.EndTimeSelect) modal.Builder {
	return &fifthStepBuilder{
		payload: p,
	}
}

func (f *fifthStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.APolicyTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          nil,
		CallbackID:      util.APolicyCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *fifthStepBuilder) BuildBlocks() slack.Blocks {
	fourthStep := BuildFourthStepSectionBlock()
	fourthCautionStep := BuildFourthStepCautionSectionBlock()
	fourthStepFirstSub := BuildFourthStepFirstSubSectionBlock()
	startDateTimeBlock := f.BuildStartDateTimeBlock()
	fourthStepSecondSub := BuildFourthStepSecondSubSectionBlock()
	endDateTimeBlock := f.BuildEndDateTimeBlock()
	fifthStep := BuildFifthStepSectionBlock()
	effectBlock := f.BuildEffectBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			fourthStep,
			fourthCautionStep,
			fourthStepFirstSub,
			startDateTimeBlock,
			fourthStepSecondSub,
			endDateTimeBlock,
			slack.NewDividerBlock(),
			fifthStep,
			effectBlock,
		},
	}
	return blocks
}

func (f *fifthStepBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APolicyStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APolicyStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APolicyStartTimeBlockActionID),
	)
}

func (f *fifthStepBuilder) BuildEndDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APolicyEndDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APolicyEndDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APolicyEndTimeBlockActionID),
	)
}

func (f *fifthStepBuilder) BuildEffectBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APolicyEffectBlockID,
		slack.NewButtonBlockElement(
			util.APolicyAllowButtonBlockActionID,
			util.APolicyAllowButtonValue,
			slack.NewTextBlockObject(util.PlainText, util.APolicyAllowButtonText, false, false),
		),
		slack.NewButtonBlockElement(
			util.APolicyDenyButtonBlockActionID,
			util.APolicyDenyButtonValue,
			slack.NewTextBlockObject(util.PlainText, util.APolicyDenyButtonText, false, false),
		),
	)
}

func (f *fifthStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.EffectSelectPrivateMetadataPayload{
		ChannelID:           f.payload.RequesterChannelID,
		ChannelName:         f.payload.RequesterChannelName,
		RealName:            f.payload.RequesterRealName,
		TimeZone:            f.payload.RequesterTimeZone,
		SelectedChannelID:   f.payload.SelectedChannelID,
		SelectedChannelName: f.payload.SelectedChannelName,
		SelectedRole:        f.payload.SelectedRole,
		SelectedRoleName:    f.payload.SelectedRoleName,
		SelectedUserID:      f.payload.SelectedUserID,
		SelectedRealName:    f.payload.SelectedRealName,
		SelectedStartDate:   f.payload.SelectedStartDate,
		SelectedStartTime:   f.payload.SelectedStartTime,
		SelectedEndDate:     f.payload.SelectedEndDate,
		SelectedEndTime:     f.payload.EndTime,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
