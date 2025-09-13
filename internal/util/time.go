package util

import (
	"fmt"
	"time"

	"github.com/gravitational/teleport/api/types"
)

func ParseDateTimeInLocation(dateStr, timeStr, timezone string) (time.Time, error) {
	dateTime := dateStr + " " + timeStr

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone format: %w", err)
	}

	t, err := time.ParseInLocation(MinuteTimeFormat, dateTime, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}
	return t, nil
}

func ParseTTLInLocation(ar types.AccessRequest, timezone string) (time.Time, error) {
	t := ar.GetAccessExpiry()
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone format: %w", err)
	}
	t = t.In(loc)
	return t, nil
}

func ParseInLocation(t time.Time, timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return t
	}
	return t.In(loc)
}
