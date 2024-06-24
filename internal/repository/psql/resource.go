package psql

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
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
		if strings.Contains(err.Error(), "resource_pkey") {
			return fmt.Errorf("%w, resource with id %s already exists", errlib.ErrResourceAlreadyExists, resource.Id)
		}
		return fmt.Errorf("%w.", errlib.ErrInternal)
	}

	return nil
}

func (s *PsqlResourceRepository) Update(ctx context.Context, resource *dto.ResourceDto) error {
	query := `UPDATE resource SET data = :data WHERE id = :id`
	if _, err := s.db.NamedExecContext(ctx, query, resource); err != nil {
		return err
	}
	return nil
}

func (s *PsqlResourceRepository) Get(ctx context.Context, id string) (*dto.ResourceDto, error) {
	var dto dto.ResourceDto
	if err := s.db.GetContext(ctx, &dto, `SELECT * FROM resource WHERE id=$1`, id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "no_data":
				return nil, fmt.Errorf("%w, resource with id %s doesn't exist", errlib.ErrNotFound, id)
			}
		}
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("%w, resource with id %s doesn't exist", errlib.ErrNotFound, id)
		}
		return nil, fmt.Errorf("%w.", errlib.ErrInternal)
	}

	return &dto, nil
}

func (s *PsqlResourceRepository) DeleteAll(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `TRUNCATE resource CASCADE`); err != nil {
		return fmt.Errorf("%w.", errlib.ErrInternal)
	}
	return nil
}

func (s *PsqlResourceRepository) DeleteById(ctx context.Context, id string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM resource WHERE id=$1`, id); err != nil {
		return fmt.Errorf("%w.", errlib.ErrInternal)
	}
	return nil
}

func (s *PsqlResourceRepository) DeleteStartsWith(ctx context.Context, prefix string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM resource WHERE id LIKE $1`, prefix+"%"); err != nil {
		return err
	}
	return nil
}
