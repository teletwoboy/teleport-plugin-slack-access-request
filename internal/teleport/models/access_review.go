package models

import "time"

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
