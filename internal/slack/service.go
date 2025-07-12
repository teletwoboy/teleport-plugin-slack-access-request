package slack

import "github.com/slack-go/slack"

/*
Service provides Slack-related business logic.

- 클라이언트 객체를 구조체로 감싸는 이유
현재 Service 계층을 사용하는 상위 계층에서 테스트를 수행할때,
상위 계층에서 하위 계층에 대한 직접적인 호출이 부담스러울 때 Service를 인터페이스로 만들면 되는데,
이는 외부(DB, API 등)와 연동된 로직이 있는 경우 그렇다.
내부 비즈니스 로직을 구현하는 경우엔,
부담이 없으니 구조체로 사용하여 부르고 외부와 연동된 객체는 mocking하면 되기에 struct로 한다.
*/
type Service struct {
	client *Client
}

func (s *Service) GetUsers() ([]User, error) {
	rawUsers, err := s.client.GetUsers()
	if err != nil {
		return nil, err
	}

	activeUsers := filterActiveUsers(rawUsers)
	return convertToSlackUsers(activeUsers), nil
}

func (s *Service) GetTeamInfo() (TeamInfo, error) {
	rawTeamInfo, err := s.client.GetTeamInfo()
	if err != nil {
		return TeamInfo{}, err
	}

	return TeamInfo{
		ID:   rawTeamInfo.ID,
		Name: rawTeamInfo.Name,
	}, nil
}

func filterActiveUsers(users []slack.User) []slack.User {
	var result []slack.User
	for _, user := range users {
		if user.Deleted {
			continue
		}
		result = append(result, user)
	}
	return result
}

func convertToSlackUsers(users []slack.User) []User {
	var result []User
	for _, user := range users {
		result = append(result, User{
			ID:       user.ID,
			Name:     user.Name,
			RealName: user.RealName,
			Email:    user.Profile.Email,
			Deleted:  user.Deleted,
		})
	}
	return result
}
