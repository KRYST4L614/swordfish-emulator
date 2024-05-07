package service

import (
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service/impl"
)

type Service struct {
	ResourceService ResourceService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		ResourceService: impl.NewResourceService(repository.ResourceRepository),
	}
}
