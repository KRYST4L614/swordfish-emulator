package impl_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
	mock_repository "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/mock"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service/impl"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
	"go.uber.org/mock/gomock"
)

var storageCollectionJson = `{
	"@odata.type": "#StorageCollection.StorageCollection",
	"Name": "Storage Collection",
	"Members@odata.count": 2,
	"Members": [
			{
					"@odata.id": "/redfish/v1/Storage/IPAttachedDrive1"
			},
			{
					"@odata.id": "/redfish/v1/Storage/IPAttachedDrive2"
			}
	],
	"@odata.id": "/redfish/v1/Storage",
	"@Redfish.Copyright": "Copyright 2015-2023 SNIA. All rights reserved."
}`

var storageJson = `{
	"@odata.type": "#Storage.v1_16_0.Storage",
	"Id": "1",
	"Name": "NVMe IP Attached Drive Configuration",
	"Description": "An NVM Express Subsystem is an NVMe device that contains one or more NVM Express controllers and may contain one or more namespaces.",
	"Status": {
			"State": "Enabled",
			"Health": "OK",
			"HealthRollup": "OK"
	},
	"Identifiers": [
			{
					"DurableNameFormat": "NQN",
					"DurableName": "nqn.2014-08.org.nvmexpress:uuid:6c5fe566-10e6-4fb6-aad4-8b4159f50245"
			}
	],
	"Controllers": {
			"@odata.id": "/redfish/v1/Storage/IPAttachedDrive1/Controllers"
	},
	"Drives": [
			{
					"@odata.id": "/redfish/v1/Chassis/EBOFEnclosure/Drives/IPAttachedDrive1"
			}
	],
	"Volumes": {
			"@odata.id": "/redfish/v1/Storage/IPAttachedDrive1/Volumes"
	},
	"Links": {
			"Enclosures": [
					{
							"@odata.id": "/redfish/v1/Chassis/EBOFEnclosure"
					}
			]
	},
	"@odata.id": "/redfish/v1/Storage/1",
	"@Redfish.Copyright": "Copyright 2015-2023 SNIA. All rights reserved."
}`

func TestStorageCollectionService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockResourceRepository(ctrl)

	dto := &dto.ResourceDto{
		Id:   "/redfish/v1/Storage",
		Data: []byte(storageCollectionJson),
	}
	mockRepo.EXPECT().Get(gomock.Any(), "/redfish/v1/Storage").AnyTimes().Return(dto, nil)
	mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errlib.ErrNotFound)

	service := impl.NewStorageCollectionService(mockRepo)

	collection, err := service.Get(context.Background(), "/redfish/v1/Storage")
	if assert.NoError(t, err) {
		assert.Equal(t, "Storage Collection", collection.Name)
		assert.Equal(t, "/redfish/v1/Storage", *collection.OdataId)
	}

	collection, err = service.Get(context.Background(), "not_existing")
	assert.Error(t, err)
	assert.Nil(t, collection)
}

func TestStorageCollectionService_AddStorage(t *testing.T) {
	var finalCollectionJson = `{
		"@odata.type": "#StorageCollection.StorageCollection",
		"Name": "Storage Collection",
		"Members@odata.count": 3,
		"Members": [
				{
						"@odata.id": "/redfish/v1/Storage/IPAttachedDrive1"
				},
				{
						"@odata.id": "/redfish/v1/Storage/IPAttachedDrive2"
				},
				{
						"@odata.id": "/redfish/v1/Storage/1"
				}
		],
		"@odata.id": "/redfish/v1/Storage",
		"@Redfish.Copyright": "Copyright 2015-2023 SNIA. All rights reserved."
	}`

	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockResourceRepository(ctrl)

	storage, _ := util.Unmarshal[domain.Storage]([]byte(storageJson))

	storageBytes, _ := util.Marshal(storage)
	storageDto := &dto.ResourceDto{
		Id:   "/redfish/v1/Storage/1",
		Data: storageBytes,
	}

	storageCollection, _ := util.Unmarshal[domain.StorageCollection]([]byte(finalCollectionJson))
	collectionBytes, _ := util.Marshal(storageCollection)
	collectionDto := &dto.ResourceDto{
		Id:   "/redfish/v1/Storage",
		Data: []byte(storageCollectionJson),
	}

	mockRepo.EXPECT().Get(gomock.Any(), "/redfish/v1/Storage").AnyTimes().Return(collectionDto, nil)
	mockRepo.EXPECT().Create(gomock.Any(), storageDto).AnyTimes().Return(nil)
	mockRepo.EXPECT().Update(gomock.Any(), &dto.ResourceDto{
		Id:   "/redfish/v1/Storage",
		Data: collectionBytes,
	}).AnyTimes().Return(nil)

	service := impl.NewStorageCollectionService(mockRepo)

	err := service.AddStorage(context.Background(), "/redfish/v1/Storage", storage)
	assert.NoError(t, err)
}
