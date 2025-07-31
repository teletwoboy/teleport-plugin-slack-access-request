package accessrequest

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	slacktypes "teleport-plugin-slack-access-request/internal/slack/types"
	"teleport-plugin-slack-access-request/internal/util"
)

type secondStepBuilder struct {
	channels []slacktypes.ReviewersChannel
	payload  *blockactions.RoleSelect
}

func NewSecondStepBuilder(c []slacktypes.ReviewersChannel, p *blockactions.RoleSelect) modal.Builder {
	return &secondStepBuilder{
		channels: c,
		payload:  p,
	}
}

func (s *secondStepBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(s.channels) == 0 {
		return nil, fmt.Errorf("role <%s> does not have any Reviewers Channel. Please contact the administrator", s.payload.Role)
	}
	blocks := s.BuildBlocks()
	privateMetadata, err := s.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.ARTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Close, false, false),
		Submit:          nil,
		CallbackID:      util.ARCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (s *secondStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock(s.payload.RequesterRealName)
	roleBlock := s.BuildRoleBlock()
	secondStep := BuildSecondStepSectionBlock()
	channelBlock := s.BuildChannelBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			roleBlock,
			slack.NewDividerBlock(),
			secondStep,
			channelBlock,
		},
	}
	return blocks
}

func (s *secondStepBuilder) BuildRoleBlock() *slack.SectionBlock {
	text := "```\n" + s.payload.Role + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (s *secondStepBuilder) BuildChannelBlock() *slack.ActionBlock {
	channelOpts := s.BuildChannelOpts()
	return slack.NewActionBlock(
		util.ARChannelActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARChannelOptionBlockActionID,
			channelOpts...,
		),
	)
}

func (s *secondStepBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOptions []*slack.OptionBlockObject
	for _, ch := range s.channels {
		id := ch.ID
		label := ch.Name
		channelOptions = append(channelOptions, slack.NewOptionBlockObject(
			id,
			slack.NewTextBlockObject(util.PlainText, label, false, false),
			nil,
		))
	}
	return channelOptions
}

func (s *secondStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.ChannelSelectPrivateMetadataPayload{
		ChannelID:     s.payload.RequesterChannelID,
		ChannelName:   s.payload.RequesterChannelName,
		RealName:      s.payload.RequesterRealName,
		RequireReason: s.payload.RequireReason,
		SelectedRole:  s.payload.Role,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
