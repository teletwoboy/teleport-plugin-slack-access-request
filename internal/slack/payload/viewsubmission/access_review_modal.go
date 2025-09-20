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

package viewsubmission

import (
	"encoding/json"
	"fmt"
)

type AccessReviewModalPayload struct {
	Type string `json:"type"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	View struct {
		PrivateMetadata string `json:"private_metadata"`
		CallbackID      string `json:"callback_id"`

		State struct {
			Values struct {
				ReasonInput struct {
					ReviewReason struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"review_reason"`
				} `json:"reason_input"`

				ReviewRadio struct {
					ReviewDecision struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`
						} `json:"selected_option,omitempty"`
					} `json:"review_decision"`
				} `json:"review_radio"`
			} `json:"values"`
		} `json:"state"`
	} `json:"view"`

	Name string
}

type AccessReviewModalPrivateMetadataPayload struct {
	ChannelID         string `json:"channel_id"`
	AccessRequestName string `json:"access_request_name"`
	MessageTs         string `json:"message_ts"`
}

type AccessReviewModal struct {
	AccessRequestName string
	Decision          string
	Reason            string
	ReviewerID        string
	ReviewerName      string
	ReviewerChannelID string
	MessageTs         string
}

func ParseAccessReviewModal(payloadStr string) (*AccessReviewModal, error) {
	var payload *AccessReviewModalPayload
	err := json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessReviewModalPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}

	return &AccessReviewModal{
		AccessRequestName: privateMetadata.AccessRequestName,
		Decision:          payload.View.State.Values.ReviewRadio.ReviewDecision.SelectedOption.Value,
		Reason:            payload.View.State.Values.ReasonInput.ReviewReason.Value,
		ReviewerID:        payload.User.ID,
		ReviewerName:      payload.User.Name,
		ReviewerChannelID: privateMetadata.ChannelID,
		MessageTs:         privateMetadata.MessageTs,
	}, nil
}
