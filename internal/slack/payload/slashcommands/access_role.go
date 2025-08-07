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

package slashcommands

import (
	"fmt"
	"net/http"
)

type AccessRole struct {
	ChannelID   string
	ChannelName string
	Command     string
	TriggerID   string
	UserID      string
	UserName    string
}

func ParseAccessRole(r *http.Request, w http.ResponseWriter) (*AccessRole, error) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return nil, fmt.Errorf("failed to parse access role payload: %w", err)
	}
	return &AccessRole{
		ChannelID:   r.FormValue("channel_id"),
		ChannelName: r.FormValue("channel_name"),
		Command:     r.FormValue("command"),
		TriggerID:   r.FormValue("trigger_id"),
		UserID:      r.FormValue("user_id"),
		UserName:    r.FormValue("user_name"),
	}, nil
}
