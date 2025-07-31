package accesspolicy

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

func BuildFirstStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APFirstStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSecondStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APSecondStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildThirdStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APThirdStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APFourthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepCautionSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APFourthStepCautionSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepFirstSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APFourthStepFirstSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSecondSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APFourthStepSecondSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFifthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APFifthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSixthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.APSixthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}
