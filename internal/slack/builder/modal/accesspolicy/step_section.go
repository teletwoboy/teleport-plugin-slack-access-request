package accesspolicy

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"

	"github.com/slack-go/slack"
)

func BuildFirstStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APFirstStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSecondStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APSecondStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildThirdStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APThirdStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APFourthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepFirstSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APFourthStepFirstSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSecondSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APFourthStepSecondSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepThirdSubSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APFourthStepThirdSubSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFifthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APFifthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSixthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(modal.APSixthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(modal.Markdown, section, false, false),
		nil, nil,
	)
}
