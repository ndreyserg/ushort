package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConnection interface {
	PingContext(ctx context.Context) error
	Close() error
}

//mockgen -destination=./internal/app/mocks/mock_db.go -package=mocks github.com/ndreyserg/ushort/internal/app/db DBConnection

type dbConn struct {
	db *sql.DB
}

func (conn *dbConn) PingContext(ctx context.Context) error {
	return conn.db.PingContext(ctx)
}

func (conn *dbConn) Close() error {
	return conn.db.Close()
}


func MakeConnect(DSN string) (DBConnection, error) {
	db, err := sql.Open("pgx", DSN)

	if err != nil {
		return nil, err
	}

	return &dbConn{
		db: db,
	}, nil
}
