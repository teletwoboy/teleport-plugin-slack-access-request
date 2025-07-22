package verifier

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/teleport"
)

type Teleport struct {
	srv teleport.Service
}

func NewTeleport(srv teleport.Service) *Teleport {
	return &Teleport{
		srv: srv,
	}
}

func (t *Teleport) VerifyExistsAccessRequestsByName(ctx context.Context, name string) error {
	exists, err := t.srv.ExistsAccessRequestByName(ctx, name)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("access request <%s> not found", name)
	}
	return nil
}

func (t *Teleport) VerifyReviewedAccessRequestByName(ctx context.Context, name string) error {
	state, err := t.srv.GetAccessRequestStateByName(ctx, name)
	if err != nil {
		return err
	}

	if state != "PENDING" {
		return fmt.Errorf("access request <%s> is Already Reviewd", name)
	}
	return nil
}
