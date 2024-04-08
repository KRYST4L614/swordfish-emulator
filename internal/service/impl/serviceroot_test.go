package impl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
	mock_repository "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/mock"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
	"go.uber.org/mock/gomock"
)

var serviceRoot = &domain.ServiceRoot{
	OdataId: util.Addr("/redfish/v1"),
	Name:    "SomeName",
	Id:      "Some Id",
}

func TestServiceRootService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockResourceRepository(ctrl)

	bytes, _ := util.Marshal(serviceRoot)
	mockRepo.EXPECT().Create(gomock.Any(), &dto.ResourceDto{Id: "/redfish/v1", Data: bytes}).AnyTimes().Return(nil)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(errlib.ErrInternal)

	service := NewServiceRootService(mockRepo)

	err := service.Create(context.Background(), serviceRoot)
	assert.NoError(t, err)

	err = service.Create(context.Background(), &domain.ServiceRoot{})
	assert.ErrorIs(t, err, errlib.ErrInternal)

	err = service.Create(context.Background(), &domain.ServiceRoot{OdataId: util.Addr("some")})
	assert.Error(t, err)
}

func TestServiceRootService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockResourceRepository(ctrl)

	bytes, _ := util.Marshal(serviceRoot)
	dto := &dto.ResourceDto{
		Id:   "/redfish/v1",
		Data: bytes,
	}
	mockRepo.EXPECT().Get(gomock.Any(), "/redfish/v1").AnyTimes().Return(dto, nil)
	mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errlib.ErrNotFound)

	service := NewServiceRootService(mockRepo)

	serviceRoot, err := service.Get(context.Background(), "/redfish/v1")
	if assert.NoError(t, err) {
		assert.Equal(t, "SomeName", serviceRoot.Name)
		assert.Equal(t, "Some Id", serviceRoot.Id)
	}

	serviceRoot, err = service.Get(context.Background(), "notexist")
	assert.Error(t, err)
	assert.Nil(t, serviceRoot)
}
