package database

import (
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database/sqlc"

	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
)

type DB struct {
	Conn    *sql.DB
	Queries *sqlc.Queries
}

func Connect() (*DB, error) {
	dsn := makeDsn()

	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

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
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		username, password, host, port, database, sslMode)
}
