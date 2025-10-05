package constant

import "time"

const (
	OutboxChannel = "outbox_notification"

	ListenMaxConcurrency = 5
	ListenQueueSize      = 50

	BackupInterval             = 10 * time.Second
	BackupMaxConcurrency       = 5
	BackupQueueSize            = 50
	BackupPullSize       int32 = 5

	AlertingDeadInterval             = 10 * time.Second
	AlertingDeadMaxConcurrency       = 5
	AlteringDeadQueueSize            = 50
	AlertingDeadPullSize       int32 = 5

	ClaimTimeout      = 3 * time.Second
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
