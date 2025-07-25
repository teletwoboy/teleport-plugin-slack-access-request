package blockactions

import (
	"encoding/json"
	"fmt"
)

type AccessRoleModalPayload struct {
	Type      string `json:"type"`
	TriggerID string `json:"trigger_id"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	View struct {
		ID              string `json:"id"`
		PrivateMetadata string `json:"private_metadata"`
		CallbackID      string `json:"callback_id"`

		State struct {
			Values struct {
				RoleBlock struct {
					RoleSelect struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option,omitempty"`
					} `json:"role_select"`
				} `json:"role_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`

	Email string
}

type AccessRoleModalPrivateMetadataPayload struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type AccessRoleModal struct {
	Email                string
	RequesterChannelID   string
	RequesterChannelName string
	RequesterID          string
	RequesterName        string
	Role                 string
	TriggerID            string
	ViewHash             string
	ViewID               string
}

func ParseAccessRoleModal(payloadStr string) (*AccessRoleModal, error) {
	var payload AccessRoleModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessRoleModalPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}

	return &AccessRoleModal{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		Role:                 payload.View.State.Values.RoleBlock.RoleSelect.SelectedOption.Value,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
	}, nil
}
