package impl

import (
	"context"
	"path/filepath"

	jsonpatch "github.com/evanphx/json-patch"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageService struct {
	repository repository.ResourceRepository
}

func NewStorageService(repo repository.ResourceRepository) *StorageService {
	return &StorageService{
		repository: repo,
	}
}

func (s *StorageService) Get(ctx context.Context, id string) (*domain.Storage, error) {
	return getResource[domain.Storage](s.repository, ctx, id)
}

func (s *StorageService) Replace(ctx context.Context, storageId string, storage *domain.Storage) (*domain.Storage, error) {
	storage.Id = filepath.Base(storageId)
	*storage.OdataId = storageId
	err := s.repository.DeleteStartsWith(ctx, storageId)
	if err != nil {
		return nil, err
	}

	data, err := util.Marshal(storage)
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, &dto.ResourceDto{
		Id:   storageId,
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *StorageService) Update(ctx context.Context, storageId string, patchData []byte) (*domain.Storage, error) {
	storage, err := s.repository.Get(ctx, storageId)
	if err != nil {
		return nil, err
	}

	patchedData, err := jsonpatch.MergePatch(storage.Data, patchData)
	if err != nil {
		return nil, errlib.ErrInternal
	}

	storage.Data = patchedData
	err = s.repository.Update(ctx, storage)
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, storageId)
}
