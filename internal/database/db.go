package database

import (
	"database/sql"
	"fmt"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	_ "github.com/lib/pq"
	"time"
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
