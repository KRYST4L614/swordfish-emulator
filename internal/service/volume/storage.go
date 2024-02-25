package volume

import (
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
)

type VolumeService struct {
	repository repository.VolumeRepository
}

func NewVolumeService(repository repository.VolumeRepository) *VolumeService {
	return &VolumeService{
		repository: repository,
	}
}

func (s *VolumeService) Ceate() {
	panic("not implemented") // TODO: Implement
}

func (s *VolumeService) Delete() {
	panic("not implemented") // TODO: Implement
}

func (s *VolumeService) Get() {
	panic("not implemented") // TODO: Implement
}
