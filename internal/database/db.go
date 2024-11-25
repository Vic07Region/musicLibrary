package database

import (
	"database/sql"
	"fmt"
	"time" //nolint:gci

	"github.com/Vic07Region/musicLibrary/internal/lib/logger" //nolint:gci
	_ "github.com/lib/pq"                                     //nolint:gci
)

type ConnectionParams struct {
	DbDriver         string
	ConnectionString string
	MaxOpenConns     int
	MsxIdleConns     int
	MaxLifeTime      time.Duration
}

func NewConnection(params ConnectionParams) (*sql.DB, error) {
	dbo, err := sql.Open(params.DbDriver, params.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	dbo.SetMaxOpenConns(params.MaxOpenConns)
	dbo.SetMaxIdleConns(params.MsxIdleConns)
	dbo.SetConnMaxLifetime(params.MaxLifeTime)
	if err = dbo.Ping(); err != nil {
		err := dbo.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return dbo, nil
}

type Queries struct {
	db    *sql.DB
	log   *logger.Logger
	debug bool
}

func NewStorage(db *sql.DB, log *logger.Logger, debug bool) *Queries {
	return &Queries{db: db, log: log, debug: debug}
}
