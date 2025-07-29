package modal

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
)

type firstStepBuilder struct {
	channels  []slack.Channel
	payload   *slashcommands.AccessPolicy
	slackUser *models.User
}

func NewFirstStepBuilder(c []slack.Channel, p *slashcommands.AccessPolicy, s *models.User) Builder {
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
		Title:           slack.NewTextBlockObject(plainText, accessPolicyTitle, false, false),
		Close:           slack.NewTextBlockObject(plainText, "Close", false, false),
		Submit:          nil,
		CallbackID:      accessPolicyCallBackID,
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
		channelActionBlockID,
		slack.NewOptionsSelectBlockElement(
			StaticSelect,
			slack.NewTextBlockObject(plainText, SelectOne, false, false),
			channelOptionBlockActionID,
			channelOpts...,
		),
	)
}

func (f *firstStepBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOpts []*slack.OptionBlockObject
	channelOpts = append(channelOpts, slack.NewOptionBlockObject(
		accessPolicyAllOptionValue,
		slack.NewTextBlockObject(plainText, accessPolicyAllOption, false, false),
		nil,
	))
	for _, channel := range f.channels {
		copiedChannel := channel
		channelOpts = append(channelOpts, slack.NewOptionBlockObject(
			copiedChannel.ID,
			slack.NewTextBlockObject(plainText, copiedChannel.Name, false, false),
			nil,
		))
	}
	return channelOpts
}

func (f *firstStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyChannelSelectPrivateMetadataPayload{
		ChannelID:   f.payload.ChannelID,
		ChannelName: f.payload.ChannelName,
		RealName:    f.slackUser.RealName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
