package service

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
)

type ServiceRootService interface {
	Create(ctx context.Context, resource *domain.ServiceRoot) error
	Get(ctx context.Context, id string) (*domain.ServiceRoot, error)
}
