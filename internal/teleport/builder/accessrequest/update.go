package accessrequest

import (
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
	return types.AccessRequestUpdate{
		RequestID: u.accessRequestName,
		State:     requestState,
		Reason:    u.reason,
	}
}

func (s *updateBuilder) BuildRequestState() types.RequestState {
	value := s.decision

	if value == "allow" {
		return types.RequestState_APPROVED
	}
	return types.RequestState_APPROVED
}
