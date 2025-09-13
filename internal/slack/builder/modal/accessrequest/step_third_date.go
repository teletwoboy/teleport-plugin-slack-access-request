package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"time"

	"github.com/slack-go/slack"
)

type thirdStepDateBuilder struct {
	payload *accessrequest.StartDateOptionSelect
	ttl     time.Time
}

func NewThirdStepDateBuilder(p *accessrequest.StartDateOptionSelect, t time.Time) modal.Builder {
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
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, t.BuildStartDateBlock()...)
	return slack.Blocks{
		BlockSet: blockSet,
	}
}

func (t *thirdStepDateBuilder) BuildStartDateBlock() []slack.Block {
	text := BuildStartDateSecondOptInfoText(t.payload.StartDateOptionName, t.ttl)
	startDateSecondOptInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	startDateTimeBlock := slack.NewActionBlock(
		util.ARequestStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestStartDateBlockActionID),
	)
	return []slack.Block{
		startDateSecondOptInfoBlock,
		startDateTimeBlock,
	}
}

func (t *thirdStepDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.StartDateSelectPrivateMetadataPayload{
		ChannelID:                   t.payload.RequesterChannelID,
		ChannelName:                 t.payload.RequesterChannelName,
		RealName:                    t.payload.RequesterRealName,
		RequireReason:               t.payload.RequireReason,
		SelectedRole:                t.payload.SelectedRole,
		SelectedChannelID:           t.payload.SelectedChannelID,
		SelectedChannelName:         t.payload.SelectedChannelName,
		SelectedStartDateOptionID:   t.payload.StartDateOptionID,
		SelectedStartDateOptionName: t.payload.StartDateOptionName,
		TTL:                         t.ttl,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
