/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package database

import (
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	driverName = "pgx"
)

type DB struct {
	Conn    *sql.DB
	Queries *sqlc.Queries
}

func Connect() (*DB, error) {
	dsn := MakeDsn()
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

func MakeDsn() string {
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
