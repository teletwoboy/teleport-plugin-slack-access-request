package modal

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

type selectEffectBuilder struct {
	payload *blockactions.AccessPolicyEndTimeSelect
}

func NewSelectEffectBuilder(p *blockactions.AccessPolicyEndTimeSelect) Builder {
	return &selectEffectBuilder{
		payload: p,
	}
}

func (a *selectEffectBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := a.BuildBlocks()
	privateMetadata, err := a.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject("plain_text", "Access Policy", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Back", false, false),
		Submit:          nil,
		CallbackID:      "access_policy_modal",
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (a *selectEffectBuilder) BuildBlocks() slack.Blocks {
	durationBlockLabel := fmt.Sprintf("*Step 4 of 6 - Duration*")
	startDateBlockLabel := fmt.Sprintf("4-1. Start Date/Time")
	endDateBlockLabel := fmt.Sprintf("4-2. End Date/Time")
	effectBlockLabel := fmt.Sprintf("*Step 5 of 6 - Effect*")
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", durationBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", startDateBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"start_date_time_block",
				slack.NewDatePickerBlockElement("access_policy_start_date_select"),
				slack.NewTimePickerBlockElement("access_policy_start_time_select"),
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", endDateBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"end_date_time_block",
				slack.NewDatePickerBlockElement("access_policy_end_date_select"),
				slack.NewTimePickerBlockElement("access_policy_end_time_select"),
			),
			slack.NewDividerBlock(),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", effectBlockLabel, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"effect_block",
				slack.NewButtonBlockElement(
					"access_policy_effect_allow_select",
					"allow",
					slack.NewTextBlockObject("plain_text", "✅ Allow", false, false),
				),
				slack.NewButtonBlockElement(
					"access_policy_effect_deny_select",
					"deny",
					slack.NewTextBlockObject("plain_text", "⛔ Deny", false, false),
				),
			),
		},
	}
	return blocks
}

func (a *selectEffectBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyEffectSelectPrivateMetadataPayload{
		ChannelID:           a.payload.RequesterChannelID,
		ChannelName:         a.payload.RequesterChannelName,
		RealName:            a.payload.RequesterRealName,
		SelectedChannelID:   a.payload.SelectedChannelID,
		SelectedChannelName: a.payload.SelectedChannelName,
		SelectedRole:        a.payload.SelectedRole,
		SelectedRoleName:    a.payload.SelectedRoleName,
		SelectedUserID:      a.payload.SelectedUserID,
		SelectedRealName:    a.payload.SelectedRealName,
		SelectedStartDate:   a.payload.SelectedStartDate,
		SelectedStartTime:   a.payload.SelectedStartTime,
		SelectedEndDate:     a.payload.SelectedEndDate,
		SelectedEndTime:     a.payload.EndTime,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
