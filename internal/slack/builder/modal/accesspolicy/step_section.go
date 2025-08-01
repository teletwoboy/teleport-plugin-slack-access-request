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
