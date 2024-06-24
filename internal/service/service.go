package service

import (
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
)

type Service struct {
	ResourceService ResourceService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		ResourceService: NewResourceService(repository.ResourceRepository, NewSimpleIdGenerator()),
	}
}
