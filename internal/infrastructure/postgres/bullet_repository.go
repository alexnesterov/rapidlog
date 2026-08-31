package postgres

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bulletRepository struct {
	pool *pgxpool.Pool
}

func NewBulletRepository(pool *pgxpool.Pool) *bulletRepository {
	return &bulletRepository{pool: pool}
}

func (r *bulletRepository) Create(ctx context.Context, bullet *entity.Bullet) error {
	query := `
		INSERT INTO bullets (id, type, signifier, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := getQueryer(ctx, r.pool).Exec(ctx, query,
		bullet.ID,
		bullet.Type,
		bullet.Signifier,
		bullet.Content,
		bullet.UserID,
		bullet.CreatedAt,
		bullet.UpdatedAt,
	)

	return err
}

func (r *bulletRepository) Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	query := `
		SELECT id, type, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE id = $1 AND user_id = $2
	`

	row := getQueryer(ctx, r.pool).QueryRow(ctx, query, id, userID)

	bullet := &entity.Bullet{}
	err := row.Scan(
		&bullet.ID,
		&bullet.Type,
		&bullet.Signifier,
		&bullet.Content,
		&bullet.UserID,
		&bullet.CreatedAt,
		&bullet.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, port.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return bullet, nil
}

func (r *bulletRepository) List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
	query := `
		SELECT id, type, signifier, content, user_id, created_at, updated_at
		FROM bullets
		WHERE user_id = $1
	`

	rows, err := getQueryer(ctx, r.pool).Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bullets, err := pgx.CollectRows(rows,
		func(row pgx.CollectableRow) (*entity.Bullet, error) {
			bullet := &entity.Bullet{}
			err := row.Scan(
				&bullet.ID,
				&bullet.Type,
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

func (r *bulletRepository) Update(ctx context.Context, bullet *entity.Bullet) error {
	query := `
		UPDATE bullets
		SET type = $2, signifier = $3, content = $4, updated_at = $5
		WHERE id = $1 AND user_id = $6
	`

	_, err := getQueryer(ctx, r.pool).Exec(ctx, query,
		bullet.ID,
		bullet.Type,
		bullet.Signifier,
		bullet.Content,
		bullet.UpdatedAt,
		bullet.UserID,
	)

	return err
}

var _ port.BulletRepository = &bulletRepository{}
