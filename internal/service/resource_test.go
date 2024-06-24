package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	mock_repository "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/mock"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
	"go.uber.org/mock/gomock"
)

var serviceRootBytes = `
{
	"@odata.type": "#ServiceRoot.v1_15_0.ServiceRoot",
	"Id": "RootService",
	"Name": "Test Service Root",
	"@odata.id": "/redfish/v1"
}`

var serviceRoot = domain.ServiceRoot{
	OdataType: util.Addr("#ServiceRoot.v1_15_0.ServiceRoot"),
	Id:        "RootService",
	Name:      "Test Service Root",
	OdataId:   util.Addr("/redfish/v1"),
}

func Test_resourceService_getResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)

	repo.EXPECT().Get(gomock.Any(), "/redfish/v1").AnyTimes().Return(
		&dto.ResourceDto{
			Id:   "/redfish/v1",
			Data: []byte(serviceRootBytes),
		},
		nil,
	)
	repo.EXPECT().Get(gomock.Any(), "/redfish/v1/not/exist").AnyTimes().Return(
		nil,
		errlib.ErrNotFound,
	)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().Return(
		nil,
		errlib.ErrInternal,
	)

	type args struct {
		r   repository.ResourceRepository
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ServiceRoot
		wantErr error
	}{
		{
			name: "Get existing resource",
			args: args{
				r:   repo,
				ctx: context.Background(),
				id:  "/redfish/v1",
			},
			want:    &serviceRoot,
			wantErr: nil,
		},
		{
			name: "Get not existing resource",
			args: args{
				r:   repo,
				ctx: context.Background(),
				id:  "/redfish/v1/not/exist",
			},
			want:    nil,
			wantErr: errlib.ErrNotFound,
		},
		{
			name: "Get with some unknown error",
			args: args{
				r:   repo,
				ctx: context.Background(),
				id:  "any",
			},
			want:    nil,
			wantErr: errlib.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResource[domain.ServiceRoot](tt.args.r, tt.args.ctx, tt.args.id)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					assert.Equal(t, tt.want, got)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_resourceService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)

	repo.EXPECT().Create(
		gomock.Any(),
		gomock.Any(),
	).MaxTimes(1).Return(nil)
	repo.EXPECT().Create(
		gomock.Any(),
		gomock.Any(),
	).AnyTimes().Return(errlib.ErrResourceAlreadyExists)

	service := NewResourceService(repo, nil)

	type args struct {
		ctx        context.Context
		resourceId string
		resource   interface{}
	}
	tests := []struct {
		name    string
		r       *resourceService
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "Create not existing resource",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1",
				resource:   serviceRoot,
			},
			want:    serviceRoot,
			wantErr: nil,
		},
		{
			name: "Create existing resource",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1",
				resource:   serviceRoot,
			},
			want:    nil,
			wantErr: errlib.ErrResourceAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Create(tt.args.ctx, tt.args.resourceId, tt.args.resource)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					assert.Equal(t, tt.want, got)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_resourceService_createCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)

	repo.EXPECT().Create(
		gomock.Any(),
		gomock.Any(),
	).MaxTimes(1).Return(nil)
	repo.EXPECT().Create(
		gomock.Any(),
		gomock.Any(),
	).AnyTimes().Return(errlib.ErrResourceAlreadyExists)

	service := NewResourceService(repo, nil)

	type args struct {
		ctx           context.Context
		collectionDto dto.CollectionDto
	}
	tests := []struct {
		name    string
		r       *resourceService
		args    args
		want    *collection
		wantErr error
	}{
		{
			name: "Create not existing collection",
			r:    service,
			args: args{
				ctx: context.Background(),
				collectionDto: dto.CollectionDto{
					OdataId:   "redfish/v1/Storage",
					OdataType: "#StorageCollection_StorageCollection",
					Name:      "Storage Collection",
				},
			},
			want: &collection{
				OdataId:           "redfish/v1/Storage",
				Members:           []domain.OdataV4IdRef{},
				Name:              "Storage Collection",
				OdataType:         "#StorageCollection_StorageCollection",
				MembersOdataCount: 0,
			},
			wantErr: nil,
		},
		{
			name: "Create existing collection",
			r:    service,
			args: args{
				ctx: context.Background(),
				collectionDto: dto.CollectionDto{
					OdataId:   "redfish/v1/Storage",
					OdataType: "#StorageCollection_StorageCollection",
					Name:      "Storage Collection",
				},
			},
			want:    nil,
			wantErr: errlib.ErrResourceAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.createCollection(tt.args.ctx, tt.args.collectionDto)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					assert.Equal(t, tt.want, got)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_resourceService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)
	serviceRootBytes, _ := util.Marshal(serviceRoot)
	repo.EXPECT().Get(gomock.Any(), "/redfish/v1").AnyTimes().Return(
		&dto.ResourceDto{
			Id:   "/redfish/v1",
			Data: serviceRootBytes,
		},
		nil,
	)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().Return(
		nil,
		errlib.ErrNotFound,
	)

	repo.EXPECT().Update(
		gomock.Any(),
		gomock.Any(),
	).AnyTimes().Return(nil)

	service := NewResourceService(repo, nil)

	updatedServiceRoot := serviceRoot
	updatedServiceRoot.Id = "New Id"

	type args struct {
		ctx        context.Context
		resourceId string
		patchData  []byte
	}
	tests := []struct {
		name    string
		r       *resourceService
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "Update existing collection",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1",
				patchData:  []byte(`{"Id": "New Id"}`),
			},
			want:    updatedServiceRoot,
			wantErr: nil,
		},
		{
			name: "Update not existing collection",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1/not/existing",
				patchData:  []byte(`{"Id": "New Id"}`),
			},
			want:    nil,
			wantErr: errlib.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Update(tt.args.ctx, tt.args.resourceId, tt.args.patchData)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					wantBytes, _ := util.Marshal(tt.want)
					gotBytes, _ := util.Marshal(got)
					assert.Equal(t, wantBytes, gotBytes)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_resourceService_Replace(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).
		AnyTimes().Return(nil)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).
		AnyTimes().Return(&dto.ResourceDto{Data: []byte(serviceRootBytes)}, nil)
	repo.EXPECT().DeleteStartsWith(gomock.Any(), gomock.Any()).
		AnyTimes().Return(nil)

	service := NewResourceService(repo, nil)

	updatedServiceRoot := serviceRoot
	updatedServiceRoot.Id = "New Id"

	type args struct {
		ctx        context.Context
		resourceId string
		resource   interface{}
	}
	tests := []struct {
		name    string
		r       ResourceService
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "Replace existing resource",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1",
				resource:   updatedServiceRoot,
			},
			want:    updatedServiceRoot,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Replace(tt.args.ctx, tt.args.resourceId, tt.args.resource)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					assert.Equal(t, tt.want, got)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_resourceService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)

	serviceRootBytes, _ := util.Marshal(serviceRoot)

	repo.EXPECT().Get(gomock.Any(), "/redfish/v1").
		AnyTimes().Return(&dto.ResourceDto{Id: "/redfish/v1", Data: []byte(serviceRootBytes)}, nil)
	repo.EXPECT().DeleteStartsWith(gomock.Any(), "/redfish/v1").
		AnyTimes().Return(nil)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).
		AnyTimes().Return(nil, errlib.ErrNotFound)
	repo.EXPECT().DeleteStartsWith(gomock.Any(), gomock.Any()).
		AnyTimes().Return(errlib.ErrNotFound)

	service := NewResourceService(repo, nil)

	type args struct {
		ctx        context.Context
		resourceId string
	}
	tests := []struct {
		name    string
		r       *resourceService
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "Delete existing resource",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1",
			},
			want:    serviceRoot,
			wantErr: nil,
		},
		{
			name: "Delete not existing resource",
			r:    service,
			args: args{
				ctx:        context.Background(),
				resourceId: "/redfish/v1/another",
			},
			want:    nil,
			wantErr: errlib.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Delete(tt.args.ctx, tt.args.resourceId)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					wantBytes, _ := util.Marshal(tt.want)
					gotBytes, _ := util.Marshal(got)
					assert.Equal(t, wantBytes, gotBytes)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_resourceService_populateResource_with_empty_Storage(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)
	service := NewResourceService(repo, NewSimpleIdGenerator())

	emptyStorage := &domain.Storage{
		Name: "Some Name",
	}

	storageCollection := "redfish/v1/Storage"
	odataType := "#Storage.v1_15_1.Storage"

	got, err := service.populateResource(0, dto.ResourceRequestDto{
		Name:            emptyStorage.Name,
		Id:              emptyStorage.Id,
		OdataType:       odataType,
		Resource:        emptyStorage,
		IdSetter:        func(id string) { emptyStorage.Id = id },
		OdataIdSetter:   func(odataId string) { emptyStorage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { emptyStorage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   storageCollection,
			Name:      "Storage Collection",
			OdataType: "#StorageCollection.StorageCollection",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, *emptyStorage.OdataId, got)
	assert.NotEmpty(t, emptyStorage.Id)
	assert.Equal(t, "redfish/v1/Storage/"+emptyStorage.Id, *emptyStorage.OdataId)
	assert.Equal(t, odataType, *emptyStorage.OdataType)
}

func Test_resourceService_populateResource_with_empty_oDataType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)
	service := NewResourceService(repo, NewSimpleIdGenerator())

	emptyStorage := &domain.Storage{
		Name: "Some Name",
	}

	storageCollection := "redfish/v1/Storage"

	got, err := service.populateResource(0, dto.ResourceRequestDto{
		Name:            emptyStorage.Name,
		Id:              emptyStorage.Id,
		OdataType:       "",
		Resource:        emptyStorage,
		IdSetter:        func(id string) { emptyStorage.Id = id },
		OdataIdSetter:   func(odataId string) { emptyStorage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { emptyStorage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   storageCollection,
			Name:      "Storage Collection",
			OdataType: "#StorageCollection.StorageCollection",
		},
	})

	assert.Error(t, err)
	assert.Empty(t, got)
}

func Test_resourceService_populateResource_with_empty_collection_ODataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)
	service := NewResourceService(repo, NewSimpleIdGenerator())

	emptyStorage := &domain.Storage{
		Name: "Some Name",
	}

	storageCollection := ""

	got, err := service.populateResource(0, dto.ResourceRequestDto{
		Name:            emptyStorage.Name,
		Id:              emptyStorage.Id,
		OdataType:       "",
		Resource:        emptyStorage,
		IdSetter:        func(id string) { emptyStorage.Id = id },
		OdataIdSetter:   func(odataId string) { emptyStorage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { emptyStorage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   storageCollection,
			Name:      "Storage Collection",
			OdataType: "#StorageCollection.StorageCollection",
		},
	})

	assert.Error(t, err)
	assert.Empty(t, got)
}

func Test_resourceService_populateResource_with_provided_Id(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)
	service := NewResourceService(repo, NewSimpleIdGenerator())
	id := "SomeId"
	emptyStorage := &domain.Storage{
		Name: "Some Name",
		Id:   id,
	}

	storageCollection := "redfish/v1/Storage"
	odataType := "#Storage.v1_15_1.Storage"

	got, err := service.populateResource(0, dto.ResourceRequestDto{
		Name:            emptyStorage.Name,
		Id:              emptyStorage.Id,
		OdataType:       odataType,
		Resource:        emptyStorage,
		IdSetter:        func(id string) { emptyStorage.Id = id },
		OdataIdSetter:   func(odataId string) { emptyStorage.OdataId = &odataId },
		OdataTypeSetter: func(odataType string) { emptyStorage.OdataType = &odataType },
		Collection: dto.CollectionDto{
			OdataId:   storageCollection,
			Name:      "Storage Collection",
			OdataType: "#StorageCollection.StorageCollection",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, *emptyStorage.OdataId, got)
	assert.Equal(t, id, emptyStorage.Id)
	assert.Equal(t, "redfish/v1/Storage/"+emptyStorage.Id, *emptyStorage.OdataId)
	assert.Equal(t, odataType, *emptyStorage.OdataType)
}

func Test_resourceService_CreateCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockResourceRepository(ctrl)

	repo.EXPECT().Create(
		gomock.Any(),
		gomock.Any(),
	).MaxTimes(1).Return(nil)
	repo.EXPECT().Create(
		gomock.Any(),
		gomock.Any(),
	).AnyTimes().Return(errlib.ErrResourceAlreadyExists)

	service := NewResourceService(repo, nil)

	type args struct {
		ctx           context.Context
		collectionDto dto.CollectionDto
	}
	tests := []struct {
		name    string
		r       *resourceService
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "Create not existing collection",
			r:    service,
			args: args{
				ctx: context.Background(),
				collectionDto: dto.CollectionDto{
					OdataId:   "redfish/v1/Storage",
					OdataType: "#StorageCollection_StorageCollection",
					Name:      "Storage Collection",
				},
			},
			want: &collection{
				OdataId:           "redfish/v1/Storage",
				Members:           []domain.OdataV4IdRef{},
				Name:              "Storage Collection",
				OdataType:         "#StorageCollection_StorageCollection",
				MembersOdataCount: 0,
			},
			wantErr: nil,
		},
		{
			name: "Create existing collection",
			r:    service,
			args: args{
				ctx: context.Background(),
				collectionDto: dto.CollectionDto{
					OdataId:   "redfish/v1/Storage",
					OdataType: "#StorageCollection_StorageCollection",
					Name:      "Storage Collection",
				},
			},
			want:    nil,
			wantErr: errlib.ErrResourceAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.CreateCollection(tt.args.ctx, tt.args.collectionDto)
			if tt.wantErr == nil {
				if assert.Nil(t, err) && assert.NotNil(t, got) {
					assert.Equal(t, tt.want, got)
				}
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
