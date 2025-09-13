package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type fifthStepDateBuilder struct {
	payload *accessrequest.RequestTTLOptionSelect
}

func NewFifthStepDateBuilder(p *accessrequest.RequestTTLOptionSelect) modal.Builder {
	return &fifthStepDateBuilder{
		payload: p,
	}
}

func (f *fifthStepDateBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (f *fifthStepDateBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, f.BuildStartDateBlock()...)
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildFourthStepSectionBlock())
	blockSet = append(blockSet, f.BuildAccessDurationBlock()...)
	blockSet = append(blockSet, slack.NewDividerBlock())
	blockSet = append(blockSet, BuildFifthStepSectionBlock())
	blockSet = append(blockSet, f.BuildRequestTTLBlock()...)
	return slack.Blocks{BlockSet: blockSet}
}

func (f *fifthStepDateBuilder) BuildStartDateBlock() []slack.Block {
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

func (f *fifthStepDateBuilder) BuildAccessDurationBlock() []slack.Block {
	if f.payload.SelectedAccessDurationOptionID == util.ARequestAccessDurationFirstOption {
		text := BuildAccessDurationFirstOptInfoText(f.payload.SelectedAccessDurationOptionName)
		accessDurationOptInfoBlock := slack.NewSectionBlock(
			slack.NewTextBlockObject(util.Markdown, text, false, false),
			nil,
			nil,
		)
		return []slack.Block{accessDurationOptInfoBlock}
	}
	text := BuildAccessDurationSecondOptInfoText(f.payload.SelectedAccessDurationOptionName, f.payload.TTL)
	accessDurationOptInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	accessDurationDateTimeBlock := slack.NewActionBlock(
		util.ARequestAccessDurationDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestAccessDurationDateBlockActionID),
		slack.NewTimePickerBlockElement(util.ARequestAccessDurationTimeBlockActionID),
	)
	return []slack.Block{accessDurationOptInfoBlock, accessDurationDateTimeBlock}
}

func (f *fifthStepDateBuilder) BuildRequestTTLBlock() []slack.Block {
	text := BuildRequestTTLSecondOptInfoText(f.payload.RequestTTLOptionName, f.payload.RequestTTL)
	requestTTLOptInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	requestTTLDateTimeBlock := slack.NewActionBlock(
		util.ARequestRequestTTLDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestRequestTTLDateBlockActionID),
	)
	return []slack.Block{requestTTLOptInfoBlock, requestTTLDateTimeBlock}
}

func (f *fifthStepDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.RequestTTLDateSelectPrivateMetadataPayload{
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
		SelectedAccessDurationOptionID:   f.payload.SelectedAccessDurationOptionID,
		SelectedAccessDurationOptionName: f.payload.SelectedAccessDurationOptionName,
		SelectedAccessDurationDate:       f.payload.SelectedAccessDurationDate,
		SelectedAccessDurationTime:       f.payload.SelectedAccessDurationTime,
		RequestTTL:                       f.payload.RequestTTL,
		SelectedRequestTTLOptionID:       f.payload.RequestTTLOptionID,
		SelectedRequestTTLOptionName:     f.payload.RequestTTLOptionName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
