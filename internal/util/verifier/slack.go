package verifier

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack"
)

type Slack struct {
	Srv slack.Service
}

func NewSlack(srv slack.Service) *Slack {
	return &Slack{
		Srv: srv,
	}
}

func (s *Slack) VerifyExistsUserByID(ctx context.Context, id, name string) error {
	exists, err := s.Srv.ExistsUserByID(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user %s not found", name)
	}
	return nil
}

func (s *Slack) VerifyExistsUserInChannelByID(id, channelID string) error {
	exists, err := s.Srv.ExistsUserInChannelByID(id, channelID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user %s not found", id)
	}
	return nil
}
