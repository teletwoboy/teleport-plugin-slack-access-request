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
