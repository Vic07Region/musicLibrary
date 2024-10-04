package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	_ "github.com/lib/pq"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func NewConnection(dbdriver, connection string) (*sql.DB, error) {
	dbo, err := sql.Open(dbdriver, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Проверьте соединение
	if err = dbo.Ping(); err != nil {
		dbo.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return dbo, nil
}

func NewStorage(db *sql.DB, log *logger.Logger) *Queries {
	return &Queries{db: db, log: log}
}

type Queries struct {
	db  *sql.DB
	log *logger.Logger
}
