package accessrequest

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

func BuildFirstStepSectionBlock(realName string) *slack.SectionBlock {
	section := fmt.Sprintf(util.ARequestFirstStepSection, realName)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildSecondStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.ARequestSecondStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}

func BuildThirdStepSectionBlock() *slack.SectionBlock {
	section := fmt.Sprintf(util.ARequestThirdStepSection)
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(util.Markdown, section, false, false),
		nil, nil,
	)
}
