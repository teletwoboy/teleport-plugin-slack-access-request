package slack

import (
	"context"
	"fmt"
	"github.com/slack-go/slack"
	"strings"
	"teleport-plugin-slack-access-request/internal/slack/message"
	"teleport-plugin-slack-access-request/internal/slack/modal"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/types"
)

type Service interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ExistsUserByID(ctx context.Context, id string) (bool, error)
	FetchUsers() ([]models.User, error)
	FetchTeamInfo() (*types.TeamInfo, error)
	FetchReviewersChannels() ([]types.ReviewersChannel, error)
	FetchAllChannels() ([]slack.Channel, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	OpenModal(triggerID string, builder modal.Builder) error
	PostMessage(channelID string, builder message.Builder) (string, string, error)
}

/*
API is interface for Slack

Go 에선 런타임 시 JVM 위에서 리플렉션과 프록시가 가능한 Java 와 다르게
컴파일 시 모든 구조체, 메서드, 타입이 고정되어 실행 중에 동작 변경이 불가능함
때문에 외부(DB, API 등) 의존성에 대한 mocking을 위해선,
여러 구현체를 가질 수 있는 Interface를 통해서만 가짜 구현체를 만들어 직접적인 외부 호출을 하지 않아도 되기에
인터페이스로 정의함.

- client.go 에서 service.go 로 옮기는 이유 : Go는 사용하는 측에서 Interface를 정의함
*/
type API interface {
	GetUsers(options ...slack.GetUsersOption) ([]slack.User, error)
	GetTeamInfo() (*slack.TeamInfo, error)
	GetConversations(params *slack.GetConversationsParameters) (channels []slack.Channel, nextCursor string, err error)
	OpenView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error)
	PostMessage(channel string, options ...slack.MsgOption) (string, string, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ExistsUserByID(ctx context.Context, id string) (bool, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

/*
Service provides Slack-related business logic.

- 클라이언트 객체를 구조체로 감싸는 이유
현재 Service 계층을 사용하는 상위 계층에서 테스트를 수행할때,
상위 계층에서 하위 계층에 대한 직접적인 호출이 부담스러울 때 Service를 인터페이스로 만들면 되는데,
이는 외부(DB, API 등)와 연동된 로직이 있는 경우 그렇다.
내부 비즈니스 로직을 구현하는 경우엔,
부담이 없으니 구조체로 사용하여 부르고 외부와 연동된 객체는 mocking하면 되기에 struct로 한다.

- 인터페이스 객체로 변경한 이유
서비스가 구조체에 의존하고 있다는 것은 서비스가 특정(Client) 어댑터에 의존하고 있다는 것임.
즉, 다르게 생긴 API 를 구현한 구조체를 사용하지 못한다는 뜻.
이를 구조화 하자면 [ Service ] --depends on--> [ *Client ] --implement--> [ API ] 이다.
기존 테스트 코드 또한 구조가
[Service] --depends on--> [*Client] --implement--> [API] <--implement-- [*mockAPI] 이며
[Service] --depends on--> [*Client] <-이 자리에 [mockAPI] 대입하였고,
이는 결과적으로 테스트 시, 무조건 *Client 와 동일한 구조체를 만들어주어야만 했던 구조였음.
(*mockAPI는 *Client 와 동일했기 때문에 client 자리에 대입이 가능했던 것. *mockAPI에 다른 Method 하나가 더 있었다면 불가능 했음.)

이 문제점은 Teleport Service 를 테스트하면서 알게됨
Slack의 반환값들은 struct 이기에 리터럴로 쉽게 필요한 mock 데이터만 구성이 가능했지만,
Teleport GetUser()의 반환값인 models.User 는 Interface 였고,
테스트 시 models.User 인터페이스를 만족하는 mock 구조체를 만들어여 했음.
models.User 의 필요 메서드는 2개였지만, 그 안의 수십개의 필수 메서드들을 구현해야하는 상황이 오게됨.
이를 위해 gomock 라이브러리를 사용하여 models.User 인터페이스를 자동으로 mock 객체로 생성하였음.
이후,
mockAPI에 [users []models.User] 을 넣어주며 특정 테스트마다의 models.User 구성을 자유롭게 하려했음.
이때 문제가 발생함!!
Client 구조체에는 [users []models.User] 라는 필드는 없음. -> 불일치 발생 -> Service 의 client 필드에 대입 불가!
물론, 직접 GetUsers 메서드에서 NewMockUser()를 하드코딩하여 이를 회피할 수 있지만,
모든 테스트에서 GetUsers 반환값을 동일하게 사용하게되어 테스트 자유도 매우 하락함.

따라서
1. 테스트 시 Client 에 대한 결합도를 줄이기 위해서
2. gomock을 활용해 테스트별 mock 데이터를 유연하게 구성하기 위해서
2. Service 가 특정 어댑터에 대한 의존도를 줄이기 위해서
Service 가 API 를 의존하도록 변경함.
*/
type service struct {
	api  API
	repo Repository
}

func NewService(api API, repo Repository) Service {
	return &service{api: api, repo: repo}
}

func (s *service) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create slack user: %w", err)
	}
	return createdUser, nil
}

