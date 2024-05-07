package impl

import (
	"context"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

func getResource[T any](r repository.ResourceRepository, ctx context.Context, id string) (*T, error) {
	dto, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	resource, err := util.Unmarshal[T](dto.Data)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
