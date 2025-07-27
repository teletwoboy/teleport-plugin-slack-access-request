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

type DeleteUserHandler struct {
	DB       *database.DB
	Clients  *container.Clients
	Services *container.Services
}

func NewDeleteUserHandler(db *database.DB, c *container.Clients, s *container.Services) *DeleteUserHandler {
	return &DeleteUserHandler{
		DB:       db,
		Clients:  c,
		Services: s,
	}
}

func (c *DeleteUserHandler) Handle(ctx context.Context, resource *types.ResourceHeader) {
	// 1. 값 준비
	username := resource.GetName()

	// 2. 검증
	slackVerifier := verifier.NewSlack(c.Services.Slack)
	teleportVerifier := verifier.NewTeleport(c.Services.Teleport)

	//    1. Teleport User가 데이터베이스에 존재하는가?
	check, err := teleportVerifier.VerifyUserNotExistsByUsername(ctx, username)
	if err != nil {
		slog.Error("failed to verify existing user", "err", err)
		return
	}
	if check {
		slog.Info("user not exists in DB", "username", username)
		return
	}

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

	// 5. Teleport User 삭제하기
	teleportUser := teleportmodels.NewUser(username)
	DeletedTeleportUser, err := txServices.Teleport.DeleteUser(ctx, teleportUser)
	if err != nil {
		slog.Error("failed to create teleport user", "err", err)
		return
	}

	// 6. Slack User 삭제하기
	DeletedSlackUser, err := txServices.Slack.DeleteUser(ctx, slackUser)
	if err != nil {
		slog.Error("failed to create slack user", "err", err)
		return
	}

	// 7. User 삭제하기
	user := usermodels.NewUser(DeletedSlackUser, DeletedTeleportUser)
	DeletedUser, err := txServices.User.DeleteUser(ctx, user)
	if err != nil {
		slog.Error("failed to create user", "err", DeletedUser)
		return
	}

	// 8. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "err", err)
		return
	}
	committed = true
	slog.Info("successfully deleted user", "username", username)
}
