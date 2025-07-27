package models

import (
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
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

func NewAccessReview(arID, uID int32, payload *viewsubmission.AccessReviewModal) *AccessReview {
	return &AccessReview{
		AccessRequestID: arID,
		ReviewerUserID:  uID,
		Reason:          payload.Reason,
		Decision:        payload.Decision,
	}
}
