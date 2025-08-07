/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package modal

import (
	"encoding/json"
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type accessReviewBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	slackUser     *slackmodels.User
	channelID     string
}

func NewAccessReviewBuilder(a *teleportmodels.AccessRequest, s *slackmodels.User, cID string) Builder {
	return &accessReviewBuilder{
		accessRequest: a,
		slackUser:     s,
		channelID:     cID,
	}
}

func (a *accessReviewBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := a.BuildBlocks()
	privateMetadata, err := a.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject("plain_text", "Access Review", false, false),
		Close:           slack.NewTextBlockObject("plain_text", "Close", false, false),
		Submit:          slack.NewTextBlockObject("plain_text", "Submit", false, false),
		CallbackID:      "access_review_modal",
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}

	return &modal, nil
}

func (a *accessReviewBuilder) BuildBlocks() slack.Blocks {
	section := a.BuildSectionBlock()
	radioBlock := a.BuildRadioBlock()
	reasonBlock := a.BuildReasonBlock()

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			section,
			radioBlock,
			reasonBlock,
		},
	}
	return blocks
}

func (a *accessReviewBuilder) BuildSectionBlock() *slack.SectionBlock {
	text := fmt.Sprintf(
		"👤 Requester          : %s\n"+
			"💬 Requester Channel  : #%s\n"+
			"🎯 Request Role       : %s\n"+
			"📝 Request Reason     : %s\n"+
			"📡 Reviewers Channel  : #%s\n"+
			"⏳ Request Expiry     : %s (UTC)\n"+
			"⏰ Role Expiry        : %s (UTC)\n"+
			"\n"+
			"📅 Created At         : %s (UTC)",
		a.slackUser.RealName,
		a.accessRequest.InputChannelName,
		a.accessRequest.Role,
		a.accessRequest.Reason,
		a.accessRequest.ReviewChannelName,
		a.accessRequest.Expires.Format(util.SecondTimeFormat),
		a.accessRequest.AccessDuration.Format(util.SecondTimeFormat),
		a.accessRequest.CreateDate.Format(util.SecondTimeFormat),
	)

	section := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("```\n%s\n```", text), false, false),
		nil, nil,
	)
	return section
}

func (a *accessReviewBuilder) BuildRadioBlock() *slack.InputBlock {
	radioOptions := []*slack.OptionBlockObject{
		slack.NewOptionBlockObject("allow", slack.NewTextBlockObject("plain_text", "✅ Allow", false, false), nil),
		slack.NewOptionBlockObject("deny", slack.NewTextBlockObject("plain_text", "⛔ Deny", false, false), nil),
	}
	radioElement := slack.NewRadioButtonsBlockElement("review_decision", radioOptions...)
	radioBlock := slack.NewInputBlock(
		"review_radio",
		slack.NewTextBlockObject("plain_text", "Choose Action", false, false),
		nil,
		radioElement,
	)
	return radioBlock
}

func (a *accessReviewBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "Enter the reason", false, false),
		"review_reason",
	)
	reasonBlock := slack.NewInputBlock(
		"reason_input",
		slack.NewTextBlockObject("plain_text", "Review Reason", false, false),
		nil,
		reasonElement,
	)
	return reasonBlock
}

func (a *accessReviewBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &viewsubmission.AccessReviewModalPrivateMetadataPayload{
		ChannelID:         a.channelID,
		AccessRequestName: a.accessRequest.Name,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
