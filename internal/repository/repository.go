package repository

import (
	provider "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/psql"
)

type Repository struct {
	ResourceRepository ResourceRepository
}

func RepositoryFactory(provider *provider.DbProvider) *Repository {
	return &Repository{
		ResourceRepository: psql.NewPsqlResourceRepository(provider),
	}
}
