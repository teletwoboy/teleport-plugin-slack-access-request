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

package accesspolicy

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

func BuildFirstStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyFirstStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSecondStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicySecondStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildThirdStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyThirdStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyFourthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepCautionSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyFourthStepCautionSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepFirstSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyFourthStepFirstSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSecondSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyFourthStepSecondSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFifthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicyFifthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSixthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APolicySixthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}
