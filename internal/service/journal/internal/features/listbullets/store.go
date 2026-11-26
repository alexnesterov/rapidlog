package listbullets

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
	query := `
		SELECT id, kind, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE user_id = $1
	`

	rows, err := transactor.GetQueryer(ctx, s.pool).Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*entity.Bullet, error) {
		bullet := &entity.Bullet{}
		err := row.Scan(
			&bullet.ID,
			&bullet.Kind,
			&bullet.Signifier,
			&bullet.Content,
			&bullet.UserID,
			&bullet.CreatedAt,
			&bullet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		return bullet, nil
	})
}
