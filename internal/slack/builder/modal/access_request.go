package modal

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	slacktypes "teleport-plugin-slack-access-request/internal/slack/types"
	teleporttypes "teleport-plugin-slack-access-request/internal/teleport/types"

	"github.com/slack-go/slack"
)

type accessRequestBuilder struct {
	accessInfo *teleporttypes.UserAccessInfo
	channels   []slacktypes.ReviewersChannel
	payload    *slashcommands.AccessRequest
	slackUser  *models.User
}

func NewAccessRequestBuilder(a *teleporttypes.UserAccessInfo, c []slacktypes.ReviewersChannel, p *slashcommands.AccessRequest, s *models.User) Builder {
	return &accessRequestBuilder{
		accessInfo: a,
		channels:   c,
		payload:    p,
		slackUser:  s,
	}
}

func (a *accessRequestBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(a.accessInfo.Roles) == 0 {
		return nil, fmt.Errorf("%s does not have any Role to request. Please contact the administrator", a.slackUser.RealName)
	}
	if len(a.channels) == 0 {
		return nil, fmt.Errorf("no available ReviewersChannel to request. Please contact the administrator")
	}
	blocks := a.BuildBlocks()
	privateMetadata, err := a.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject("plain_text", "Access Request", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Close", false, false),
		Submit:          slack.NewTextBlockObject("plain_text", "Submit", false, false),
		CallbackID:      "access_request_modal",
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (a *accessRequestBuilder) BuildBlocks() slack.Blocks {
	roleOptions := a.BuildRoleOpts()
	channelOptions := a.BuildChannelOpts()
	reasonBlock := a.BuildReasonBlock()

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewInputBlock(
				"role_block",
				slack.NewTextBlockObject("plain_text", "Request Role", false, false),
				nil,
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject("plain_text", "Select one", false, false),
					"role_select",
					roleOptions...,
				),
			),
			slack.NewInputBlock(
				"channel_block",
				slack.NewTextBlockObject("plain_text", "Reviewers Channel", false, false),
				nil,
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject("plain_text", "Select one", false, false),
					"channel_select",
					channelOptions...,
				),
			),
		},
	}

	blocks.BlockSet = append(blocks.BlockSet, reasonBlock)
	return blocks
}

func (a *accessRequestBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	for _, role := range a.accessInfo.Roles {
		r := role
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			r,
			slack.NewTextBlockObject("plain_text", role, false, false),
			nil,
		))
	}
	return roleOpts
}

func (a *accessRequestBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOptions []*slack.OptionBlockObject
	for _, ch := range a.channels {
		id := ch.ID
		label := ch.Name
		channelOptions = append(channelOptions, slack.NewOptionBlockObject(
			id,
			slack.NewTextBlockObject("plain_text", label, false, false),
			nil,
		))
	}
	return channelOptions
}

func (a *accessRequestBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "Enter the reason", false, false),
		"reason_input",
	)
	reasonBlock := slack.NewInputBlock(
		"reason_block",
		slack.NewTextBlockObject("plain_text", "Request Reason", false, false),
		nil,
		reasonElement,
	)

	if !a.accessInfo.RequireReason {
		reasonBlock.Optional = true
	}
	return reasonBlock
}

func (a *accessRequestBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &viewsubmission.AccessRequestModalPrivateMetadataPayload{
		ChannelID:   a.payload.ChannelID,
		ChannelName: a.payload.ChannelName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}

	return string(jsonBytes), nil
}
