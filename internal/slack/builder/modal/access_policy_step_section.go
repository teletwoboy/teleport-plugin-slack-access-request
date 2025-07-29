package modal

import (
	"fmt"
	"github.com/slack-go/slack"
)

func BuildFirstStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(firstStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSecondStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(secondStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(Markdown, section, false, false),
		nil, nil,
	)
}

func BuildThirdStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(thirdStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFourthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(fourthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(Markdown, section, false, false),
		nil, nil,
	)
}

func BuildFifthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(fifthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSixthStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(sixthStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(Markdown, section, false, false),
		nil, nil,
	)
}
