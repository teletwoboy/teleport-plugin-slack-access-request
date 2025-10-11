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

package res

import (
	"context"
	"log/slog"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/message"
)

func ErrorMessageToSlack(ctx context.Context, s slack.Service, channelID string, err error) {
	msg := message.NewErrorBuilder(err)
	_, _, postErr := s.PostMessageContext(ctx, channelID, msg)
	if postErr != nil {
		slog.Error("failed to post msg to slack", "err", postErr)
		return
	}
}
