package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type firstStepBuilder struct {
	channels  []slack.Channel
	payload   *slashcommands.AccessPolicy
	slackUser *models.User
}

func NewFirstStepBuilder(c []slack.Channel, p *slashcommands.AccessPolicy, s *models.User) modal.Builder {
	return &firstStepBuilder{
		channels:  c,
		payload:   p,
		slackUser: s,
	}
}

func (f *firstStepBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(f.channels) == 0 {
		return nil, fmt.Errorf("no channels found. Please contact the administrator")
	}
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.APTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, "Close", false, false),
		Submit:          nil,
		CallbackID:      util.APCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (f *firstStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock()
	channelBlock := f.BuildChannelBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			channelBlock,
		},
	}
	return blocks
}

func (f *firstStepBuilder) BuildChannelBlock() *slack.ActionBlock {
	channelOpts := f.BuildChannelOpts()
	return slack.NewActionBlock(
		util.APChannelActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.APChannelOptionBlockActionID,
			channelOpts...,
		),
	)
}

func (f *firstStepBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOpts []*slack.OptionBlockObject
	channelOpts = append(channelOpts, slack.NewOptionBlockObject(
		util.APAllOptionValue,
		slack.NewTextBlockObject(util.PlainText, util.APAllOption, false, false),
		nil,
	))
	for _, channel := range f.channels {
		copiedChannel := channel
		channelOpts = append(channelOpts, slack.NewOptionBlockObject(
			copiedChannel.ID,
			slack.NewTextBlockObject(util.PlainText, copiedChannel.Name, false, false),
			nil,
		))
	}
	return channelOpts
}

func (f *firstStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.ChannelSelectPrivateMetadataPayload{
		ChannelID:   f.payload.ChannelID,
		ChannelName: f.payload.ChannelName,
		RealName:    f.slackUser.RealName,
		TimeZone:    f.slackUser.TimeZone,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
