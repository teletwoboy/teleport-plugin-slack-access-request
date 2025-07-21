package accessrequest

import (
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"

	"github.com/google/uuid"
	"github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
)

type CreateBuilder interface {
	Build() types.AccessRequest
}

type v3Builder struct {
	Payload *viewsubmission.AccessRequestModal
}

func NewV3Builder(p *viewsubmission.AccessRequestModal) CreateBuilder {
	return &v3Builder{
		Payload: p,
	}
}

func (v *v3Builder) Build() types.AccessRequest {
	email := v.Payload.Email
	roles := v.Payload.Role
	reason := v.Payload.Reason
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
