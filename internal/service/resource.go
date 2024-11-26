package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	jsonpatch "github.com/evanphx/json-patch"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type resourceService struct {
	repo      repository.ResourceRepository
	generator IdGenerator
}

type collection struct {
	OdataId           domain.OdataV4Id      `json:"@odata.id"`
	Members           []domain.OdataV4IdRef `json:"Members"`
	Name              domain.ResourceName   `json:"Name"`
	OdataType         domain.OdataV4Type    `json:"@odata.type"`
	MembersOdataCount domain.OdataV4Count   `json:"Members@odata.count"`
}

func NewResourceService(repository repository.ResourceRepository, generator IdGenerator) *resourceService {
	return &resourceService{
		repo:      repository,
		generator: generator,
	}
}

// Get - just gets resource by 'resourceId'
func (r *resourceService) Get(ctx context.Context, resourceId string) (interface{}, error) {
	return getResource[interface{}](r.repo, ctx, resourceId)
}

// Create - creates resource with 'resourceId'
func (r *resourceService) Create(ctx context.Context, resourceId string, resource interface{}) (interface{}, error) {
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

// CreateCollection is just wrapper to return untyped collection and satisfy interface
func (r *resourceService) CreateCollection(ctx context.Context, collectionDto dto.CollectionDto) (interface{}, error) {
	return r.createCollection(ctx, collectionDto)
}

// createCollection creates empty collection with zero members contained
func (r *resourceService) createCollection(ctx context.Context, collectionDto dto.CollectionDto) (*collection, error) {
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

	// TODO: have to add Reference to a collection in owning resource

	if err != nil {
		return nil, err
	}

	return newCollection, nil
}

// Update uses 'patchData' to update resource with 'resourceId'. Returns updated resource.
//
// 'patchData' should contain new state of resource fields. Read 'MergePatch' https://datatracker.ietf.org/doc/html/rfc7386
// for more info.
func (r *resourceService) Update(ctx context.Context, resourceId string, patchData []byte) (interface{}, error) {
	resourceDto, err := r.repo.Get(ctx, resourceId)
	if err != nil {
		return nil, err
	}

	patchedData, err := jsonpatch.MergePatch(resourceDto.Data, patchData)
	if err != nil {
		return nil, errlib.ErrInternal
	}

	resourceDto.Data = patchedData

	resource, err := util.Unmarshal[interface{}](resourceDto.Data)
	if err != nil {
		return nil, err
	}

	//TODO: Remake signature of Update method

	// err = validate(r, ctx, resource)
	// if err != nil {
	// 	return nil, err
	// }

	err = r.repo.Update(ctx, resourceDto)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// Replace - replaces old resources placed by 'resourceId' with new provided data
func (r *resourceService) Replace(ctx context.Context, resourceId string, resource interface{}) (interface{}, error) {
	err := validate(r, ctx, resource)
	if err != nil {
		return nil, err
	}

	_, err = r.Delete(ctx, resourceId)
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

// Delete just removes resource by 'resourceId'
func (r *resourceService) Delete(ctx context.Context, resourceId string) (interface{}, error) {
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

// AddResourceToCollection adds new resource in provided collection
func (r *resourceService) AddResourceToCollection(ctx context.Context, resourceDto dto.ResourceRequestDto) (interface{}, error) {
	collection, err := getResource[collection](r.repo, ctx, resourceDto.Collection.OdataId)
	if err != nil {
		return nil, err
	}

	err = validate(r, ctx, resourceDto.Resource)

	switch {
	case err == nil:
		// Do nothing if no error happened
	case errors.Is(err, errlib.ErrNotFound):
		collection, err = r.createCollection(ctx, resourceDto.Collection)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	collection.MembersOdataCount += 1

	resourceId, err := r.populateResource(collection.MembersOdataCount, resourceDto)
	if err != nil {
		return nil, err
	}

	collection.Members = append(collection.Members, domain.OdataV4IdRef{OdataId: &resourceId})

	createdResource, err := r.Create(ctx, resourceId, resourceDto.Resource)
	if err != nil {
		return nil, err
	}

	collectionPatchBytes, err := util.Marshal(collection)
	if err != nil {
		return nil, err
	}

	_, err = r.Update(ctx, resourceDto.Collection.OdataId, collectionPatchBytes)
	if err != nil {
		return nil, err
	}

	return createdResource, nil
}

// DeleteResourceFromCollection deletes resource with 'resourceId' from provided collection.
//
// Returns deleted resource.
func (r *resourceService) DeleteResourceFromCollection(ctx context.Context, collectionId, resourceId string) (interface{}, error) {
	collection, err := getResource[collection](r.repo, ctx, collectionId)
	if err != nil {
		return nil, err
	}

	resource, err := r.Delete(ctx, resourceId)
	if err != nil {
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

// populateResource - populates resource with some internal properties, such as Id
//
// Returns OdataId of populated resource
func (r *resourceService) populateResource(count int64, resourceDto dto.ResourceRequestDto) (string, error) {
	if resourceDto.OdataType == "" || resourceDto.Collection.OdataId == "" {
		return "", errlib.ErrInternal
	}

	if resourceDto.Id == "" {
		id, err := r.generator.Generate(uint64(count))
		if err != nil {
			return "", err
		}

		resourceDto.Id = id
		resourceDto.IdSetter(id)
	}

	odataId := resourceDto.Collection.OdataId + "/" + resourceDto.Id
	resourceDto.OdataIdSetter(odataId)
	resourceDto.OdataTypeSetter(resourceDto.OdataType)

	return odataId, nil
}

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

func validate(service ResourceService, context context.Context, resource any) error {
	v := reflect.TypeOf(resource)
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		field := v.Field(i)
		if field.Type.String() == "*[]domain.OdataV4IdRef" {
			value := reflect.ValueOf(resource)
			fieldValue := value.FieldByName(field.Name).Interface()
			ids, ok := fieldValue.(*[]domain.OdataV4IdRef)
			if !ok {
				return errlib.ErrInternal
			}
			if ids == nil {
				continue
			}
			for _, i := range *ids {
				_, err := service.Get(context, *i.OdataId)
				if err != nil {
					return fmt.Errorf("validation resource error: %w", err)
				}
			}
		}

		if field.Type.String() == "*domain.OdataV4IdRef" {
			value := reflect.ValueOf(resource)
			fieldValue := value.FieldByName(field.Name).Interface()
			id, ok := fieldValue.(*domain.OdataV4IdRef)
			if !ok {
				return errlib.ErrInternal
			}
			if id == nil {
				continue
			}
			_, err := service.Get(context, *id.OdataId)
			slog.Info(*id.OdataId)
			if err != nil {
				return fmt.Errorf("validation resource error: %w", err)
			}
		}
	}
	return nil
}
