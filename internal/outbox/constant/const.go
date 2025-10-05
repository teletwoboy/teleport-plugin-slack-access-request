package constant

import "time"

const (
	OutboxChannel = "outbox_notification"
	MaxConcurrent = 10

	BackupInterval             = 10 * time.Second
	BackupMaxConcurrency       = 16
	BackupPullSize       int32 = 10

	AlertingDeadInterval             = 10 * time.Second
	AlertingDeadMaxConcurrency       = 16
	AlertingDeadPullSize       int32 = 10

	MaxRetries  = 1
	NextTrySecs = 5.0

	Pending    = "pending"
	Processing = "processing"
	Failed     = "failed"
	Done       = "done"
	Dead       = "dead"
)
