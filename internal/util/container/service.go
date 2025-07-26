package container

import (
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
)

type Services struct {
	SeedInit seedinit.Service
	Slack    slack.Service
	Teleport teleport.Service
	User     user.Service
}

func NewServices(clients *Clients, repos *Repositories) *Services {
	slackSrv := slack.NewService(clients.Slack, repos.Slack)
	teleportSrv := teleport.NewService(clients.Teleport, repos.Teleport)
	userSrv := user.NewService(repos.User, slackSrv, teleportSrv)
	seedInitSrv := seedinit.NewService(repos.SeedInit, slackSrv, teleportSrv, userSrv)
	return &Services{
		SeedInit: seedInitSrv,
		Slack:    slackSrv,
		Teleport: teleportSrv,
		User:     userSrv,
	}
}
