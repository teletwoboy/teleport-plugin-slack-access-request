package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"

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
		Title:           slack.NewTextBlockObject(util.PlainText, util.APTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          slack.NewTextBlockObject(util.PlainText, util.Submit, false, false),
		CallbackID:      util.APCallBackID,
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
	if s.payload.SelectedChannelName != util.APAllOption {
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
			"🕐 Start Date        : %s (UTC)\n"+
			"🕐 End Date          : %s (UTC)\n"+
			"⚙️ Effect            : %s",
		s.payload.RequesterRealName,
		s.payload.RequesterChannelName,
		selectedChannelName,
		s.payload.SelectedRoleName,
		s.payload.SelectedRealName,
		s.payload.SelectedStartDate.Format(util.SecondTimeFormat),
		s.payload.SelectedEndDate.Format(util.SecondTimeFormat),
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
		slack.NewTextBlockObject(util.PlainText, util.APTitleElemBlockText, false, false),
		util.APTitleElemBlockActionID,
	)
	reasonBlock := slack.NewInputBlock(
		util.APTitleBlockID,
		slack.NewTextBlockObject(util.PlainText, util.APTitleBlockText, false, false),
		nil,
		reasonElement,
	)
	return reasonBlock
}

func (s *summaryBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(util.PlainText, util.APReasonElemBlockText, false, false),
		util.APReasonElemBlockActionID,
	)
	reasonBlock := slack.NewInputBlock(
		util.APReasonBlockID,
		slack.NewTextBlockObject(util.PlainText, util.APReasonBlockText, false, false),
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
		TimeZone:            s.payload.RequesterTimeZone,
		SelectedChannelID:   s.payload.SelectedChannelID,
		SelectedChannelName: s.payload.SelectedChannelName,
		SelectedRole:        s.payload.SelectedRole,
		SelectedRoleName:    s.payload.SelectedRoleName,
		SelectedUserID:      s.payload.SelectedUserID,
		SelectedRealName:    s.payload.SelectedRealName,
		SelectedStartDate:   s.payload.SelectedStartDate,
		SelectedEndDate:     s.payload.SelectedEndDate,
		SelectedEffect:      s.payload.Effect,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
