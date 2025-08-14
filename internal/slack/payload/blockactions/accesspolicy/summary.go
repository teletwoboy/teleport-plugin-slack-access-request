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

import "time"

type SummaryPrivateMetadataPayload struct {
	ChannelID           string    `json:"channel_id"`
	ChannelName         string    `json:"channel_name"`
	RealName            string    `json:"real_name"`
	TimeZone            string    `json:"time_zone"`
	SelectedChannelID   string    `json:"selected_channel_id"`
	SelectedChannelName string    `json:"selected_channel_name"`
	SelectedRole        string    `json:"selected_role"`
	SelectedRoleName    string    `json:"selected_role_name"`
	SelectedUserID      string    `json:"selected_user_id"`
	SelectedRealName    string    `json:"selected_real_name"`
	SelectedStartDate   time.Time `json:"selected_start_date"`
	SelectedEndDate     time.Time `json:"selected_end_date"`
	SelectedEffect      string    `json:"selected_effect"`
}
