package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"

	"github.com/slack-go/slack"
)

type secondStepBuilder struct {
	payload *accesspolicy.ChannelSelect
	roles   map[string]struct{}
}

func NewSecondStepBuilder(p *accesspolicy.ChannelSelect, r map[string]struct{}) modal.Builder {
	return &secondStepBuilder{
		payload: p,
		roles:   r,
	}
}

func (s *secondStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := s.BuildBlocks()
	privateMetadata, err := s.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(modal.PlainText, modal.APTitle, false, false),
		Close:           slack.NewTextBlockObject(modal.PlainText, "Close", false, false),
		Submit:          nil,
		CallbackID:      modal.APCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (s *secondStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock()
	channelBlock := s.BuildChannelBlock()
	secondStep := BuildSecondStepSectionBlock()
	roleBlock := s.BuildRoleBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			channelBlock,
			slack.NewDividerBlock(),
			secondStep,
			roleBlock,
		},
	}
	return blocks
}

func (s *secondStepBuilder) BuildChannelBlock() *slack.SectionBlock {
	text := "```\n" + s.payload.ChannelName + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (s *secondStepBuilder) BuildRoleBlock() *slack.ActionBlock {
	roleOpts := s.BuildRoleOpts()
	return slack.NewActionBlock(
		modal.APRoleActionBlockID,
		slack.NewOptionsSelectBlockElement(
			modal.StaticSelect,
			slack.NewTextBlockObject(modal.PlainText, modal.SelectOne, false, false),
			modal.APRoleOptionBlockActionID,
			roleOpts...,
		),
	)
}

func (s *secondStepBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	roleOpts = append(roleOpts, slack.NewOptionBlockObject(
		modal.APAllOptionValue,
		slack.NewTextBlockObject(modal.PlainText, modal.APAllOption, false, false),
		nil,
	))
	for r := range s.roles {
		copiedRole := r
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			copiedRole,
			slack.NewTextBlockObject(modal.PlainText, copiedRole, false, false),
			nil,
		))
	}
	return roleOpts
}

func (s *secondStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.RoleSelectPrivateMetadataPayload{
		ChannelID:           s.payload.RequesterChannelID,
		ChannelName:         s.payload.RequesterChannelName,
		RealName:            s.payload.RequesterRealName,
		SelectedChannelID:   s.payload.ChannelID,
		SelectedChannelName: s.payload.ChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
