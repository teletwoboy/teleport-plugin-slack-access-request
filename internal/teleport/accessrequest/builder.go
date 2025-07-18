package accessrequest

import (
	"github.com/google/uuid"
	"github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
	"teleport-plugin-slack-access-request/internal/slack/modal"
)

type Builder interface {
	Build() types.AccessRequest
}

type V3Builder struct {
	Payload *modal.AccessRequestViewSubmissionPayload
}

func NewV3Builder(p *modal.AccessRequestViewSubmissionPayload) *V3Builder {
	return &V3Builder{
		Payload: p,
	}
}

func (v *V3Builder) Build() types.AccessRequest {
	email := v.Payload.Email
	roles := v.Payload.View.State.Values.RoleBlock.RoleSelect.SelectedOption.Text.Text
	reason := v.Payload.View.State.Values.ReasonBlock.ReasonInput.Value
	return &types.AccessRequestV3{
		Kind:    types.KindAccessRequest,
		Version: types.V2,
		Metadata: types.Metadata{
			Name:      uuid.NewString(),
			Namespace: defaults.Namespace,
		},
		Spec: types.AccessRequestSpecV3{
			User:          email,
			Roles:         []string{roles},
			RequestReason: reason,
		},
	}
}
