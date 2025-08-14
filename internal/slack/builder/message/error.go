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

package message

import "github.com/slack-go/slack"

type errorBuilder struct {
	Err error
}

func NewErrorBuilder(err error) Builder {
	return &errorBuilder{
		Err: err,
	}
}

func (e *errorBuilder) Build() slack.MsgOption {
	text := ":warning: Error Occurred \n```\n [Error] " + e.Err.Error() + "```"
	return slack.MsgOptionText(text, false)
}
