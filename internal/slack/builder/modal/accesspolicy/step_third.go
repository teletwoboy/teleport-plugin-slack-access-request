package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type thirdStepBuilder struct {
	payload    *accesspolicy.RoleSelect
	slackUsers []models.User
}

func NewThirdStepBuilder(p *accesspolicy.RoleSelect, s []models.User) modal.Builder {
	return &thirdStepBuilder{
		payload:    p,
		slackUsers: s,
	}
}

func (t *thirdStepBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := t.BuildBlocks()
	privateMetadata, err := t.BuildPrivateMetadata()
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

func (t *thirdStepBuilder) BuildBlocks() slack.Blocks {
	firstStep := BuildFirstStepSectionBlock()
	channelBlock := t.BuildChannelBlock()
	secondStep := BuildSecondStepSectionBlock()
	roleBlock := t.BuildRoleBlock()
	thirdStep := BuildThirdStepSectionBlock()
	userBlock := t.BuildUserBlock()
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			firstStep,
			channelBlock,
			slack.NewDividerBlock(),
			secondStep,
			roleBlock,
			slack.NewDividerBlock(),
			thirdStep,
			userBlock,
		},
	}
	return blocks
}

func (t *thirdStepBuilder) BuildChannelBlock() *slack.SectionBlock {
	text := "```\n" + t.payload.SelectedChannelName + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (t *thirdStepBuilder) BuildRoleBlock() *slack.SectionBlock {
	text := "```\n" + t.payload.RoleName + "\n```"
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
}

func (t *thirdStepBuilder) BuildUserBlock() *slack.ActionBlock {
	userOpts := t.BuildUserOpts()
	return slack.NewActionBlock(
		util.APUserActionBlockID,
		slack.NewOptionsSelectBlockElement(
			util.StaticSelect,
			slack.NewTextBlockObject(util.PlainText, util.SelectOne, false, false),
			util.APUserOptionBlockActionID,
			userOpts...,
		),
	)
}

func (t *thirdStepBuilder) BuildUserOpts() []*slack.OptionBlockObject {
	var userOpts []*slack.OptionBlockObject
	userOpts = append(userOpts, slack.NewOptionBlockObject(
		util.APAllOptionValue,
		slack.NewTextBlockObject(util.PlainText, util.APAllOption, false, false),
		nil,
	))
	for _, u := range t.slackUsers {
		copiedUser := u
		userOpts = append(userOpts, slack.NewOptionBlockObject(
			copiedUser.ID,
			slack.NewTextBlockObject(util.PlainText, copiedUser.RealName, false, false),
			nil,
		))
	}
	return userOpts
}

func (t *thirdStepBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accesspolicy.UserSelectPrivateMetadataPayload{
		ChannelID:           t.payload.RequesterChannelID,
		ChannelName:         t.payload.RequesterChannelName,
		RealName:            t.payload.RequesterRealName,
		TimeZone:            t.payload.RequesterTimeZone,
		SelectedChannelID:   t.payload.SelectedChannelID,
		SelectedChannelName: t.payload.SelectedChannelName,
		SelectedRole:        t.payload.Role,
		SelectedRoleName:    t.payload.RoleName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
