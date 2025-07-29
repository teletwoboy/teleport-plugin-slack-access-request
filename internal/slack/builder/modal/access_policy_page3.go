package modal

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

type summaryBuilder struct {
	payload *blockactions.AccessPolicyEffectSelect
}

func NewSummaryBuilder(p *blockactions.AccessPolicyEffectSelect) Builder {
	return &summaryBuilder{
		payload: p,
	}
}

func (a *summaryBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := a.BuildBlocks()
	privateMetadata, err := a.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject("plain_text", "Access Policy", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Back", false, false),
		Submit:          slack.NewTextBlockObject("plain_text", "Submit", false, false),
		CallbackID:      "access_policy_modal",
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (a *summaryBuilder) BuildBlocks() slack.Blocks {
	section := a.BuildSectionBlock()
	reason := a.BuildReasonBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			section,
			slack.NewDividerBlock(),
			reason,
		},
	}
	return blocks
}

func (a *summaryBuilder) BuildSectionBlock() *slack.SectionBlock {
	text := fmt.Sprintf(
		"🙋 Requester         : %s\n"+
			"💬 Requester Channel : #%s\n"+
			"\n"+
			"📥 Target Channel    : #%s\n"+
			"🏷️ Target Role       : %s\n"+
			"👤 Target User       : %s\n"+
			"\n"+
			"🕐 Start Date        : %s\n"+
			"🕐 End Date          : %s\n"+
			"\n"+
			"⚙️ Effect            : %s",
		a.payload.RequesterRealName,
		a.payload.RequesterChannelName,
		a.payload.SelectedChannelName,
		a.payload.SelectedRoleName,
		a.payload.SelectedRealName,
		a.payload.SelectedStartDate+" "+a.payload.SelectedStartTime,
		a.payload.SelectedEndDate+" "+a.payload.SelectedEndTime,
		a.payload.Effect,
	)

	section := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("```\n%s\n```", text), false, false),
		nil, nil,
	)
	return section
}

func (a *summaryBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "Enter the reason", false, false),
		"review_reason",
	)
	reasonBlock := slack.NewInputBlock(
		"reason_input",
		slack.NewTextBlockObject("plain_text", "Review Reason", false, false),
		nil,
		reasonElement,
	)
	return reasonBlock
}

func (a *summaryBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicySummaryPrivateMetadataPayload{
		ChannelID:           a.payload.RequesterChannelID,
		ChannelName:         a.payload.RequesterChannelName,
		RealName:            a.payload.RequesterRealName,
		SelectedChannelID:   a.payload.SelectedChannelID,
		SelectedChannelName: a.payload.SelectedChannelName,
		SelectedRole:        a.payload.SelectedRole,
		SelectedRoleName:    a.payload.SelectedRoleName,
		SelectedUserID:      a.payload.SelectedUserID,
		SelectedRealName:    a.payload.SelectedRealName,
		SelectedStartDate:   a.payload.SelectedStartDate,
		SelectedStartTime:   a.payload.SelectedStartTime,
		SelectedEndDate:     a.payload.SelectedEndDate,
		SelectedEndTime:     a.payload.SelectedEndTime,
		SelectedEffect:      a.payload.Effect,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
