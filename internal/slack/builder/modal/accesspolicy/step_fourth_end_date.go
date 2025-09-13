package accesspolicy

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"
)

type fourthStepEndDateBuilder struct {
	payload *accesspolicy.StartTimeSelect
}

func NewFourthStepEndDateBuilder(p *accesspolicy.StartTimeSelect) modal.Builder {
	return &fourthStepEndDateBuilder{
		payload: p,
	}
}

func (f *fourthStepEndDateBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (f *fourthStepEndDateBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildFourthStepSectionBlock())
	blockSet = append(blockSet, f.BuildDurationBlock()...)
	return slack.Blocks{BlockSet: blockSet}
}

func (f *fourthStepEndDateBuilder) BuildDurationBlock() []slack.Block {
	fourthCautionStep := BuildFourthStepCautionSectionBlock()
	fourthStepFirstSub := BuildFourthStepFirstSubSectionBlock()
	startDateTimeBlock := f.BuildStartDateTimeBlock()
	fourthStepSecondSub := BuildFourthStepSecondSubSectionBlock()
	endDateTimeBlock := f.BuildEndDateTimeBlock()
	return []slack.Block{fourthCautionStep, fourthStepFirstSub, startDateTimeBlock, fourthStepSecondSub, endDateTimeBlock}
}

func (f *fourthStepEndDateBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APolicyStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APolicyStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APolicyStartTimeBlockActionID),
	)
}

func (f *fourthStepEndDateBuilder) BuildEndDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APolicyEndDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APolicyEndDateBlockActionID),
	)
}

func (f *fourthStepEndDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.EndDateSelectPrivateMetadataPayload{
		ChannelID:           f.payload.RequesterChannelID,
		ChannelName:         f.payload.RequesterChannelName,
		RealName:            f.payload.RequesterRealName,
		TimeZone:            f.payload.RequesterTimezone,
		SelectedChannelID:   f.payload.SelectedChannelID,
		SelectedChannelName: f.payload.SelectedChannelName,
		SelectedRole:        f.payload.SelectedRole,
		SelectedRoleName:    f.payload.SelectedRoleName,
		SelectedUserID:      f.payload.SelectedUserID,
		SelectedRealName:    f.payload.SelectedRealName,
		SelectedStartDate:   f.payload.SelectedStartDate,
		SelectedStartTime:   f.payload.StartTime,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
