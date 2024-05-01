package service

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
)

//go:generate mockgen --build_flags=--mod=mod -destination mock/mock_service.go . ServiceRootService,StorageService,StorageCollectionService

type ServiceRootService interface {
	Create(ctx context.Context, resource *domain.ServiceRoot) error
	Get(ctx context.Context, id string) (*domain.ServiceRoot, error)
}

type StorageService interface {
	Get(ctx context.Context, id string) (*domain.Storage, error)
	Replace(ctx context.Context, storageId string, storage *domain.Storage) (*domain.Storage, error)
	Update(ctx context.Context, storageId string, patchData []byte) (*domain.Storage, error)
}

type StorageCollectionService interface {
	Get(ctx context.Context, id string) (*domain.StorageCollection, error)
	AddStorage(ctx context.Context, collectionId string, storage *domain.Storage) error
	DeleteStorage(ctx context.Context, storageId string) (*domain.Storage, error)
}
