package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database/sqlc"

	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
)

/*
DB is struct for database connection
이 구조는 다음 장점을 가짐
1. DB 커넥샨을 Wrapping한 Queries 를 사용 시마다 생성하지 않고, 하나로 여러 서비스에서 사용 가능
2. Queries 를 만드려면 Conn 이 있어야 하기에 함께 캡슐링하여 관리 용이
3. 만약 Queries 에 없는 저수준 메서드(트랜잭션 처리 등)를 *sql.DB 를 통해 사용해야 하는 경우 용이
4. Service가 의존하는 Interface의 구현체의 내부 필드로 2가지를 주면됨.

계획중인 서비스 비즈니스 계층 긴 구조는
[Service] --depends on--> [Repository interface] <--implements-- [Repository struct (DB 사용)]
이며, sqlc는 아래와 같다.
[*sqlc.Queries] --wraps--> [*sql.DB] --implements--> [DBTX interface]

이 구조의 이점은
1. sqlc가 아닌 다른 것으로 바꿔도, Service 단에 영향 많이 X -> 구현만 잘해놓은 구조체가 있으면 됨
2. 테스트 용이 -> Repository interface 자체를 Mocking 하면 테스트 매우 쉬움
*/
type DB struct {
	Conn    *sql.DB
	Queries *sqlc.Queries
}

func Connect() (*DB, error) {
	dsn := makeDsn()

	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	slog.Info("successfully connected to database")

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	slog.Info("successfully pinged to database")

	return &DB{
		Conn:    conn,
		Queries: sqlc.New(conn),
	}, nil
}

// makeDsn is function for PlainText Connection to Database.
// need to add(refactor) function for TLS Connection to Database!
func makeDsn() string {
	host := config.Cfg.Database.Host
	port := config.Cfg.Database.Port
	username := config.Cfg.Database.Username
	password := config.Cfg.Database.Password
	database := config.Cfg.Database.Database
	sslMode := config.Cfg.Database.SslMode

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		username, password, host, port, database, sslMode)
}
