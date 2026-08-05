package postgres

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresBulletRepository struct {
	db *pgxpool.Pool
}

func NewPostgresBulletRepository(db *pgxpool.Pool) *postgresBulletRepository {
	return &postgresBulletRepository{db: db}
}

func (r *postgresBulletRepository) Create(ctx context.Context, bullet *entity.Bullet) error {
	return nil
}

func (r *postgresBulletRepository) List(ctx context.Context) ([]*entity.Bullet, error) {
	return nil, nil
}

func (r *postgresBulletRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Bullet, error) {
	return nil, nil
}

func (r *postgresBulletRepository) Update(ctx context.Context, bullet *entity.Bullet) error {
	return nil
}

var _ port.BulletRepository = &postgresBulletRepository{}
