// Package slack provides business logic for interacting with Slack,
// such as sending and receiving messages, and retrieving channel lists.
package slack

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/config"

	"github.com/slack-go/slack"
)

/*
Client : 외부 API를 직접 사용하지 않고, API 호출을 감싼 Adapter 계층
바로 Service 에서 API 인터페이스 필드를 가질수도 있지만,
중간 다리 역할을 둔 이유는 외부 호출에 대한 책임을 Service 계층으로부터 분리해서
Service 로직은 내부 비즈니스 로직에만 집중하게 하여 변동성을 줄이고,
외부 연결 책임지는 Client 를 통해 외부 호출 로직을 한 곳에서 관리(캡슐화)하기 위함.
아래 GetUsers 와 같은 외부 연동 함수들을 Service 계층에서 구현한다면 복잡해짐
*/
type Client struct {
	api API
}

// Init initializes and returns a new slack client using token
func Init() (*Client, error) {
	token := config.Cfg.Slack.Token
	api := slack.New(token)

	_, err := api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("failed to perform slack auth test: %w", err)
	}

	return &Client{api: api}, nil
}

func (c *Client) GetUsers(options ...slack.GetUsersOption) ([]slack.User, error) {
	return c.api.GetUsers(options...)
}

func (c *Client) GetTeamInfo() (*slack.TeamInfo, error) {
	return c.api.GetTeamInfo()
}

func (c *Client) GetConversations(params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	return c.api.GetConversations(params)
}

func (c *Client) OpenView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error) {
	return c.api.OpenView(triggerID, view)
}
