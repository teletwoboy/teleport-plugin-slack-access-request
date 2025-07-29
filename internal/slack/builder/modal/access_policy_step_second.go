package modal

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

type secondStepBuilder struct {
	payload *blockactions.AccessPolicyChannelSelect
	roles   map[string]struct{}
}

func NewSecondStepBuilder(p *blockactions.AccessPolicyChannelSelect, r map[string]struct{}) Builder {
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
		Title:           slack.NewTextBlockObject(plainText, accessPolicyTitle, false, false),
		Close:           slack.NewTextBlockObject(plainText, "Close", false, false),
		Submit:          nil,
		CallbackID:      accessPolicyCallBackID,
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
		slack.NewTextBlockObject(Markdown, text, false, false),
		nil,
		nil,
	)
}

func (s *secondStepBuilder) BuildRoleBlock() *slack.ActionBlock {
	roleOpts := s.BuildRoleOpts()
	return slack.NewActionBlock(
		roleActionBlockID,
		slack.NewOptionsSelectBlockElement(
			StaticSelect,
			slack.NewTextBlockObject(plainText, SelectOne, false, false),
			roleOptionBlockActionID,
			roleOpts...,
		),
	)
}

func (s *secondStepBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	roleOpts = append(roleOpts, slack.NewOptionBlockObject(
		accessPolicyAllOptionValue,
		slack.NewTextBlockObject(plainText, accessPolicyAllOption, false, false),
		nil,
	))
	for r := range s.roles {
		copiedRole := r
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			copiedRole,
			slack.NewTextBlockObject(plainText, copiedRole, false, false),
			nil,
		))
	}
	return roleOpts
}

func (s *secondStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyRoleSelectPrivateMetadataPayload{
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
