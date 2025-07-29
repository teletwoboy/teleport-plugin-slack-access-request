package modal

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

type selectStartDateBuilder struct {
	payload *blockactions.AccessPolicyUserSelect
}

func NewSelectStartDateBuilder(p *blockactions.AccessPolicyUserSelect) Builder {
	return &selectStartDateBuilder{
		payload: p,
	}
}

func (a *selectStartDateBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (a *selectStartDateBuilder) BuildBlocks() slack.Blocks {
	durationBlockLabel := fmt.Sprintf("*Step 4 of 6 - Duration*")
	startDateBlockLabel := fmt.Sprintf("4-1. Start Date/Time")
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
			),
		},
	}
	return blocks
}

func (a *selectStartDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyStartDateSelectPrivateMetadataPayload{
		ChannelID:           a.payload.RequesterChannelID,
		ChannelName:         a.payload.RequesterChannelName,
		RealName:            a.payload.RequesterRealName,
		SelectedChannelID:   a.payload.SelectedChannelID,
		SelectedChannelName: a.payload.SelectedChannelName,
		SelectedRole:        a.payload.SelectedRole,
		SelectedRoleName:    a.payload.SelectedRoleName,
		SelectedUserID:      a.payload.UserID,
		SelectedRealName:    a.payload.RealName,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -----------------------------------------------------------------------------------

type selectStartTimeBuilder struct {
	payload *blockactions.AccessPolicyStartDateSelect
}

func NewSelectStartTimeBuilder(p *blockactions.AccessPolicyStartDateSelect) Builder {
	return &selectStartTimeBuilder{
		payload: p,
	}
}

func (a *selectStartTimeBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (a *selectStartTimeBuilder) BuildBlocks() slack.Blocks {
	durationBlockLabel := fmt.Sprintf("*Step 4 of 6 - Duration*")
	startDateBlockLabel := fmt.Sprintf("4-1. Start Date/Time")
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
		},
	}
	return blocks
}

func (a *selectStartTimeBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyStartTimeSelectPrivateMetadataPayload{
		ChannelID:           a.payload.RequesterChannelID,
		ChannelName:         a.payload.RequesterChannelName,
		RealName:            a.payload.RequesterRealName,
		SelectedChannelID:   a.payload.SelectedChannelID,
		SelectedChannelName: a.payload.SelectedChannelName,
		SelectedRole:        a.payload.SelectedRole,
		SelectedRoleName:    a.payload.SelectedRoleName,
		SelectedUserID:      a.payload.SelectedUserID,
		SelectedRealName:    a.payload.SelectedRealName,
		SelectedStartDate:   a.payload.StartDate,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -----------------------------------------------------------------------------------

type selectEndDateBuilder struct {
	payload *blockactions.AccessPolicyStartTimeSelect
}

func NewSelectEndDateBuilder(p *blockactions.AccessPolicyStartTimeSelect) Builder {
	return &selectEndDateBuilder{
		payload: p,
	}
}

func (a *selectEndDateBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (a *selectEndDateBuilder) BuildBlocks() slack.Blocks {
	durationBlockLabel := fmt.Sprintf("*Step 4 of 6 - Duration*")
	startDateBlockLabel := fmt.Sprintf("4-1. Start Date/Time")
	endDateBlockLabel := fmt.Sprintf("4-2. End Date/Time")
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
			),
		},
	}
	return blocks
}

func (a *selectEndDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyEndDateSelectPrivateMetadataPayload{
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
		SelectedStartTime:   a.payload.StartTime,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -----------------------------------------------------------------------------------

type selectEndTimeBuilder struct {
	payload *blockactions.AccessPolicyEndDateSelect
}

func NewSelectEndTimeBuilder(p *blockactions.AccessPolicyEndDateSelect) Builder {
	return &selectEndTimeBuilder{
		payload: p,
	}
}

func (a *selectEndTimeBuilder) Build() (*slack.ModalViewRequest, error) {
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

func (a *selectEndTimeBuilder) BuildBlocks() slack.Blocks {
	durationBlockLabel := fmt.Sprintf("*Step 4 of 6 - Duration*")
	startDateBlockLabel := fmt.Sprintf("4-1. Start Date/Time")
	endDateBlockLabel := fmt.Sprintf("4-2. End Date/Time")
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
		},
	}
	return blocks
}

func (a *selectEndTimeBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &blockactions.AccessPolicyEndTimeSelectPrivateMetadataPayload{
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
		SelectedEndDate:     a.payload.EndDate,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}

// -----------------------------------------------------------------------------------
