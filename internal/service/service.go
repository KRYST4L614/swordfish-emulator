package service

import (
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service/volume"
)

type Service struct {
	VolumeService VolumeService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		VolumeService: volume.NewVolumeService(repository.VolumeRepo),
	}
}
