package repository

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
)

//go:generate mockgen --build_flags=--mod=mod -destination mock/mock_repository.go . ResourceRepository

type ResourceRepository interface {
	Create(ctx context.Context, resource *dto.ResourceDto) error
	Get(ctx context.Context, id string) (*dto.ResourceDto, error)
	Update(ctx context.Context, resource *dto.ResourceDto) error
	DeleteAll(ctx context.Context) error
	DeleteById(ctx context.Context, id string) error
	DeleteStartsWith(ctx context.Context, prefix string) error
}
