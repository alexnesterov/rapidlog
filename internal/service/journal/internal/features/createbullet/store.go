package createbullet

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Insert(ctx context.Context, bullet *entity.Bullet) error {
	query := `
		INSERT INTO bullets (id, kind, signifier, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := transactor.GetQueryer(ctx, s.pool).Exec(ctx, query,
		bullet.ID,
		bullet.Kind,
		bullet.Signifier,
		bullet.Content,
		bullet.UserID,
		bullet.CreatedAt,
		bullet.UpdatedAt,
	)
	return err
}
