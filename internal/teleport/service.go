package teleport

import (
	"context"
	"fmt"

	"github.com/gravitational/teleport/api/types"
)

type Service struct {
	api API
}

func (s *Service) GetUsersWithoutSecrets(ctx context.Context) ([]User, error) {
	rawUsers, err := s.api.GetUsers(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	humanUsers := filterHumanUsers(rawUsers)
	return convertToUsers(humanUsers), nil
}

func filterHumanUsers(users []types.User) []types.User {
	var humanUsers []types.User
	for _, user := range users {
		if !user.IsBot() {
			humanUsers = append(humanUsers, user)
		}
	}
	return humanUsers
}

func convertToUsers(users []types.User) []User {
	var result []User
	for _, user := range users {
		result = append(result, User{
			Username: user.GetName(),
		})
	}
	return result
}
