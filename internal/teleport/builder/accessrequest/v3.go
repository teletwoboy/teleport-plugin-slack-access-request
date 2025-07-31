package accessrequest

import (
	"github.com/google/uuid"
	"github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/models"
)

type CreateBuilder interface {
	Build() types.AccessRequest
}

type v3Builder struct {
	Payload      *viewsubmission.AccessRequestModal
	TeleportUser *models.User
}

func NewV3Builder(p *viewsubmission.AccessRequestModal, t *models.User) CreateBuilder {
	return &v3Builder{
		Payload:      p,
		TeleportUser: t,
	}
}

func (v *v3Builder) Build() types.AccessRequest {
	username := v.TeleportUser.Username
	roles := v.Payload.SelectedRole
	reason := v.Payload.Reason
	return &types.AccessRequestV3{
		Kind:    types.KindAccessRequest,
		Version: types.V2,
		Metadata: types.Metadata{
			Name:      uuid.NewString(),
			Namespace: defaults.Namespace,
		},
		Spec: types.AccessRequestSpecV3{
			User:          username,
			Roles:         []string{roles},
			RequestReason: reason,
		},
	}
}
