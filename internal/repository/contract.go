package repository

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
)

type VolumeRepository interface {
	Create()
	Delete()
	Modify()
}

type ServiceRootRepository interface {
	Create(ctx context.Context, serviceRoot domain.ServiceRoot) error
	Get(ctx context.Context, id string) (*domain.ServiceRoot, error)
}
