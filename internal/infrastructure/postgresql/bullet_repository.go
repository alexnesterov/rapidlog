package postgresql

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bulletPostgresql struct {
	db *pgxpool.Pool
}

func NewPostgresqlBulletRepository(db *pgxpool.Pool) *bulletPostgresql {
	return &bulletPostgresql{db: db}
}

func (r *bulletPostgresql) Create(bullet *entity.Bullet) error {
	return nil
}

func (r *bulletPostgresql) List() ([]*entity.Bullet, error) {
	return nil, nil
}

func (r *bulletPostgresql) Get(id uuid.UUID) (*entity.Bullet, error) {
	return nil, nil
}

func (r *bulletPostgresql) Update(bullet *entity.Bullet) error {
	return nil
}

var _ port.BulletRepository = &bulletPostgresql{}
