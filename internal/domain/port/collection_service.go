package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateCollectionRequest struct {
	Topic string
}

type UpdateCollectionRequest struct {
	ID    uuid.UUID
	Topic string
}

type CollectionService interface {
	CreateCollection(req CreateCollectionRequest) (*entity.Collection, error)
	ListCollections() ([]*entity.Collection, error)
	ReadCollection(id uuid.UUID) (*entity.Collection, error)
	UpdateCollection(req UpdateCollectionRequest) (*entity.Collection, error)
	DeleteCollection(id uuid.UUID) error
}
