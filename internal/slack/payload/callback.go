package payload

type Callback struct {
	Type    string           `json:"type"`
	View    CallbackView     `json:"view"`
	Actions []CallbackAction `json:"actions"`
}

type CallbackView struct {
	CallbackID string `json:"callback_id"`
}

type CallbackAction struct {
	ActionID string `json:"action_id,omitempty"`
}
