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
	err := teleportVerifier.VerifyUserExistsByUsername(ctx, username)
	if err != nil {
		slog.Error("failed to verify existing user", "err", err)
		return
	}

	//    2. Username을 갖는 Slack User가 데이터베이스에 존재하는지 확인하고 해당하는 User의 정보를 가져옴
	teleportUserInfo, err := c.Services.Teleport.GetUserByUsername(ctx, username)
	if err != nil {
		slog.Error("failed to get teleport user", "err", err)
	}
	userInfo, err := c.Services.User.GetUserByTeleportUserID(ctx, teleportUserInfo.TeleportUserID)
	if err != nil {
		slog.Error("failed to get user", "err", err)
	}
	slackUser, err := slackVerifier.VerifyUserExistsBySlackID(ctx, userInfo.SlackUser.SlackUserID)
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

	// # user state가 존재할 경우 삭제하기
	exist, err := teleportVerifier.VerifyUserState(ctx, username)
	if err != nil {
		slog.Error("not found teleport user state", "err", err)
	}
	if exist {
		err := txServices.Teleport.DeleteUserLoginState(ctx, username)
		if err != nil {
			slog.Error("failed to delete user state", "err", err)
			return
		}
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
