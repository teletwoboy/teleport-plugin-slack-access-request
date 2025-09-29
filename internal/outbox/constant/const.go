package constant

import "time"

const (
	OutboxPollInterval = 2 * time.Second
	MaxAttempts        = 3

	Pending string = "pending"
)
