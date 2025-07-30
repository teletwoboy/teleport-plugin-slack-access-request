package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"

	"github.com/slack-go/slack"
)

type summaryBuilder struct {
	payload *accesspolicy.EffectSelect
}

func NewSummaryBuilder(p *accesspolicy.EffectSelect) modal.Builder {
	return &summaryBuilder{
		payload: p,
	}
}

func (s *summaryBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := s.BuildBlocks()
	privateMetadata, err := s.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(modal.PlainText, modal.APTitle, false, false),
		Close:           slack.NewTextBlockObject(modal.PlainText, modal.Back, false, false),
		Submit:          slack.NewTextBlockObject(modal.PlainText, modal.Submit, false, false),
		CallbackID:      modal.APCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (s *summaryBuilder) BuildBlocks() slack.Blocks {
	sixthStep := BuildSixthStepSectionBlock()
	summary := s.BuildSummaryBlock()
	title := s.BuildTitleBlock()
	reason := s.BuildReasonBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			sixthStep,
			summary,
			slack.NewDividerBlock(),
			title,
			slack.NewDividerBlock(),
			reason,
		},
	}
	return blocks
}

func (s *summaryBuilder) BuildSummaryBlock() *slack.SectionBlock {
	var selectedChannelName string
	if s.payload.SelectedChannelName != modal.APAllOption {
		selectedChannelName = "#" + s.payload.SelectedChannelName
	} else {
		selectedChannelName = s.payload.SelectedChannelName
	}

	text := fmt.Sprintf(
		"🙋 Requester         : %s\n"+
			"💬 Requester Channel : #%s\n"+
			"\n"+
			"📥 Target Channel    : %s\n"+
			"🏷️ Target Role       : %s\n"+
			"👤 Target User       : %s\n"+
			"\n"+
			"🌍 Time Zone         : %s\n"+
			"🕐 Start Date        : %s\n"+
			"🕐 End Date          : %s\n"+
			"\n"+
			"⚙️ Effect            : %s",
		s.payload.RequesterRealName,
		s.payload.RequesterChannelName,
		selectedChannelName,
		s.payload.SelectedRoleName,
		s.payload.SelectedRealName,
		s.payload.SelectedTimeZone,
		s.payload.SelectedStartDate+" "+s.payload.SelectedStartTime,
		s.payload.SelectedEndDate+" "+s.payload.SelectedEndTime,
		s.payload.Effect,
	)

	section := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("```\n%s\n```", text), false, false),
		nil, nil,
	)
	return section
}

func (s *summaryBuilder) BuildTitleBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(modal.PlainText, modal.APTitleElemBlockText, false, false),
		modal.APTitleElemBlockActionID,
	)
	reasonBlock := slack.NewInputBlock(
		modal.APTitleBlockID,
		slack.NewTextBlockObject(modal.PlainText, modal.APTitleBlockText, false, false),
		nil,
		reasonElement,
	)
	return reasonBlock
}

func (s *summaryBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(modal.PlainText, modal.APReasonElemBlockText, false, false),
		modal.APReasonElemBlockActionID,
	)
	reasonBlock := slack.NewInputBlock(
		modal.APReasonBlockID,
		slack.NewTextBlockObject(modal.PlainText, modal.APReasonBlockText, false, false),
		nil,
		reasonElement,
	)
	return reasonBlock
}

func (s *summaryBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.SummaryPrivateMetadataPayload{
		ChannelID:           s.payload.RequesterChannelID,
		ChannelName:         s.payload.RequesterChannelName,
		RealName:            s.payload.RequesterRealName,
		SelectedChannelID:   s.payload.SelectedChannelID,
		SelectedChannelName: s.payload.SelectedChannelName,
		SelectedRole:        s.payload.SelectedRole,
		SelectedRoleName:    s.payload.SelectedRoleName,
		SelectedUserID:      s.payload.SelectedUserID,
		SelectedRealName:    s.payload.SelectedRealName,
		SelectedTimeZone:    s.payload.SelectedTimeZone,
		SelectedStartDate:   s.payload.SelectedStartDate,
		SelectedStartTime:   s.payload.SelectedStartTime,
		SelectedEndDate:     s.payload.SelectedEndDate,
		SelectedEndTime:     s.payload.SelectedEndTime,
		SelectedEffect:      s.payload.Effect,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
