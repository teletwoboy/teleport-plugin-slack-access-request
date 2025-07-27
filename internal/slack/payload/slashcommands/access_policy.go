package slashcommands

import (
	"fmt"
	"net/http"
)

type AccessPolicy struct {
	ChannelID   string
	ChannelName string
	Command     string
	TriggerID   string
	UserID      string
	UserName    string
}

func ParseAccessPolicy(r *http.Request, w http.ResponseWriter) (*AccessPolicy, error) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, fmt.Errorf("failed to parse access policy form: %w", err)
	}
	return &AccessPolicy{
		ChannelID:   r.Form.Get("channel_id"),
		ChannelName: r.Form.Get("channel_name"),
		Command:     r.Form.Get("command"),
		TriggerID:   r.Form.Get("trigger_id"),
		UserID:      r.Form.Get("user_id"),
		UserName:    r.Form.Get("user_name"),
	}, nil
}
