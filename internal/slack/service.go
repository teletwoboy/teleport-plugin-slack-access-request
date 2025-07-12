package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

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
		return nil, fmt.Errorf("failed to get users from Slack API: %w", err)
	}

	activeUsers := filterActiveUsers(rawUsers)
	return convertToUsers(activeUsers), nil
}

func (s *Service) GetTeamInfo() (TeamInfo, error) {
	rawTeamInfo, err := s.client.GetTeamInfo()
	if err != nil {
		return TeamInfo{}, fmt.Errorf("failed to get team info from Slack API: %w", err)
	}

	return TeamInfo{
		ID:   rawTeamInfo.ID,
		Name: rawTeamInfo.Name,
	}, nil
}

func (s *Service) GetReviewersChannels() ([]ReviewersChannel, error) {
	channels, err := s.GetAllChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all channels: %w", err)
	}

	filteredReviewersChannels := filterReviewersChannels(channels)
	filteredJoinedChannels := filterJoinedChannels(filteredReviewersChannels)
	return convertToReviewersChannels(filteredJoinedChannels), nil
}

func (s *Service) GetAllChannels() ([]slack.Channel, error) {
	var channels []slack.Channel
	params := &slack.GetConversationsParameters{
		ExcludeArchived: true,
	}

	for {
		rawChannels, nextCursor, err := s.client.GetConversations(params)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversations (cursor=%s): %w", params.Cursor, err)
		}
		channels = append(channels, rawChannels...)
		if nextCursor == "" {
			break
		}
	}

	return channels, nil
}

// --- Internal Util Functions related to User ---
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

func convertToUsers(users []slack.User) []User {
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

// --- Internal Util Functions related to ReviewersChannel ---
func filterReviewersChannels(channels []slack.Channel) []slack.Channel {
	var filteredChannels []slack.Channel
	for _, channel := range channels {
		if strings.HasSuffix(channel.Name, "-reviewers") {
			filteredChannels = append(filteredChannels, channel)
		}
	}
	return filteredChannels
}

func filterJoinedChannels(channels []slack.Channel) []slack.Channel {
	var filteredChannels []slack.Channel
	for _, channel := range channels {
		if channel.IsMember {
			filteredChannels = append(filteredChannels, channel)
		}
	}
	return filteredChannels
}

func convertToReviewersChannels(channels []slack.Channel) []ReviewersChannel {
	var reviewersChannels []ReviewersChannel
	for _, channel := range channels {
		reviewersChannels = append(reviewersChannels, ReviewersChannel{
			ID:       channel.ID,
			Name:     channel.Name,
			IsMember: channel.IsMember,
		})
	}
	return reviewersChannels
}
