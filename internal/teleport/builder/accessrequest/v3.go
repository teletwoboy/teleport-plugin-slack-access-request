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
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"time"

	"github.com/google/uuid"
	"github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
)

type CreateBuilder interface {
	Build() types.AccessRequest
}

type v3Builder struct {
	Payload      *viewsubmission.AccessRequestModal
	TeleportUser *models.User
	DryRun       bool
}

func NewV3Builder(p *viewsubmission.AccessRequestModal, t *models.User) CreateBuilder {
	return &v3Builder{
		Payload:      p,
		TeleportUser: t,
		DryRun:       false,
	}
}

func NewV3BuilderChg(p *viewsubmission.AccessRequestModal, username string) CreateBuilder {
	return &v3Builder{
		Payload:      p,
		TeleportUser: models.NewUserWithUsername(username),
		DryRun:       false,
	}
}

func NewV3DryRunBuilder(r string, s, a, rT time.Time, t *models.User) CreateBuilder {
	return &v3Builder{
		Payload: &viewsubmission.AccessRequestModal{
			SelectedRole:                   r,
			SelectedStartDateTime:          s,
			SelectedAccessDurationDateTime: a,
			SelectedRequestTTLDateTime:     rT,
		},
		TeleportUser: t,
		DryRun:       true,
	}
}

func (v *v3Builder) Build() types.AccessRequest {
	req := &types.AccessRequestV3{
		Kind:    types.KindAccessRequest,
		Version: types.V2,
		Metadata: types.Metadata{
			Name:      uuid.NewString(),
			Namespace: defaults.Namespace,
		},
		Spec: types.AccessRequestSpecV3{
			User:          v.TeleportUser.Username,
			Roles:         []string{v.Payload.SelectedRole},
			RequestReason: v.Payload.Reason,
			Expires:       v.Payload.SelectedAccessDurationDateTime,
			DryRun:        v.DryRun,
		},
	}

	if !v.Payload.SelectedStartDateTime.IsZero() {
		req.Spec.AssumeStartTime = &v.Payload.SelectedStartDateTime
	}

	if !v.Payload.SelectedRequestTTLDateTime.IsZero() {
		req.Metadata.Expires = &v.Payload.SelectedRequestTTLDateTime
	}
	return req
}
