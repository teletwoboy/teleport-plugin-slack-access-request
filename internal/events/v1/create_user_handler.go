package v1

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/database"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	usermodels "teleport-plugin-slack-access-request/internal/user/models"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"

	"github.com/gravitational/teleport/api/types"
)

type CreateUserHandler struct {
	DB       *database.DB
	Clients  *container.Clients
	Services *container.Services
}

func NewCreateUserHandler(db *database.DB, c *container.Clients, s *container.Services) *CreateUserHandler {
	return &CreateUserHandler{
		DB:       db,
		Clients:  c,
		Services: s,
	}
}

func (c *CreateUserHandler) Handle(ctx context.Context, resource *types.UserV2) {
	// 1. 값 준비
	username := resource.GetName()

	// 2. 검증
	slackVerifier := verifier.NewSlack(c.Services.Slack)
	teleportVerifier := verifier.NewTeleport(c.Services.Teleport)
	//    1. Teleport User가 데이터베이스에 존재하지 않는가?
	check, err := teleportVerifier.VerifyUserNotExistsByUsername(ctx, username)
	if err != nil {
		slog.Error("failed to verify existing user", "err", err)
		return
	}
	if !check {
		slog.Info("user already exists in DB", "username", username)
		return
	}

	//    2. Username을 갖는 Slack User가 Slack에 존재하는가?
	slackUser, err := slackVerifier.VerifyUserExistsByUsernameFromClient(username)
	if err != nil {
		slog.Error("failed to verify existing user", "err", err)
		return
	}

	// 3. 트랜잭션 시작하기
	tx, err := c.DB.Conn.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		return
	}
	committed := false
	defer func(tx *sql.Tx) {
		if !committed {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				slog.Error("failed to rollback transaction", "err", err)
			}
		}
	}(tx)

	// 4. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := c.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(c.Clients, txRepos)

	// 5. Teleport User 저장하기
	teleportUser := teleportmodels.NewUser(username)
	createdTeleportUser, err := txServices.Teleport.CreateUser(ctx, teleportUser)
	if err != nil {
		slog.Error("failed to create teleport user", "err", err)
		return
	}

	// 6. Slack User 저장하기
	createdSlackUser, err := txServices.Slack.CreateUser(ctx, slackUser)
	if err != nil {
		slog.Error("failed to create slack user", "err", err)
		return
	}

	// 7. User 저장하기
	user := usermodels.NewUser(createdSlackUser, createdTeleportUser)
	createdUser, err := txServices.User.CreateUser(ctx, user)
	if err != nil {
		slog.Error("failed to create user", "err", err)
		return
	}

	// 8. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "err", err)
		return
	}
	committed = true
	slog.Info("successfully created user", "username", createdUser.TeleportUser.Username)
}
