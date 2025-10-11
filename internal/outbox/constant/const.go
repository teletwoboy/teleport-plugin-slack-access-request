/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
