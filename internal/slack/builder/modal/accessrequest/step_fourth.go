package accessrequest

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
)

type fourthStepBuilder struct {
	payload *accessrequest.StartTimeSelect
}

func NewFourthStepBuilder(p *accessrequest.StartTimeSelect) modal.Builder {
	return &fourthStepBuilder{
		payload: p,
	}
}

func (f *fourthStepBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (f *fourthStepBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, f.BuildStartDateBlock()...)
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildFourthStepSectionBlock())
	blockSet = append(blockSet, f.BuildAccessDurationBlock()...)
	return slack.Blocks{BlockSet: blockSet}
}

func (f *fourthStepBuilder) BuildStartDateBlock() []slack.Block {
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

func (f *fourthStepBuilder) BuildAccessDurationBlock() []slack.Block {
	text := BuildAccessDurationInfoText(f.payload.TTL)
	accessDurationInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	accessDurationOpts := f.BuildAccessDurationOpts()
	accessDurationOptsBlock := slack.NewActionBlock(
		util.ARequestAccessDurationOptionActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARequestAccessDurationOptionBlockActionID,
			accessDurationOpts...,
		),
	)
	return []slack.Block{accessDurationInfoBlock, accessDurationOptsBlock}
}

func (f *fourthStepBuilder) BuildAccessDurationOpts() []*slack.OptionBlockObject {
	var accessDurationOpts []*slack.OptionBlockObject
	accessDurationOpts = append(accessDurationOpts,
		slack.NewOptionBlockObject(
			util.ARequestAccessDurationFirstOption,
			slack.NewTextBlockObject(util.PlainText, util.ARequestAccessDurationFirstOption, false, false),
			nil,
		),
		slack.NewOptionBlockObject(
			util.ARequestAccessDurationSecondOption,
			slack.NewTextBlockObject(util.PlainText, util.ARequestAccessDurationSecondOption, false, false),
			nil,
		),
	)
	return accessDurationOpts
}

func (f *fourthStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.AccessDurationOptionSelectPrivateMetadataPayload{
		ChannelID:                   f.payload.RequesterChannelID,
		ChannelName:                 f.payload.RequesterChannelName,
		RealName:                    f.payload.RequesterRealName,
		RequireReason:               f.payload.RequireReason,
		SelectedRole:                f.payload.SelectedRole,
		SelectedChannelID:           f.payload.RequesterChannelID,
		SelectedChannelName:         f.payload.RequesterChannelName,
		SelectedStartDateOptionID:   f.payload.SelectedStartDateOptionID,
		SelectedStartDateOptionName: f.payload.SelectedStartDateOptionName,
		TTL:                         f.payload.TTL,
		SelectedStartDate:           f.payload.SelectedStartDate,
		SelectedStartTime:           f.payload.StartTime,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
