package service

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
)

//go:generate mockgen --build_flags=--mod=mod -destination mock/mock_service.go . ServiceRootService

type ServiceRootService interface {
	Create(ctx context.Context, resource *domain.ServiceRoot) error
	Get(ctx context.Context, id string) (*domain.ServiceRoot, error)
}
