package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alexnesterov/rapidlog-api/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(dsn string) error {
	source, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		return fmt.Errorf("init migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, strings.Replace(dsn, "postgres://", "pgx5://", 1))
	if err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