func (s *service) ExistsUserByID(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.ExistsUserByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (s *service) FetchUsers() ([]models.User, error) {
	rawUsers, err := s.api.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users from Slack API: %w", err)
	}

	activeUsers := filterActiveUsers(rawUsers)
	return convertToUsers(activeUsers), nil
}

func (s *service) FetchTeamInfo() (*types.TeamInfo, error) {
	rawTeamInfo, err := s.api.GetTeamInfo()
	if err != nil {
		return &types.TeamInfo{}, fmt.Errorf("failed to get team info from Slack API: %w", err)
	}
	return &types.TeamInfo{
		ID:   rawTeamInfo.ID,
		Name: rawTeamInfo.Name,
	}, nil
}

func (s *service) FetchReviewersChannels() ([]types.ReviewersChannel, error) {
	channels, err := s.FetchAllChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all channels: %w", err)
	}

	reviewersChannels := filterReviewersChannels(channels)
	joinedChannels := filterJoinedChannels(reviewersChannels)
	return convertToReviewersChannels(joinedChannels), nil
}

func (s *service) FetchAllChannels() ([]slack.Channel, error) {
	var channels []slack.Channel
	params := &slack.GetConversationsParameters{
		ExcludeArchived: true,
		Types:           []string{"public_channel", "private_channel"},
	}

	for {
		rawChannels, nextCursor, err := s.api.GetConversations(params)
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

func (s *service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID (%s): %w", id, err)
	}
	return user, nil
}

func (s *service) OpenModal(triggerID string, builder modal.Builder) error {
	builtModal, err := builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build modal: %w", err)
	}

	_, err = s.api.OpenView(triggerID, *builtModal)
	return err
}

func (s *service) PostMessage(channelID string, builder message.Builder) (string, string, error) {
	msgOption := builder.Build()
	return s.api.PostMessage(channelID, msgOption)
}

// --- Internal Util Functions related to User ---
func filterActiveUsers(users []slack.User) []slack.User {
	var activeUsers []slack.User
	for _, user := range users {
		copiedUser := user
		if !copiedUser.Deleted {
			activeUsers = append(activeUsers, copiedUser)
		}
	}
	return activeUsers
}

func convertToUsers(users []slack.User) []models.User {
	var result []models.User
	for _, user := range users {
		copiedUser := user
		result = append(result, models.User{
			ID:       copiedUser.ID,
			Name:     copiedUser.Name,
			RealName: copiedUser.RealName,
			Email:    copiedUser.Profile.Email,
			Deleted:  copiedUser.Deleted,
		})
	}
	return result
}

// --- Internal Util Functions related to ReviewersChannel ---
func filterReviewersChannels(channels []slack.Channel) []slack.Channel {
	var reviewersChannels []slack.Channel
	for _, channel := range channels {
		copiedChannel := channel
		if strings.HasSuffix(copiedChannel.Name, "-reviewers") {
			reviewersChannels = append(reviewersChannels, copiedChannel)
		}
	}
	return reviewersChannels
}

func filterJoinedChannels(channels []slack.Channel) []slack.Channel {
	var joinedChannels []slack.Channel
	for _, channel := range channels {
		copiedChannel := channel
		if copiedChannel.IsMember {
			joinedChannels = append(joinedChannels, copiedChannel)
		}
	}
	return joinedChannels
}

func convertToReviewersChannels(channels []slack.Channel) []types.ReviewersChannel {
	var result []types.ReviewersChannel
	for _, channel := range channels {
		copiedChannel := channel
		result = append(result, types.ReviewersChannel{
			ID:       copiedChannel.ID,
			Name:     copiedChannel.Name,
			IsMember: copiedChannel.IsMember,
		})
	}
	return result
}
