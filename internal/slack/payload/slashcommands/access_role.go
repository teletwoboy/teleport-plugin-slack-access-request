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
