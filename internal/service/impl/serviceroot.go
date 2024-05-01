package impl

import (
	"context"
	"fmt"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
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

	if resource.OdataId == nil {
		return fmt.Errorf("%w: ODataId in ServiceRoot is nil", errlib.ErrInternal)
	}

	err = s.repository.Create(ctx, &dto.ResourceDto{
		Id:   *resource.OdataId,
		Data: bytes,
	})

	return err
}

func (s *ServiceRootService) Get(ctx context.Context, id string) (*domain.ServiceRoot, error) {
	return getResource[domain.ServiceRoot](s.repository, ctx, id)
}
