package completebullet

import (
	"context"
	"errors"

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

func (s *PostgresStore) Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	query := `
		SELECT id, kind, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE id = $1 AND user_id = $2
	`

	bullet := &entity.Bullet{}
	err := transactor.GetQueryer(ctx, s.pool).QueryRow(ctx, query, id, userID).Scan(
		&bullet.ID,
		&bullet.Kind,
		&bullet.Signifier,
		&bullet.Content,
		&bullet.UserID,
		&bullet.CreatedAt,
		&bullet.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return bullet, nil
}

func (s *PostgresStore) UpdateSignifier(ctx context.Context, bullet *entity.Bullet) error {
	query := `
		UPDATE bullets
		SET signifier = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
	`

	tag, err := transactor.GetQueryer(ctx, s.pool).Exec(ctx, query,
		bullet.Signifier,
		bullet.UpdatedAt,
		bullet.ID,
		bullet.UserID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return entity.ErrNotFound
	}

	return nil
}
