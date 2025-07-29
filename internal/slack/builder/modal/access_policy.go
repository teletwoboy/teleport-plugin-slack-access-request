package modal

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

type selectRoleBuilder struct {
	channels []slack.Channel
	payload  *blockactions.AccessPolicyChannelSelect
	roles    map[string]struct{}
}

func NewSelectRoleBuilder(c []slack.Channel, p *blockactions.AccessPolicyChannelSelect, r map[string]struct{}) Builder {
	return &selectRoleBuilder{
		channels: c,
		payload:  p,
		roles:    r,
	}
}

func (a *selectRoleBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(a.channels) == 0 {
		return nil, fmt.Errorf("no channels found. Please contact the administrator")
	}
	blocks := a.BuildBlocks()
	privateMetadata, err := a.BuildPrivateMetadata()
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

func (a *selectRoleBuilder) BuildBlocks() slack.Blocks {
	channelOptions := a.BuildChannelOpts()
	roleOptions := a.BuildRoleOpts()
	channelSectionBlockLabel := fmt.Sprintf("*Target Channel*")
	roleSectionBlockLabel := fmt.Sprintf("*Target Role*")
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", channelSectionBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"channel_block",
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject(plainText, "Select one", false, false),
					"access_policy_channel_select",
					channelOptions...,
				),
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", roleSectionBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"role_block",
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject(plainText, "Select one", false, false),
					"access_policy_role_select",
					roleOptions...,
				),
			),
		},
	}
	return blocks
}

func (a *selectRoleBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOpts []*slack.OptionBlockObject
	channelOpts = append(channelOpts, slack.NewOptionBlockObject(
		"*",
		slack.NewTextBlockObject(plainText, "* (all)", false, false),
		nil,
	))
	for _, c := range a.channels {
		copiedChannel := c
		channelOpts = append(channelOpts, slack.NewOptionBlockObject(
			copiedChannel.ID,
			slack.NewTextBlockObject(plainText, copiedChannel.Name, false, false),
			nil,
		))
	}
	return channelOpts
}

func (a *selectRoleBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	roleOpts = append(roleOpts, slack.NewOptionBlockObject(
		"*",
		slack.NewTextBlockObject(plainText, "* (all)", false, false),
		nil,
	))
	for r := range a.roles {
		copiedRole := r
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			copiedRole,
			slack.NewTextBlockObject(plainText, copiedRole, false, false),
			nil,
		))
	}
	return roleOpts
}

func (a *selectRoleBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyRoleSelectPrivateMetadataPayload{
		ChannelID:   a.payload.ChannelID,
		ChannelName: a.payload.ChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
