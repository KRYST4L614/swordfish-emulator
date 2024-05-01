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

func TestStorageHandler_getStorage_success(t *testing.T) {
	storage, err := util.Unmarshal[domain.Storage]([]byte(storageJson))

	assert.NoError(t, err)

	ctrl := gomock.NewController(t)

	storageService := mock_service.NewMockStorageService(ctrl)
	storageService.EXPECT().Get(gomock.Any(), "/redfish/v1/Storage/1").AnyTimes().Return(storage, nil)

	handler := handler.NewHandler(&service.Service{
		StorageService: storageService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("GET", "/redfish/v1/Storage/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code) {
		receivedStorage, _ := util.Unmarshal[domain.Storage](rr.Body.Bytes())
		assert.Equal(t, storage, receivedStorage)
	}
}

func TestStorageHandler_getStorage_fail(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageService := mock_service.NewMockStorageService(ctrl)
	storageService.EXPECT().Get(gomock.Any(), "/redfish/v1/Storage/1").AnyTimes().Return(nil, errlib.ErrNotFound)

	handler := handler.NewHandler(&service.Service{
		StorageService: storageService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("GET", "/redfish/v1/Storage/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestStorageHandler_replaceStorage_success(t *testing.T) {
	storage, err := util.Unmarshal[domain.Storage]([]byte(storageJson))

	assert.NoError(t, err)

	ctrl := gomock.NewController(t)

	storageService := mock_service.NewMockStorageService(ctrl)
	storageService.EXPECT().Replace(gomock.Any(), "/redfish/v1/Storage/1", gomock.Any()).AnyTimes().Return(storage, nil)

	handler := handler.NewHandler(&service.Service{
		StorageService: storageService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("PUT", "/redfish/v1/Storage/1", strings.NewReader(storageJson))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code) {
		receivedStorage, _ := util.Unmarshal[domain.Storage](rr.Body.Bytes())
		assert.Equal(t, storage, receivedStorage)
	}
}

func TestStorageHandler_replaceStorage_fail(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageService := mock_service.NewMockStorageService(ctrl)
	storageService.EXPECT().Replace(gomock.Any(), "/redfish/v1/Storage/1", gomock.Any()).AnyTimes().Return(nil, errlib.ErrNotFound)

	handler := handler.NewHandler(&service.Service{
		StorageService: storageService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("PUT", "/redfish/v1/Storage/1", strings.NewReader(storageJson))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestStorageHandler_updateStorage_success(t *testing.T) {
	storage, err := util.Unmarshal[domain.Storage]([]byte(storageJson))

	assert.NoError(t, err)

	ctrl := gomock.NewController(t)

	storageService := mock_service.NewMockStorageService(ctrl)
	storageService.EXPECT().Update(gomock.Any(), "/redfish/v1/Storage/1", gomock.Any()).AnyTimes().Return(storage, nil)

	handler := handler.NewHandler(&service.Service{
		StorageService: storageService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("PATCH", "/redfish/v1/Storage/1", strings.NewReader(storageJson))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code) {
		receivedStorage, _ := util.Unmarshal[domain.Storage](rr.Body.Bytes())
		assert.Equal(t, storage, receivedStorage)
	}
}

func TestStorageHandler_updateStorage_fail(t *testing.T) {
	ctrl := gomock.NewController(t)

	storageService := mock_service.NewMockStorageService(ctrl)
	storageService.EXPECT().Update(gomock.Any(), "/redfish/v1/Storage/1", gomock.Any()).AnyTimes().Return(nil, errlib.ErrNotFound)

	handler := handler.NewHandler(&service.Service{
		StorageService: storageService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("PATCH", "/redfish/v1/Storage/1", strings.NewReader(storageJson))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
