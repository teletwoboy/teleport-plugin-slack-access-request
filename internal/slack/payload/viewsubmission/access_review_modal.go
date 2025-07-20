package viewsubmission

type AccessRequestReviewModalPayload struct {
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

	Email string
}

type AccessReviewModalPrivateMetadataPayload struct {
	ChannelID         string `json:"channel_id"`
	AccessRequestName string `json:"access_request_name"`
}
