package modal

type AccessRequestViewSubmissionPayload struct {
	Type  string      `json:"type"`
	User  UserPayload `json:"user"`
	View  ViewPayload `json:"view"`
	Email string
}

type UserPayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ViewPayload struct {
	PrivateMetadata string       `json:"private_metadata"`
	CallbackID      string       `json:"callback_id"`
	State           StatePayload `json:"state"`
}

type PrivateMetadataPayload struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type StatePayload struct {
	Values ValuesPayload `json:"values"`
}

type ValuesPayload struct {
	ChannelBlock ChannelBlockPayload `json:"channel_block"`
	ReasonBlock  ReasonBlockPayload  `json:"reason_block"`
	RoleBlock    RoleBlockPayload    `json:"role_block"`
}

type ChannelBlockPayload struct {
	ChannelSelect ChannelSelectPayload `json:"channel_select"`
}

type ReasonBlockPayload struct {
	ReasonInput ReasonInputPayload `json:"reason_input"`
}

type RoleBlockPayload struct {
	RoleSelect RoleSelectPayload `json:"role_select"`
}

type ChannelSelectPayload struct {
	Type           string `json:"type"`
	SelectedOption *struct {
		Value string `json:"value"`
		Text  struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"text"`
	} `json:"selected_option,omitempty"`
}

type ReasonInputPayload struct {
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

type RoleSelectPayload struct {
	Type           string `json:"type"`
	SelectedOption *struct {
		Value string `json:"value"`
		Text  struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"text"`
	} `json:"selected_option,omitempty"`
}
