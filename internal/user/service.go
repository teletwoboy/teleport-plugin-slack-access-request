package user

import (
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
)

type Service struct {
}

func (s *Service) MapUsersByUsername(teleportUsers []teleport.User, slackUsers []slack.User) []User {
	var users []User
	for _, teleportUser := range teleportUsers {
		for _, slackUser := range slackUsers {
			if teleportUser.Username == slackUser.Email {
				tu := teleportUser // range 안에서 포인터 변수를 그냥 사용하는 경우,
				su := slackUser    // 모든 포인터가 마지막만 가리키는 버그 발생
				users = append(users, User{
					TeleportUser: &tu,
					SlackUser:    &su,
				})
			}
		}
	}

	return users
}
