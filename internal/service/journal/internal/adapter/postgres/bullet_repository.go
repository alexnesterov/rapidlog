package postgres

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BulletRepository struct {
	pool *pgxpool.Pool
}

func NewBulletRepository(pool *pgxpool.Pool) *BulletRepository {
	return &BulletRepository{pool: pool}
}

func (r *BulletRepository) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Bullet, error) {
	query := `
		SELECT id, kind, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE id = $1 AND user_id = $2
	`

	row := transactor.GetQueryer(ctx, r.pool).QueryRow(ctx, query, id, userID)

	bullet := &domain.Bullet{}
	err := row.Scan(
		&bullet.ID,
		&bullet.Kind,
		&bullet.Signifier,
		&bullet.Content,
		&bullet.UserID,
		&bullet.CreatedAt,
		&bullet.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return bullet, nil
}

func (r *BulletRepository) List(ctx context.Context, userID uuid.UUID) ([]*domain.Bullet, error) {
	query := `
		SELECT id, kind, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE user_id = $1
	`

	rows, err := transactor.GetQueryer(ctx, r.pool).Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bullets, err := pgx.CollectRows(rows,
		func(row pgx.CollectableRow) (*domain.Bullet, error) {
			bullet := &domain.Bullet{}
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

	if err != nil {
		return nil, err
	}

	return bullets, nil
}

func (r *BulletRepository) Save(ctx context.Context, bullet *domain.Bullet) error {
	query := `
		INSERT INTO bullets (id, kind, signifier, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET kind = EXCLUDED.kind,
			signifier = EXCLUDED.signifier,
			content = EXCLUDED.content,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at
		WHERE bullets.user_id = EXCLUDED.user_id
	`

	tag, err := transactor.GetQueryer(ctx, r.pool).Exec(ctx, query,
		bullet.ID,
		bullet.Kind,
		bullet.Signifier,
		bullet.Content,
		bullet.UserID,
		bullet.CreatedAt,
		bullet.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrOwnerMismatch
	}

	return nil
}
