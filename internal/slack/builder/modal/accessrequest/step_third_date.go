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

package accessrequest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type thirdStepDateBuilder struct {
	payload *accessrequest.StartDateOptionSelect
	ttl     time.Time
}

func NewThirdStepDateBuilder(p *accessrequest.StartDateOptionSelect, t time.Time) modal.Builder {
	return &thirdStepDateBuilder{
		payload: p,
		ttl:     t,
	}
}

func (t *thirdStepDateBuilder) Build() (*slack.ModalViewRequest, error) {
	blocks := t.BuildBlocks()
	privateMetadata, err := t.BuildPrivateMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to build private metadata: %w", err)
	}

	modal := &slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(util.PlainText, util.ARequestTitle, false, false),
		Close:           slack.NewTextBlockObject(util.PlainText, util.Back, false, false),
		Submit:          nil,
		CallbackID:      util.ARequestCallBackID,
		Blocks:          blocks,
		PrivateMetadata: privateMetadata,
	}
	return modal, nil
}

func (t *thirdStepDateBuilder) BuildBlocks() slack.Blocks {
	var blockSet []slack.Block
	blockSet = append(blockSet, BuildThirdStepSectionBlock())
	blockSet = append(blockSet, t.BuildStartDateBlock()...)
	return slack.Blocks{
		BlockSet: blockSet,
	}
}

func (t *thirdStepDateBuilder) BuildStartDateBlock() []slack.Block {
	text := BuildStartDateSecondOptInfoText(t.payload.StartDateOptionName, t.ttl)
	startDateSecondOptInfoBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, text, false, false),
		nil,
		nil,
	)
	startDateTimeBlock := slack.NewActionBlock(
		util.ARequestStartDateTimeBlockID,
		slack.NewDatePickerBlockElement(util.ARequestStartDateBlockActionID),
	)
	return []slack.Block{
		startDateSecondOptInfoBlock,
		startDateTimeBlock,
	}
}

func (t *thirdStepDateBuilder) BuildPrivateMetadata() (string, error) {
	privateMetadata := &accessrequest.StartDateSelectPrivateMetadataPayload{
		ChannelID:                   t.payload.RequesterChannelID,
		ChannelName:                 t.payload.RequesterChannelName,
		RealName:                    t.payload.RequesterRealName,
		RequireReason:               t.payload.RequireReason,
		SelectedRole:                t.payload.SelectedRole,
		SelectedChannelID:           t.payload.SelectedChannelID,
		SelectedChannelName:         t.payload.SelectedChannelName,
		SelectedStartDateOptionID:   t.payload.StartDateOptionID,
		SelectedStartDateOptionName: t.payload.StartDateOptionName,
		TTL:                         t.ttl,
	}

	jsonBytes, err := json.Marshal(privateMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private metadata: %w", err)
	}
	return string(jsonBytes), nil
}
