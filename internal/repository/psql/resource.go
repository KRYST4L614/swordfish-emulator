package psql

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
)

type PsqlResourceRepository struct {
	db *provider.DbProvider
}

func NewPsqlResourceRepository(db *provider.DbProvider) *PsqlResourceRepository {
	return &PsqlResourceRepository{
		db: db,
	}
}

func (s *PsqlResourceRepository) Create(ctx context.Context, resource *dto.ResourceDto) error {
	query := `INSERT INTO resource (id, data) VALUES (:id, :data)`
	_, err := s.db.NamedExecContext(ctx, query, resource)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return fmt.Errorf("%w. %s", errlib.ErrResourceExists, err.Error())
			}
		}
		return fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	return nil
}

func (s *PsqlResourceRepository) Get(ctx context.Context, id string) (*dto.ResourceDto, error) {
	var dto dto.ResourceDto
	if err := s.db.GetContext(ctx, &dto, `SELECT * FROM resource WHERE id=$1`, id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "no_data":
				return nil, fmt.Errorf("%w. %s", errlib.ErrNotFound, err.Error())
			}
		}
		return nil, fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	return &dto, nil
}
