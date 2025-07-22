package verifier

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"

	"github.com/gravitational/teleport/api/types"
)

type Teleport struct {
	srv teleport.Service
}

func NewTeleport(srv teleport.Service) *Teleport {
	return &Teleport{
		srv: srv,
	}
}

func (t *Teleport) VerifyAccessRequestExists(ctx context.Context, name string) error {
	exists, err := t.srv.ExistsAccessRequestByName(ctx, name)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("access request <%s> not found", name)
	}
	return nil
}

func (t *Teleport) VerifyAccessRequestNotReviewedFromCluster(ctx context.Context, name string) error {
	builder := accessrequest.NewFilterBuilder(name)
	accessRequests, err := t.srv.FetchAccessRequests(ctx, builder)
	if err != nil {
		return err
	}

	var accessRequest types.AccessRequest
	for _, a := range accessRequests {
		accessRequest = a
	}

	if accessRequest.GetState() == types.RequestState_APPROVED || accessRequest.GetState() == types.RequestState_DENIED {
		return fmt.Errorf("access request <%s> is already reviewed", name)
	}
	return nil
}

func (t *Teleport) VerifyAccessRequestNotReviewedFromDB(ctx context.Context, name string) error {
	state, err := t.srv.GetAccessRequestStateByName(ctx, name)
	if err != nil {
		return err
	}

	if state != "PENDING" {
		return fmt.Errorf("access request <%s> is already reviewed", name)
	}
	return nil
}
