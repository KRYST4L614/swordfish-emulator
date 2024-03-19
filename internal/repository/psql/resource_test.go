package psql

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
)

var serviceRootBytes = `
{
	"@odata.type": "#ServiceRoot.v1_15_0.ServiceRoot",
	"Id": "RootService",
	"Name": "Root Service",
	"RedfishVersion": "1.18.0",
	"UUID": "92384634-2938-2342-8820-489239905423",
	"Chassis": {
		"@odata.id": "/redfish/v1/Chassis"
	},
	"Fabrics": {
		"@odata.id": "/redfish/v1/Fabrics"
	},
	"Managers": {
		"@odata.id": "/redfish/v1/Managers"
	},
	"SessionService": {
		"@odata.id": "/redfish/v1/SessionService"
	},
	"Registries": {
		"@odata.id": "/redfish/v1/Registries"
	},
	"Storage": {
		"@odata.id": "/redfish/v1/Storage"
	},
	"Links": {
		"Sessions": {
			"@odata.id": "/redfish/v1/SessionService/Sessions"
		}
	},
	"@odata.id": "/redfish/v1"
}`

func TestResourceRepository_Get_NoError(t *testing.T) {
	var resourceRows = sqlxmock.NewRows([]string{"id", "data"}).
		AddRow("/redfish/v1", []byte(serviceRootBytes))

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectQuery(`^SELECT (.+) FROM resource WHERE id=(.+)$`).
		WillReturnRows(resourceRows)

	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	root, err := repo.Get(context.Background(), "/redfish/v1")

	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	if assert.NotNil(t, root) {
		assert.Equal(t, "/redfish/v1", root.Id)
	}
}

func TestServiceRootRepository_Get_WithError(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectQuery(`^SELECT (.+) FROM resource WHERE id=(.+)$`).
		WithArgs("/redfish/v1").
		WillReturnError(&pq.Error{Code: pq.ErrorCode("internal")}).
		WillReturnError(&pq.Error{Code: pq.ErrorCode("02000")})

	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	root, err := repo.Get(context.Background(), "/redfish/v1")

	assert.Error(t, err)
	assert.Nil(t, root)
	assert.ErrorIs(t, err, errlib.ErrNotFound)

	root, err = repo.Get(context.Background(), "/redfish/v1")

	assert.Error(t, err)
	assert.Nil(t, root)
	assert.ErrorIs(t, err, errlib.ErrInternal)
}

func TestServiceRootRepository_Create_NoError(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectExec("INSERT INTO resource").
		WithArgs("/redfish/v1", sqlxmock.AnyArg()).
		WillReturnResult(sqlxmock.NewResult(1, 1))
	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	root := &dto.ResourceDto{
		Id:   "/redfish/v1",
		Data: []byte("somedata"),
	}
	err = repo.Create(context.Background(), root)
	assert.NoError(t, err)
}

func TestServiceRootRepository_Create_WithError(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectExec("INSERT INTO resource").
		WithArgs("/redfish/v1", sqlxmock.AnyArg()).
		WillReturnError(&pq.Error{Code: pq.ErrorCode("internal")}).
		WillReturnError(&pq.Error{Code: pq.ErrorCode("23505")})
	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	root := &dto.ResourceDto{
		Id:   "/redfish/v1",
		Data: []byte("somedata"),
	}

	err = repo.Create(context.Background(), root)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errlib.ErrResourceExists)

	err = repo.Create(context.Background(), root)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errlib.ErrInternal)
}
