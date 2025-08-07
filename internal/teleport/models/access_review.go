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

package models

import (
	"time"
)

type AccessReview struct {
	AccessReviewID  int32
	AccessRequestID int32
	ReviewerUserID  int32
	Reason          string
	Decision        string
	UseYn           bool
	CreateCode      string
	CreateDate      time.Time
	UpdateCode      string
	UpdateDate      time.Time
	DeleteCode      string
	DeleteDate      time.Time
	Version         int64
}

func NewAccessReview(arID, uID int32, reason, decision string) *AccessReview {
	return &AccessReview{
		AccessRequestID: arID,
		ReviewerUserID:  uID,
		Reason:          reason,
		Decision:        decision,
	}
}
