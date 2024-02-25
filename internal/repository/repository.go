package repository

import (
	provider "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/psql"
)

type Repository struct {
	VolumeRepo VolumeRepository
}

func RepositoryFactory(provider *provider.DbProvider) *Repository {
	return &Repository{
		VolumeRepo: psql.NewVolumeRepository(provider),
	}
}
