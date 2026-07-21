package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type CollectionRepository interface {
	Create(collection *entity.Collection) error
	List() ([]*entity.Collection, error)
	Read(id uuid.UUID) (*entity.Collection, error)
	Update(collection *entity.Collection) error
	Delete(id uuid.UUID) error
}
