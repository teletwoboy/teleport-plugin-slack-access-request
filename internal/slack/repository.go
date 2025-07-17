package slack

import (
	"context"
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
)

/*
PostgresRepository implements Repository interface

- 왜 client.go의 Client 는 동일하게 외부와 통신하는 struct 인데 interface를 가지고 있고, PostgresRepository 는 구현체를 가지고 있는가?
일단 비유를 통해 설명하자면,
Client 는 통역사이며, PostgresRepository 는 사서 선생님.

1. Client
나 (한글로 GetUser만 아는 Service) --> 통역사 (영어로 GetUser 번역하는 Client) --> 외국인(영어로 GetUser 수행하는 *slack.Client)
즉, Client 는 동작은 알필요 없는 단순히 통역만 할줄 아는 사람.
이런 Client 가 Interface를 가지는 이유는
1. 특정 외국인(*slack.Client)에 의지하지 않고, 영어 및 GetUser 수행 가능한 어떤 외국인(디스코드 등)이던 상관없도록 하기 위함.
2. GetUser 말고도 다른 많은 일을 하는 외국인에게 필요한 기능만 가진 '나' 의 부탁을 들어주기 위함(Service 중심적)
=> 그렇기에 Adapter 패턴이 적용됨

2. PostgresRepository
나 (책 빌림,반납,찾기요청,예약을 하는 Service) --> 사서 선생님(그 모든 요청을 듣고 수행하는 PostgresRepository)
즉, 내가 요청 가능한 기능은 분야에 정해져 있고, 수행까지 하는 역할
이런 PostgresRepository 가 구조체를 가지는 이유는
1. 나(Service)는 여러 사서 선생님(MySQL 등)한테 부탁할 수 있지만, 다른 역할의 선생님이 아닌 "사서가 붙은" 선생님만 가능함.
=> 그렇기에 Repository 패턴이 적용됨
*/
type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

// CreateUser creates a new Slack user in the database.
// This operation executes a single INSERT statement and does not require an explicit transaction.
func (r *PostgresRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	baseEntity := database.MarkCreate()

	createSlackUserParams := sqlc.CreateSlackUserParams{
		ID:         user.ID,
		Name:       user.Name,
		RealName:   sql.NullString{String: user.RealName, Valid: user.RealName != ""},
		Email:      user.Email,
		Deleted:    user.Deleted,
		UseYn:      baseEntity.UseYn,
		CreateCode: baseEntity.CreateCode,
		CreateDate: baseEntity.CreateDate,
		Version:    baseEntity.Version,
	}

	createdSlackUser, err := r.q.CreateSlackUser(ctx, createSlackUserParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create slack user in DB: %w", err)
	}

	return &User{
		SlackUserID: createdSlackUser.SlackUserID,
		ID:          createdSlackUser.ID,
		Name:        createdSlackUser.Name,
		RealName:    createdSlackUser.RealName.String,
		Email:       createdSlackUser.Email,
		Deleted:     createdSlackUser.Deleted,
		UseYn:       createdSlackUser.UseYn,
		CreateCode:  createdSlackUser.CreateCode,
		CreateDate:  createdSlackUser.CreateDate,
		Version:     createdSlackUser.Version,
	}, nil
}
