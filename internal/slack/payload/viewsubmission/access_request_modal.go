package viewsubmission

import (
	"encoding/json"
	"fmt"
)

type AccessRequestModalPayload struct {
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
				ChannelBlock struct {
					ChannelSelect struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option,omitempty"`
					} `json:"channel_select"`
				} `json:"channel_block"`

				ReasonBlock struct {
					ReasonInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"reason_input"`
				} `json:"reason_block"`
			} `json:"values"`
		} `json:"state"`
	} `json:"view"`

	Email string
}

type AccessRequestModalPrivateMetadataPayload struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Role        string `json:"role"`
}

type AccessRequestModal struct {
	Username             string
	Reason               string
	RequesterChannelID   string
	RequesterChannelName string
	RequesterID          string
	RequesterName        string
	ReviewersChannelID   string
	ReviewersChannelName string
	Role                 string
}

func ParseAccessRequestModal(payloadStr string) (*AccessRequestModal, error) {
	var payload AccessRequestModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessRequestModalPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}

	return &AccessRequestModal{
		Reason:               payload.View.State.Values.ReasonBlock.ReasonInput.Value,
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		ReviewersChannelID:   payload.View.State.Values.ChannelBlock.ChannelSelect.SelectedOption.Value,
		ReviewersChannelName: payload.View.State.Values.ChannelBlock.ChannelSelect.SelectedOption.Text.Text,
		Role:                 privateMetadata.Role,
	}, nil
}
