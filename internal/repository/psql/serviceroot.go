package psql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
)

type ServiceRootRepository struct {
	db *provider.DbProvider
}

func NewServiceRootRepository(db *provider.DbProvider) *ServiceRootRepository {
	return &ServiceRootRepository{
		db: db,
	}
}

type serviceRootDto struct {
	Id   string `db:"id"`
	Data string `db:"data"`
}

func (s *ServiceRootRepository) Create(ctx context.Context, serviceRoot domain.ServiceRoot) error {
	data, err := json.Marshal(serviceRoot)
	if err != nil {
		return fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	dto := &serviceRootDto{
		Id:   serviceRoot.ODataId,
		Data: string(data),
	}

	query := `INSERT INTO resource (id, data) VALUES (:id, :data)`
	_, err = s.db.NamedExecContext(ctx, query, dto)
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

func (s *ServiceRootRepository) Get(ctx context.Context, id string) (*domain.ServiceRoot, error) {
	var dto serviceRootDto
	if err := s.db.GetContext(ctx, &dto, `SELECT * FROM resource WHERE id=$1`, id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "no_data":
				return nil, fmt.Errorf("%w. %s", errlib.ErrNotFound, err.Error())
			}
		}
		return nil, fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	serviceRoot := &domain.ServiceRoot{}
	err := json.Unmarshal([]byte(dto.Data), &serviceRoot)
	if err != nil {
		return nil, fmt.Errorf("%w. %s", errlib.ErrInternal, err.Error())
	}

	return serviceRoot, nil
}
