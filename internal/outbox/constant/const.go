package constant

import "time"

const (
	OutboxPollInterval = 2 * time.Second
	MaxRetries         = 1

	Pending = "pending"
)
