package migrate

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

const (
	POSTGRES = "postgres"
	SQLITE   = "sqlite"
)

func ApplyMigrations(db *sql.DB, migrationsDir string, dbdriver string) error {
	// Инициализируем Goose
	if err := goose.SetDialect(dbdriver); err != nil {
		panic(err)
	}
	// Проходим по всем файлам миграций
	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
