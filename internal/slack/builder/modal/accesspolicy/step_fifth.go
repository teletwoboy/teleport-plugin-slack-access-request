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
		Title:           slack.NewTextBlockObject(util.PlainText, util.APTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          nil,
		CallbackID:      util.APCallBackID,
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
		util.APStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APStartTimeBlockActionID),
	)
}

func (f *fifthStepBuilder) BuildEndDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APEndDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APEndDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APEndTimeBlockActionID),
	)
}

func (f *fifthStepBuilder) BuildEffectBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APEffectBlockID,
		slack.NewButtonBlockElement(
			util.APAllowButtonBlockActionID,
			util.APAllowButtonValue,
			slack.NewTextBlockObject(util.PlainText, util.APAllowButtonText, false, false),
		),
		slack.NewButtonBlockElement(
			util.APDenyButtonBlockActionID,
			util.APDenyButtonValue,
			slack.NewTextBlockObject(util.PlainText, util.APDenyButtonText, false, false),
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
