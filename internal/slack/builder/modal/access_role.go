package modal

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/types"

	"github.com/slack-go/slack"
)

type accessRoleBuilder struct {
	accessInfo *types.UserAccessInfo
	payload    *slashcommands.AccessRole
	slackUser  *models.User
}

func NewAccessRoleBuilder(a *types.UserAccessInfo, p *slashcommands.AccessRole, s *models.User) Builder {
	return &accessRoleBuilder{
		accessInfo: a,
		payload:    p,
		slackUser:  s,
	}
}

func (a *accessRoleBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(a.accessInfo.Roles) == 0 {
		return nil, fmt.Errorf("%s does not have any Role to request. Please contact the administrator", a.slackUser.RealName)
	}
	blocks := a.BuildBlocks()
	privateMetadata, err := a.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject("plain_text", "Access Role", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Close", false, false),
		Submit:          nil,
		CallbackID:      "access_role_modal",
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (a *accessRoleBuilder) BuildBlocks() slack.Blocks {
	roleOptions := a.BuildRoleOpts()
	sectionBlockLabel := fmt.Sprintf(":lock: *%s's requestable roles*", a.slackUser.RealName)
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", sectionBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"role_block",
				slack.NewOptionsSelectBlockElement(
					"static_select",
					slack.NewTextBlockObject("plain_text", "Select one", false, false),
					"role_select",
					roleOptions...,
				),
			),
		},
	}
	return blocks
}

func (a *accessRoleBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
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

func (a *accessRoleBuilder) BuildPrivateMetadata() (string, error) {
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
