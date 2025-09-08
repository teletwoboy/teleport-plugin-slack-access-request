package accessrequest

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"time"
)

type thirdStepBuilder struct {
	payload *blockactions.ChannelSelect
}

func NewThirdStepBuilder(p *blockactions.ChannelSelect) modal.Builder {
	return &thirdStepBuilder{
		payload: p,
	}
}

func (t *thirdStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := t.BuildBlocks()
	privateMetadata, err := t.BuildPrivateMetadata()
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

func (t *thirdStepBuilder) BuildBlocks() slack.Blocks {
	thirdStep := BuildThirdStepSectionBlock()
	startDateOptionBlock := t.BuildStartDateOptionBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			thirdStep,
			startDateOptionBlock,
		},
	}
	return blocks
}

func (t *thirdStepBuilder) BuildStartDateOptionBlock() *slack.ActionBlock {
	startDateOpts := t.BuildStartDateOpts()
	return slack.NewActionBlock(
		util.ARequestStartDateOptionActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARequestStartDateOptionOptionBlockActionID,
			startDateOpts...,
		),
	)
}

func (t *thirdStepBuilder) BuildStartDateOpts() []*slack.OptionBlockObject {
	var startDateOpts []*slack.OptionBlockObject
	startDateOpts = append(startDateOpts,
		slack.NewOptionBlockObject(
			util.ARequestStartDateFirstOption,
			slack.NewTextBlockObject(util.PlainText, util.ARequestStartDateFirstOption, false, false),
			nil,
		),
		slack.NewOptionBlockObject(
			util.ARequestStartDateSecondOption,
			slack.NewTextBlockObject(util.PlainText, util.ARequestStartDateSecondOption, false, false),
			nil,
		),
	)
	return startDateOpts
}

func (t *thirdStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.StartDateOptionSelectPrivateMetadataPayload{
		ChannelID:           t.payload.RequesterChannelID,
		ChannelName:         t.payload.RequesterChannelName,
		RealName:            t.payload.RequesterRealName,
		RequireReason:       t.payload.RequireReason,
		SelectedRole:        t.payload.SelectedRole,
		SelectedChannelID:   t.payload.RequesterChannelID,
		SelectedChannelName: t.payload.RequesterChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -------------------------------------------------------------------------

type thirdStepDateBuilder struct {
	payload *blockactions.StartDateOptionSelect
	ttl     time.Time
}

func NewThirdStepDateBuilder(p *blockactions.StartDateOptionSelect, t time.Time) modal.Builder {
	return &thirdStepDateBuilder{
		payload: p,
		ttl:     t,
	}
}

func (t *thirdStepDateBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := t.BuildBlocks()
	privateMetadata, err := t.BuildPrivateMetadata()
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

func (t *thirdStepDateBuilder) BuildBlocks() slack.Blocks {
	thirdStep := BuildThirdStepSectionBlock()
	startDateOptionBlock := t.BuildStartDateOptionBlock()
	startDateTimeBlock := t.BuildStartDateTimeBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			thirdStep,
			startDateOptionBlock,
			startDateTimeBlock,
		},
	}
	return blocks
}

func (t *thirdStepDateBuilder) BuildStartDateOptionBlock() *slack.SectionBlock {
	text := "```\n"
	text += t.payload.StartDateOptionName + "\n"
	text += "\n" + "💡 You can select a time until " + t.ttl.String()
	text += "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (t *thirdStepDateBuilder) BuildStartDateTimeBlock() *slack.ActionBlock {
	return slack.NewActionBlock(
		util.ARequestStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestStartDateBlockActionID),
	)
}

func (t *thirdStepDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.StartDateSelectPrivateMetadataPayload{
		ChannelID:                   t.payload.RequesterChannelID,
		ChannelName:                 t.payload.RequesterChannelName,
		RealName:                    t.payload.RequesterRealName,
		RequireReason:               t.payload.RequireReason,
		SelectedRole:                t.payload.SelectedRole,
		SelectedChannelID:           t.payload.RequesterChannelID,
		SelectedChannelName:         t.payload.RequesterChannelName,
		SelectedStartDateOptionID:   t.payload.StartDateOptionID,
		SelectedStartDateOptionName: t.payload.StartDateOptionName,
		TTL:                         t.ttl.String(),
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
