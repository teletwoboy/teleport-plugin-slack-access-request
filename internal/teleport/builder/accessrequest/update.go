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
	if u.decision == "allow" {
		return types.RequestState_APPROVED
	}
	return types.RequestState_DENIED
}
