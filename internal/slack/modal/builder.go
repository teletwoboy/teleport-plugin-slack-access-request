package modal

import (
	"fmt"
	slacktypes "teleport-plugin-slack-access-request/internal/slack/types"
	teleporttypes "teleport-plugin-slack-access-request/internal/teleport/types"

	"github.com/slack-go/slack"
)

// Builder 또한 MessageBuilder 와 동일
type Builder interface {
	Build() slack.ModalViewRequest
}

type AccessRequestBuilder struct {
	AccessInfo *teleporttypes.UserAccessInfo
	Channels   []slacktypes.ReviewersChannel
	ChannelID  string
}

func NewAccessRequestBuilder(a *teleporttypes.UserAccessInfo, c []slacktypes.ReviewersChannel, cID string) *AccessRequestBuilder {
	return &AccessRequestBuilder{
		AccessInfo: a,
		Channels:   c,
		ChannelID:  cID,
	}
}

func (a *AccessRequestBuilder) Build() slack.ModalViewRequest {
	blocks := a.BuildBlocks()

	modal := slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject("plain_text", "Access Request", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "닫기", false, false),
		Submit:          slack.NewTextBlockObject("plain_text", "요청", false, false),
		CallbackID:      "access_request_modal",
		Blocks:          blocks,
		PrivateMetadata: a.ChannelID,
	}

	return modal
}

func (a *AccessRequestBuilder) BuildBlocks() slack.Blocks {
	roleOptions := a.BuildRoleOpts()
	channelOptions := a.BuildChannelOpts()
	reasonBlock := a.BuildReasonBlock()

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewInputBlock(
				"role_block",
				slack.NewTextBlockObject("plain_text", "요청할 Role", false, false),
				nil,
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject("plain_text", "선택하세요", false, false),
					"role_select",
					roleOptions...,
				),
			),
			slack.NewInputBlock(
				"channel_block",
				slack.NewTextBlockObject("plain_text", "요청 보낼 채널", false, false),
				nil,
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject("plain_text", "선택하세요", false, false),
					"channel_select",
					channelOptions...,
				),
			),
		},
	}

	blocks.BlockSet = append(blocks.BlockSet, reasonBlock)
	return blocks
}

func (a *AccessRequestBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	for _, role := range a.AccessInfo.Roles {
		r := role
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			r,
			slack.NewTextBlockObject("plain_text", role, false, false),
			nil,
		))
	}
	return roleOpts
}

func (a *AccessRequestBuilder) BuildChannelOpts() []*slack.OptionBlockObject {
	var channelOptions []*slack.OptionBlockObject
	for _, ch := range a.Channels {
		id := ch.ID
		name := ch.Name
		label := fmt.Sprintf("%s", name)
		channelOptions = append(channelOptions, slack.NewOptionBlockObject(
			id,
			slack.NewTextBlockObject("plain_text", label, false, false),
			nil,
		))
	}
	return channelOptions
}

func (a *AccessRequestBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonBlock := slack.NewInputBlock(
		"reason_block",
		slack.NewTextBlockObject("plain_text", "요청 이유를 입력하세요", false, false),
		nil,
		slack.NewPlainTextInputBlockElement(
			slack.NewTextBlockObject("plain_text", "입력하세요", false, false),
			"reason_input",
		),
	)
	if !a.AccessInfo.RequireReason {
		reasonBlock.Optional = true
	}
	return reasonBlock
}
