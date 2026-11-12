package postgres

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
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

func (r *BulletRepository) Create(ctx context.Context, bullet *entity.Bullet) error {
	query := `
		INSERT INTO bullets (id, kind, signifier, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := transactor.GetQueryer(ctx, r.pool).Exec(ctx, query,
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

func (r *BulletRepository) Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	query := `
		SELECT id, kind, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE id = $1 AND user_id = $2
	`

	row := transactor.GetQueryer(ctx, r.pool).QueryRow(ctx, query, id, userID)

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

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return bullet, nil
}

func (r *BulletRepository) List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
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
		func(row pgx.CollectableRow) (*entity.Bullet, error) {
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

	if err != nil {
		return nil, err
	}

	return bullets, nil
}

func (r *BulletRepository) Update(ctx context.Context, bullet *entity.Bullet) error {
	query := `
		UPDATE bullets
		SET kind = $2, signifier = $3, content = $4, created_at = $5, updated_at = $6
		WHERE id = $1 AND user_id = $7
	`

	_, err := transactor.GetQueryer(ctx, r.pool).Exec(ctx, query,
		bullet.ID,
		bullet.Kind,
		bullet.Signifier,
		bullet.Content,
		bullet.CreatedAt,
		bullet.UpdatedAt,
		bullet.UserID,
	)

	return err
}
