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

package verifier

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/trace"
)

type Teleport struct {
	srv teleport.Service
}

func NewTeleport(srv teleport.Service) *Teleport {
	return &Teleport{
		srv: srv,
	}
}

func (t *Teleport) VerifyUserExistsByUsername(ctx context.Context, username string) error {
	exists, err := t.srv.ExistsUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user <%s> not found in DB ", username)
	}
	return nil
}

func (t *Teleport) VerifyUserNotExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := t.srv.ExistsUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}
	return true, nil
}

func (t *Teleport) VerifyAccessRequestFromCluster(ctx context.Context, name string) error {
	builder := accessrequest.NewFilterBuilder(name)
	accessRequests, err := t.srv.FetchAccessRequests(ctx, builder)
	if err != nil {
		return err
	}

	var accessRequest types.AccessRequest
	for _, a := range accessRequests {
		copied := a.Copy()
		if copied.GetName() == name {
			accessRequest = copied
		}
	}

	if accessRequest == nil {
		return fmt.Errorf("access request <%s> not found in cluster", name)
	}

	if accessRequest.GetState() == types.RequestState_APPROVED || accessRequest.GetState() == types.RequestState_DENIED {
		return fmt.Errorf("access request <%s> is already reviewed", name)
	}
	return nil
}

func (t *Teleport) VerifyAccessRequestFromDB(ctx context.Context, name string) error {
	accessRequest, err := t.srv.GetAccessRequestByName(ctx, name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("access request <%s> not found in DB", name)
		default:
			return err
		}
	}

	if accessRequest.State != "PENDING" {
		return fmt.Errorf("access request <%s> is already reviewed", name)
	}
	return nil
}

func (t *Teleport) VerifyUserLoginStateExists(ctx context.Context, name string) (bool, error) {
	_, err := t.srv.GetUserLoginState(ctx, name)
	if err != nil {
		if trace.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
