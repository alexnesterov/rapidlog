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

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, created_at)
		VALUES ($1, $2)
	`

	_, err := getQueryer(ctx, r.pool).Exec(ctx, query,
		user.ID,
		user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Get(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, created_at
		FROM users
		WHERE id = $1
	`

	row := getQueryer(ctx, r.pool).QueryRow(ctx, query, id)

	user := &entity.User{}
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, port.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

var _ port.UserRepository = &userRepository{}
