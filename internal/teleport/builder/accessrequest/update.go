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

package accessrequest

import (
	"time"

	"github.com/gravitational/teleport/api/types"
)

type UpdateBuilder interface {
	Build() types.AccessRequestUpdate
}

type updateBuilder struct {
	accessRequestName string
	decision          string
	reason            string
}

func NewUpdateBuilder(a, d, r string) UpdateBuilder {
	return &updateBuilder{
		accessRequestName: a,
		decision:          d,
		reason:            r,
	}
}

func (u *updateBuilder) Build() types.AccessRequestUpdate {
	requestState := u.BuildRequestState()
	now := time.Now().Add(time.Second)
	return types.AccessRequestUpdate{
		RequestID:       u.accessRequestName,
		State:           requestState,
		Reason:          u.reason,
		AssumeStartTime: &now,
	}
}

func (u *updateBuilder) BuildRequestState() types.RequestState {
	if u.decision == types.RequestState_APPROVED.String() {
		return types.RequestState_APPROVED
	}
	return types.RequestState_DENIED
}
