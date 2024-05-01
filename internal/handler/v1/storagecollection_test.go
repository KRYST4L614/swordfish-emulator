package v1_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/handler"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
	mock_service "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service/mock"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
	"go.uber.org/mock/gomock"
)

var storageCollectionJson string = `{
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

var storageJson string = `{
	"@odata.type": "#Storage.v1_16_0.Storage",
	"Id": "TestId",
	"Name": "NVMe IP Attached Drive Configuration",
	"Description": "An NVM Express Subsystem is an NVMe device that contains one or more NVM Express controllers and may contain one or more namespaces.",
	"@odata.id": "/redfish/v1/Storage/TestId",
	"@Redfish.Copyright": "Copyright 2015-2023 SNIA. All rights reserved."
}`

func TestStorageCollectionHandler_createStorage_success(t *testing.T) {
	storage, err := util.Unmarshal[domain.Storage]([]byte(storageJson))

	assert.NoError(t, err)

	ctrl := gomock.NewController(t)

	collectionService := mock_service.NewMockStorageCollectionService(ctrl)
	collectionService.EXPECT().AddStorage(gomock.Any(), "/redfish/v1/Storage", gomock.Any()).AnyTimes().Return(nil)

	handler := handler.NewHandler(&service.Service{
		StorageCollectionService: collectionService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("POST", "/redfish/v1/Storage", strings.NewReader(storageJson))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusCreated, rr.Code) {
		receivedStorage, _ := util.Unmarshal[domain.Storage](rr.Body.Bytes())
		assert.Equal(t, storage, receivedStorage)
	}
}

func TestStorageCollectionHandler_createStorage_fail(t *testing.T) {
	ctrl := gomock.NewController(t)

	collectionService := mock_service.NewMockStorageCollectionService(ctrl)
	collectionService.EXPECT().AddStorage(gomock.Any(), "/redfish/v1/Storage", gomock.Any()).AnyTimes().Return(errlib.ErrResourceAlreadyExists)

	handler := handler.NewHandler(&service.Service{
		StorageCollectionService: collectionService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("POST", "/redfish/v1/Storage", strings.NewReader(storageJson))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestStorageCollectionHandler_getStorageCollection_success(t *testing.T) {
	storageCollection, err := util.Unmarshal[domain.StorageCollection]([]byte(storageCollectionJson))

	assert.NoError(t, err)

	ctrl := gomock.NewController(t)

	collectionService := mock_service.NewMockStorageCollectionService(ctrl)
	collectionService.EXPECT().Get(gomock.Any(), "/redfish/v1/Storage").AnyTimes().Return(storageCollection, nil)

	handler := handler.NewHandler(&service.Service{
		StorageCollectionService: collectionService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("GET", "/redfish/v1/Storage", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code) {
		receivedStorageCollection, _ := util.Unmarshal[domain.StorageCollection](rr.Body.Bytes())
		assert.Equal(t, storageCollection, receivedStorageCollection)
	}
}

func TestStorageCollectionHandler_getStorageCollection_fail(t *testing.T) {
	ctrl := gomock.NewController(t)

	collectionService := mock_service.NewMockStorageCollectionService(ctrl)
	collectionService.EXPECT().Get(gomock.Any(), "/redfish/v1/Storage").AnyTimes().Return(nil, errlib.ErrNotFound)

	handler := handler.NewHandler(&service.Service{
		StorageCollectionService: collectionService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("GET", "/redfish/v1/Storage", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestStorageCollectionHandler_deleteStorage_success(t *testing.T) {
	storage, err := util.Unmarshal[domain.Storage]([]byte(storageJson))
	assert.NoError(t, err)
	ctrl := gomock.NewController(t)

	collectionService := mock_service.NewMockStorageCollectionService(ctrl)
	collectionService.EXPECT().DeleteStorage(gomock.Any(), "/redfish/v1/Storage/1").AnyTimes().Return(storage, nil)

	handler := handler.NewHandler(&service.Service{
		StorageCollectionService: collectionService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("DELETE", "/redfish/v1/Storage/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code) {
		receivedStorage, _ := util.Unmarshal[domain.Storage](rr.Body.Bytes())
		assert.Equal(t, storage, receivedStorage)
	}
}

func TestStorageCollectionHandler_deleteStorage_fail(t *testing.T) {
	ctrl := gomock.NewController(t)

	collectionService := mock_service.NewMockStorageCollectionService(ctrl)
	collectionService.EXPECT().DeleteStorage(gomock.Any(), "/redfish/v1/Storage/1").AnyTimes().Return(nil, errlib.ErrNotFound)

	handler := handler.NewHandler(&service.Service{
		StorageCollectionService: collectionService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("DELETE", "/redfish/v1/Storage/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
