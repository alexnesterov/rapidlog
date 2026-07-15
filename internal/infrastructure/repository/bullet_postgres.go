package repository

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bulletPostgres struct {
	db *pgxpool.Pool
}

func NewPostgresBulletRepository(db *pgxpool.Pool) *bulletPostgres {
	return &bulletPostgres{db: db}
}

func (r *bulletPostgres) Create(bullet *entity.Bullet) error {
	return nil
}

func (r *bulletPostgres) List() ([]*entity.Bullet, error) {
	return nil, nil
}

func (r *bulletPostgres) Read(id uuid.UUID) (*entity.Bullet, error) {
	return nil, nil
}

func (r *bulletPostgres) Update(bullet *entity.Bullet) error {
	return nil
}

func (r *bulletPostgres) Delete(id uuid.UUID) error {
	return nil
}

var _ port.BulletRepository = &bulletPostgres{}
