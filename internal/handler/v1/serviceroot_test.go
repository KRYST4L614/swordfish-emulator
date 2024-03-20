package v1_test

import (
	"net/http"
	"net/http/httptest"
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

var serviceRoot = &domain.ServiceRoot{
	Base: domain.Base{
		Name: "SomeName",
		Id:   "Some Id",
		InlineODataId: domain.InlineODataId{
			ODataId: "/redfish/v1",
		},
	},
}

func TestServiceRootHandler_getServiceRoot_success(t *testing.T) {
	ctrl := gomock.NewController(t)

	rootService := mock_service.NewMockServiceRootService(ctrl)
	rootService.EXPECT().Get(gomock.Any(), "/redfish/v1").AnyTimes().Return(serviceRoot, nil)
	rootService.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errlib.ErrInternal)

	handler := handler.NewHandler(&service.Service{
		ServiceRootService: rootService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("GET", "/redfish/v1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code) {
		receivedRoot, _ := util.Unmarshal[domain.ServiceRoot](rr.Body.Bytes())
		assert.Equal(t, serviceRoot, receivedRoot)
	}
}

func TestServiceRootHandler_getServiceRoot_failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	rootService := mock_service.NewMockServiceRootService(ctrl)
	rootService.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errlib.ErrInternal)

	handler := handler.NewHandler(&service.Service{
		ServiceRootService: rootService,
	})

	router := mux.NewRouter()
	handler.SetRouter(router)

	req := httptest.NewRequest("GET", "/redfish/v1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusInternalServerError, rr.Code) {
		_, err := util.Unmarshal[errlib.JSONError](rr.Body.Bytes())
		assert.NoError(t, err)
	}
}
