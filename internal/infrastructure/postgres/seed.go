package postgres

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed seed.sql
var seedSQL []byte

func Seed(ctx context.Context, pool *pgxpool.Pool) error {
	query := `SELECT count(*) FROM bullets`

	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = pool.Exec(ctx, string(seedSQL))
	return err
}
