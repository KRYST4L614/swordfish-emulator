package impl

import (
	"context"
	"slices"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type ResourceService struct {
	repo repository.ResourceRepository
}

type collection struct {
	OdataId           domain.OdataV4Id      `json:"@odata.id"`
	Members           []domain.OdataV4IdRef `json:"Members"`
	Name              domain.ResourceName   `json:"Name"`
	OdataType         domain.OdataV4Type    `json:"@odata.type"`
	MembersOdataCount domain.OdataV4Count   `json:"Members@odata.count"`
}

func NewResourceService(repository repository.ResourceRepository) *ResourceService {
	return &ResourceService{
		repo: repository,
	}
}

func (r *ResourceService) Get(ctx context.Context, resourceId string) (interface{}, error) {
	return getResource[interface{}](r.repo, ctx, resourceId)
}

func (r *ResourceService) Create(ctx context.Context, resourceId string, resource interface{}) (interface{}, error) {
	bytes, err := util.Marshal(resource)
	if err != nil {
		return nil, err
	}

	err = r.repo.Create(ctx, &dto.ResourceDto{
		Id:   resourceId,
		Data: bytes,
	})
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *ResourceService) CreateCollection(ctx context.Context, collectionDto dto.CollectionDto) (interface{}, error) {
	newCollection := &collection{
		OdataId:           collectionDto.OdataId,
		OdataType:         collectionDto.OdataType,
		Members:           make([]domain.OdataV4IdRef, 0),
		MembersOdataCount: 0,
		Name:              collectionDto.Name,
	}

	bytes, err := util.Marshal(newCollection)
	if err != nil {
		return nil, err
	}

	err = r.repo.Create(ctx, &dto.ResourceDto{
		Id:   collectionDto.OdataId,
		Data: bytes,
	})
	if err != nil {
		return nil, err
	}

	return newCollection, nil
}

func (r *ResourceService) Update(ctx context.Context, resourceId string, patchData []byte) (interface{}, error) {
	resourceDto, err := r.repo.Get(ctx, resourceId)
	if err != nil {
		return nil, err
	}

	patchedData, err := jsonpatch.MergePatch(resourceDto.Data, patchData)
	if err != nil {
		return nil, errlib.ErrInternal
	}

	resourceDto.Data = patchedData
	err = r.repo.Update(ctx, resourceDto)
	if err != nil {
		return nil, err
	}

	resource, err := util.Unmarshal[interface{}](resourceDto.Data)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *ResourceService) Replace(ctx context.Context, resourceId string, resource interface{}) (interface{}, error) {
	_, err := r.Delete(ctx, resourceId)
	if err != nil {
		return nil, err
	}

	data, err := util.Marshal(resource)
	if err != nil {
		return nil, err
	}

	err = r.repo.Create(ctx, &dto.ResourceDto{
		Id:   resourceId,
		Data: data,
	})

	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *ResourceService) Delete(ctx context.Context, resourceId string) (interface{}, error) {
	resource, err := r.Get(ctx, resourceId)
	if err != nil {
		return nil, err
	}

	err = r.repo.DeleteStartsWith(ctx, resourceId)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *ResourceService) AddResourceToCollection(ctx context.Context, collectionId, resourceId string, resource interface{}) (interface{}, error) {
	collection, err := getResource[collection](r.repo, ctx, collectionId)
	if err != nil {
		return nil, err
	}

	createdResource, err := r.Create(ctx, resourceId, resource)
	if err != nil {
		return nil, err
	}

	collection.Members = append(collection.Members, domain.OdataV4IdRef{OdataId: &resourceId})
	collection.MembersOdataCount += 1

	collectionPatchBytes, err := util.Marshal(collection)
	if err != nil {
		return nil, err
	}

	_, err = r.Update(ctx, collectionId, collectionPatchBytes)

	return createdResource, err
}

func (r *ResourceService) DeleteResourceFromCollection(ctx context.Context, collectionId, resourceId string) (interface{}, error) {
	collection, err := getResource[collection](r.repo, ctx, collectionId)
	if err != nil {
		return nil, err
	}

	resource, err := r.Delete(ctx, resourceId)
	if err != nil {
		logrus.Tracef("Resource wasn't deleted due to error: %e", err)
		return nil, err
	}

	collection.Members = slices.DeleteFunc(collection.Members, func(ref domain.OdataV4IdRef) bool {
		return *ref.OdataId == resourceId
	})
	collection.MembersOdataCount -= 1

	collectionPatchBytes, err := util.Marshal(collection)
	if err != nil {
		return nil, err
	}

	_, err = r.Update(ctx, collectionId, collectionPatchBytes)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
