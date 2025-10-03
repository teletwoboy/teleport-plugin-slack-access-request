package constant

import "time"

const (
	PollInterval       = 1 * time.Second
	PollSize     int32 = 1
	MaxRetries         = 1
	NextTrySecs        = 2.0

	Pending    = "pending"
	Processing = "processing"
	Failed     = "failed"
	Done       = "done"
	Dead       = "dead"
)
