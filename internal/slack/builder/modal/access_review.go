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
	"teleport-plugin-slack-access-request/internal/slack/builder"
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
		Title:           slack.NewTextBlockObject(util.PlainText, "Access Review", false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, "Close", false, false),
		Submit:          slack.NewTextBlockObject(util.PlainText, util.Submit, false, false),
		CallbackID:      "access_review_modal",
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}

	return &modal, nil
}

func (a *accessReviewBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, a.BuildSectionBlock())
	blockSet = append(blockSet, a.BuildRadioBlock())
	blockSet = append(blockSet, a.BuildReasonBlock())
	return slack.Blocks{BlockSet: blockSet}
}

func (a *accessReviewBuilder) BuildSectionBlock() *slack.SectionBlock {
	text := builder.BuildAccessReviewText(a.accessRequest, a.slackUser)
	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil, nil,
	)
	return section
}

func (a *accessReviewBuilder) BuildRadioBlock() *slack.InputBlock {
	radioOptions := []*slack.OptionBlockObject{
		slack.NewOptionBlockObject("allow", slack.NewTextBlockObject(util.PlainText, "✅ Allow", false, false), nil),
		slack.NewOptionBlockObject("deny", slack.NewTextBlockObject(util.PlainText, "⛔ Deny", false, false), nil),
	}
	radioElement := slack.NewRadioButtonsBlockElement("review_decision", radioOptions...)
	radioBlock := slack.NewInputBlock(
		"review_radio",
		slack.NewTextBlockObject(util.PlainText, "Choose Action", false, false),
		nil,
		radioElement,
	)
	return radioBlock
}

func (a *accessReviewBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(util.PlainText, "Enter the reason", false, false),
		"review_reason",
	)
	reasonBlock := slack.NewInputBlock(
		"reason_input",
		slack.NewTextBlockObject(util.PlainText, "Review Reason", false, false),
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
