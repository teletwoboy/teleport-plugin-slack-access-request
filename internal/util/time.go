package util

import (
	"fmt"
	"github.com/gravitational/teleport/api/types"
	"time"
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
