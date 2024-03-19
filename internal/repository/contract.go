package repository

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
)

//go:generate mockgen --build_flags=--mod=mod -destination mock/mock_repository.go . ResourceRepository

type ResourceRepository interface {
	Create(ctx context.Context, resource *dto.ResourceDto) error
	Get(ctx context.Context, id string) (*dto.ResourceDto, error)
}
