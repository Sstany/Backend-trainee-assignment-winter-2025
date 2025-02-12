package repo

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, migrationsDir embed.FS) error {
	goose.SetBaseFS(migrationsDir)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set 'postgres' dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}

	return nil
}
