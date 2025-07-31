package accessrequest

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/teleport/types"
	"teleport-plugin-slack-access-request/internal/util"
)

type firstStepBuilder struct {
	accessInfo *types.UserAccessInfo
	payload    *slashcommands.AccessRole
	slackUser  *models.User
}

func NewFirstStepBuilder(a *types.UserAccessInfo, p *slashcommands.AccessRole, s *models.User) modal.Builder {
	return &firstStepBuilder{
		accessInfo: a,
		payload:    p,
		slackUser:  s,
	}
}

func (f *firstStepBuilder) Build() (*slack.ModalViewRequest, error) {
	if len(f.accessInfo.Roles) == 0 {
		return nil, fmt.Errorf("<%s> does not have any Role to request. Please contact the administrator", f.slackUser.RealName)
	}
	blocks := f.BuildBlocks()
	privateMetadata, err := f.BuildPrivateMetadata()
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

func (f *firstStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock(f.slackUser.RealName)
	roleBlock := f.BuildRoleBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			roleBlock,
		},
	}
	return blocks
}

func (f *firstStepBuilder) BuildRoleBlock() *slack.ActionBlock {
	roleOpts := f.BuildRoleOpts()
	return slack.NewActionBlock(
		util.ARRoleActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.ARRoleOptionBlockActionID,
			roleOpts...,
		),
	)
}

func (f *firstStepBuilder) BuildRoleOpts() []*slack.OptionBlockObject {
	var roleOpts []*slack.OptionBlockObject
	for _, role := range f.accessInfo.Roles {
		r := role
		roleOpts = append(roleOpts, slack.NewOptionBlockObject(
			r,
			slack.NewTextBlockObject("plain_text", role, false, false),
			nil,
		))
	}
	return roleOpts
}

func (f *firstStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.RoleSelectPrivateMetadataPayload{
		ChannelID:     f.payload.ChannelID,
		ChannelName:   f.payload.ChannelName,
		RealName:      f.slackUser.RealName,
		RequireReason: f.accessInfo.RequireReason,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
