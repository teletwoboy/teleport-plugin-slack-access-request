package constant

import "time"

const (
	OutboxChannel        = "outbox_notification"
	ListenMaxConcurrency = 16

	BackupInterval             = 10 * time.Second
	BackupMaxConcurrency       = 16
	BackupPullSize       int32 = 10

	AlertingDeadInterval             = 10 * time.Second
	AlertingDeadMaxConcurrency       = 16
	AlertingDeadPullSize       int32 = 10

	ProcessingTimeout = 3 * time.Second
	DeadTimeout       = 3 * time.Second

	MaxRetries  = 1
	NextTrySecs = 5.0

	Pending    = "pending"
	Processing = "processing"
	Failed     = "failed"
	Done       = "done"
	Dead       = "dead"
)
