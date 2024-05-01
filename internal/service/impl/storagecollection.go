package impl

import (
	"context"
	"slices"

	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type StorageCollectionService struct {
	repository repository.ResourceRepository
}

func NewStorageCollectionService(repo repository.ResourceRepository) *StorageCollectionService {
	return &StorageCollectionService{
		repository: repo,
	}
}

func (s *StorageCollectionService) Get(ctx context.Context, id string) (*domain.StorageCollection, error) {
	return getResource[domain.StorageCollection](s.repository, ctx, id)
}

func (s *StorageCollectionService) AddStorage(ctx context.Context, collectionId string, storage *domain.Storage) error {
	*storage.OdataId = collectionId + "/" + storage.Id
	err := s.createStorage(ctx, storage)
	if err != nil {
		return err
	}

	storageCollection, err := s.Get(ctx, collectionId)
	if err != nil {
		return err
	}

	*storageCollection.Members = append(*storageCollection.Members, domain.OdataV4IdRef{OdataId: storage.OdataId})
	*storageCollection.MembersOdataCount = *storageCollection.MembersOdataCount + 1

	collectionBytes, err := util.Marshal(storageCollection)
	if err != nil {
		return err
	}
	err = s.repository.Update(ctx, &dto.ResourceDto{
		Id:   collectionId,
		Data: collectionBytes,
	})

	return err
}

func (s *StorageCollectionService) DeleteStorage(ctx context.Context, storageId string) (*domain.Storage, error) {
	storage, err := getResource[domain.Storage](s.repository, ctx, storageId)
	if err != nil {
		return nil, err
	}

	err = s.deleteStorage(ctx, storageId)
	if err != nil {
		logrus.Tracef("Storage wasn't deleted due to error: %e", err)
		return nil, err
	}

	collectionId := getParent(storageId)
	storageCollection, err := s.Get(ctx, collectionId)
	if err != nil {
		return nil, err
	}

	*storageCollection.Members = slices.DeleteFunc(*storageCollection.Members, func(ref domain.OdataV4IdRef) bool {
		return *ref.OdataId == storageId
	})
	*storageCollection.MembersOdataCount = *storageCollection.MembersOdataCount - 1

	collectionBytes, err := util.Marshal(storageCollection)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(ctx, &dto.ResourceDto{
		Id:   collectionId,
		Data: collectionBytes,
	})

	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *StorageCollectionService) createStorage(ctx context.Context, resource *domain.Storage) error {
	bytes, err := util.Marshal(resource)
	if err != nil {
		return err
	}

	err = s.repository.Create(ctx, &dto.ResourceDto{
		Id:   *resource.OdataId,
		Data: bytes,
	})

	return err
}

func (s *StorageCollectionService) deleteStorage(ctx context.Context, storageId string) error {
	return s.repository.DeleteStartsWith(ctx, storageId)
}
