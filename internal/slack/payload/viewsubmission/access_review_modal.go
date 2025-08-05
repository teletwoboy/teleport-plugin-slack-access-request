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
}

type AccessReviewModal struct {
	AccessRequestName string
	Decision          string
	Reason            string
	ReviewerID        string
	ReviewerName      string
	ReviewerChannelID string
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
	}, nil
}
