package accessrequest

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/util"
)

type thirdStepBuilder struct {
	payload *blockactions.ChannelSelect
}

func NewThirdStepBuilder(p *blockactions.ChannelSelect) modal.Builder {
	return &thirdStepBuilder{payload: p}
}

func (f *thirdStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.ARTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          slack.NewTextBlockObject(util.PlainText, util.Submit, false, false),
		CallbackID:      util.ARCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *thirdStepBuilder) BuildBlocks() slack.Blocks {
	thirdStep := BuildThirdStepSectionBlock()
	summarySection := f.BuildSummaryBlock()
	reasonBlock := f.BuildReasonBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			thirdStep,
			summarySection,
			slack.NewDividerBlock(),
			reasonBlock,
		},
	}
	return blocks
}

func (f *thirdStepBuilder) BuildSummaryBlock() *slack.SectionBlock {
	text := "```\n"
	text += fmt.Sprintf("🙋 Requester         : %s\n", f.payload.RequesterRealName)
	text += fmt.Sprintf("💬 Requester Channel : #%s\n", f.payload.RequesterChannelName)
	text += "\n"
	text += fmt.Sprintf("🏷️ Requested Role    : %s\n", f.payload.SelectedRole)
	text += fmt.Sprintf("📥 Reviewres Channel : %s", f.payload.ChannelName)
	text += "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", text, false, false),
		nil, nil,
	)
}

func (f *thirdStepBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(util.PlainText, util.ARReasonElemBlockTest, false, false),
		util.ARReasonElemBlockActionID,
	)
	reasonBlock := slack.NewInputBlock(
		util.ARReasonBlockID,
		slack.NewTextBlockObject(util.PlainText, util.ARReasonBlockText, false, false),
		nil,
		reasonElement,
	)

	if !f.payload.RequireReason {
		reasonBlock.Optional = true
	}
	return reasonBlock
}

func (f *thirdStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.SummaryPrivateMetadataPayload{
		ChannelID:           f.payload.RequesterChannelID,
		ChannelName:         f.payload.RequesterChannelName,
		RealName:            f.payload.RequesterRealName,
		RequireReason:       f.payload.RequireReason,
		SelectedRole:        f.payload.SelectedRole,
		SelectedChannelID:   f.payload.ChannelID,
		SelectedChannelName: f.payload.ChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
