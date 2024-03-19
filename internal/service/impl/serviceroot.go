package impl

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type ServiceRootService struct {
	repository repository.ResourceRepository
}

func NewServiceRootService(repo repository.ResourceRepository) *ServiceRootService {
	return &ServiceRootService{
		repository: repo,
	}
}

func (s *ServiceRootService) Create(ctx context.Context, resource *domain.ServiceRoot) error {
	bytes, err := util.Marshal(resource)
	if err != nil {
		return err
	}

	err = s.repository.Create(ctx, &dto.ResourceDto{
		Id:   resource.ODataId,
		Data: bytes,
	})

	return err
}

func (s *ServiceRootService) Get(ctx context.Context, id string) (*domain.ServiceRoot, error) {
	dto, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	resource, err := util.Unmarshal[domain.ServiceRoot](dto.Data)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
