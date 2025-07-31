package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type fourthStepStartDateBuilder struct {
	payload *accesspolicy.UserSelect
}

func NewFourthStepStartDateBuilder(p *accesspolicy.UserSelect) modal.Builder {
	return &fourthStepStartDateBuilder{
		payload: p,
	}
}

func (f *fourthStepStartDateBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (f *fourthStepStartDateBuilder) BuildBlocks() slack.Blocks {
	fourthStep := BuildFourthStepSectionBlock()
	fourthCautionStep := BuildFourthStepCautionSectionBlock()
	fourthStepFirstSub := BuildFourthStepFirstSubSectionBlock()
	startDateTimeBlock := f.BuildStartDateTimeBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			fourthStep,
			fourthCautionStep,
			fourthStepFirstSub,
			startDateTimeBlock,
		},
	}
	return blocks
}

func (f *fourthStepStartDateBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APStartDateBlockActionID),
	)
}

func (f *fourthStepStartDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.StartDateSelectPrivateMetadataPayload{
		ChannelID:           f.payload.RequesterChannelID,
		ChannelName:         f.payload.RequesterChannelName,
		RealName:            f.payload.RequesterRealName,
		TimeZone:            f.payload.RequesterTimeZone,
		SelectedChannelID:   f.payload.SelectedChannelID,
		SelectedChannelName: f.payload.SelectedChannelName,
		SelectedRole:        f.payload.SelectedRole,
		SelectedRoleName:    f.payload.SelectedRoleName,
		SelectedUserID:      f.payload.UserID,
		SelectedRealName:    f.payload.RealName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -------------------------------------------------------------------------

type fourthStepStartTimeBuilder struct {
	payload *accesspolicy.StartDateSelect
}

func NewFourthStepStartTimeBuilder(p *accesspolicy.StartDateSelect) modal.Builder {
	return &fourthStepStartTimeBuilder{
		payload: p,
	}
}

func (f *fourthStepStartTimeBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (f *fourthStepStartTimeBuilder) BuildBlocks() slack.Blocks {
	fourthStep := BuildFourthStepSectionBlock()
	fourthCautionStep := BuildFourthStepCautionSectionBlock()
	fourthStepFirstSub := BuildFourthStepFirstSubSectionBlock()
	startDateTimeBlock := f.BuildStartDateTimeBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			fourthStep,
			fourthCautionStep,
			fourthStepFirstSub,
			startDateTimeBlock,
		},
	}
	return blocks
}

func (f *fourthStepStartTimeBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APStartTimeBlockActionID),
	)
}

func (f *fourthStepStartTimeBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.StartTimeSelectPrivateMetadataPayload{
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
		SelectedStartDate:   f.payload.StartDate,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -------------------------------------------------------------------------

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
		Title:           slack.NewTextBlockObject(util.PlainText, util.APTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          nil,
		CallbackID:      util.APCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *fourthStepEndDateBuilder) BuildBlocks() slack.Blocks {
	fourthStep := BuildFourthStepSectionBlock()
	fourthCautionStep := BuildFourthStepCautionSectionBlock()
	fourthStepFirstSub := BuildFourthStepFirstSubSectionBlock()
	startDateTimeBlock := f.BuildStartDateTimeBlock()
	fourthStepSecondSub := BuildFourthStepSecondSubSectionBlock()
	endDateTimeBlock := f.BuildEndDateTimeBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			fourthStep,
			fourthCautionStep,
			fourthStepFirstSub,
			startDateTimeBlock,
			fourthStepSecondSub,
			endDateTimeBlock,
		},
	}
	return blocks
}

func (f *fourthStepEndDateBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APStartTimeBlockActionID),
	)
}

func (f *fourthStepEndDateBuilder) BuildEndDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APEndDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APEndDateBlockActionID),
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

// -------------------------------------------------------------------------

type fourthStepEndTimeBuilder struct {
	payload *accesspolicy.EndDateSelect
}

func NewFourthStepEndTimeBuilder(p *accesspolicy.EndDateSelect) modal.Builder {
	return &fourthStepEndTimeBuilder{
		payload: p,
	}
}

func (f *fourthStepEndTimeBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (f *fourthStepEndTimeBuilder) BuildBlocks() slack.Blocks {
	fourthStep := BuildFourthStepSectionBlock()
	fourthCautionStep := BuildFourthStepCautionSectionBlock()
	fourthStepFirstSub := BuildFourthStepFirstSubSectionBlock()
	startDateTimeBlock := f.BuildStartDateTimeBlock()
	fourthStepSecondSub := BuildFourthStepSecondSubSectionBlock()
	endDateTimeBlock := f.BuildEndDateTimeBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			fourthStep,
			fourthCautionStep,
			fourthStepFirstSub,
			startDateTimeBlock,
			fourthStepSecondSub,
			endDateTimeBlock,
		},
	}
	return blocks
}

func (f *fourthStepEndTimeBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APStartDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APStartTimeBlockActionID),
	)
}

func (f *fourthStepEndTimeBuilder) BuildEndDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.APEndDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.APEndDateBlockActionID),
		slack.NewTimePickerBlockElement(util.APEndTimeBlockActionID),
	)
}

func (f *fourthStepEndTimeBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.EndTimeSelectPrivateMetadataPayload{
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
		SelectedEndDate:     f.payload.EndDate,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
