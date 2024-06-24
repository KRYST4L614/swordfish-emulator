package service

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
)

//go:generate mockgen --build_flags=--mod=mod -destination mock/mock_service.go . ResourceService

type ResourceService interface {
	Get(ctx context.Context, resourceId string) (interface{}, error)
	Create(ctx context.Context, resourceId string, resource interface{}) (interface{}, error)
	Update(ctx context.Context, resourceId string, patchData []byte) (interface{}, error)
	Replace(ctx context.Context, resourceId string, resource interface{}) (interface{}, error)
	Delete(ctx context.Context, resourceId string) (interface{}, error)
	CreateCollection(ctx context.Context, collectionDto dto.CollectionDto) (interface{}, error)
	AddResourceToCollection(ctx context.Context, resourceDto dto.ResourceRequestDto) (interface{}, error)
	DeleteResourceFromCollection(ctx context.Context, collectionId, resourceId string) (interface{}, error)
}

type IdGenerator interface {
	Generate(count uint64) (string, error)
}
