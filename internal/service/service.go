package service

import (
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service/impl"
)

type Service struct {
	ServiceRootService ServiceRootService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		ServiceRootService: impl.NewServiceRootService(repository.ResourceRepository),
	}
}
