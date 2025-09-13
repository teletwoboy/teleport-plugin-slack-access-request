package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type thirdStepBuilder struct {
	payload *accessrequest.ChannelSelect
}

func NewThirdStepBuilder(p *accessrequest.ChannelSelect) modal.Builder {
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
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, t.BuildStartDateBlock()...)
	return slack.Blocks{BlockSet: blockSet}
}

func (t *thirdStepBuilder) BuildStartDateBlock() []slack.Block {
	text := BuildStartDateInfoText()
	startDateInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	startDateOpts := t.BuildStartDateOpts()
	startDateOptsBlock := slack.NewActionBlock(
		util.ARequestStartDateOptionActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARequestStartDateOptionOptionBlockActionID,
			startDateOpts...,
		),
	)
	return []slack.Block{startDateInfoBlock, startDateOptsBlock}
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
	privateMetadata := &accessrequest.StartDateOptionSelectPrivateMetadataPayload{
		ChannelID:           t.payload.RequesterChannelID,
		ChannelName:         t.payload.RequesterChannelName,
		RealName:            t.payload.RequesterRealName,
		RequireReason:       t.payload.RequireReason,
		SelectedRole:        t.payload.SelectedRole,
		SelectedChannelID:   t.payload.ChannelID,
		SelectedChannelName: t.payload.ChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
