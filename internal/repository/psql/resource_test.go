package psql

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
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
		WillReturnError(errors.New("another")).
		WillReturnError(errors.New("resource_pkey"))
	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	root := &dto.ResourceDto{
		Id:   "/redfish/v1",
		Data: []byte("somedata"),
	}

	err = repo.Create(context.Background(), root)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errlib.ErrResourceAlreadyExists)

	err = repo.Create(context.Background(), root)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errlib.ErrInternal)
}

func TestServiceRootRepository_DeleteAll(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectExec("TRUNCATE resource CASCADE").WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectExec("TRUNCATE resource CASCADE").WillReturnError(&pq.Error{Code: pq.ErrorCode("internal")})

	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	err = repo.DeleteAll(context.Background())
	assert.NoError(t, err)

	err = repo.DeleteAll(context.Background())
	assert.Error(t, err)
}

func TestServiceRootRepository_DeleteById(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectExec("DELETE FROM resource WHERE id=(.*)").
		WithArgs("some_id").
		WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM resource WHERE id=(.*)").
		WillReturnError(&pq.Error{Code: pq.ErrorCode("internal")})

	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	err = repo.DeleteById(context.Background(), "some_id")
	assert.NoError(t, err)

	err = repo.DeleteById(context.Background(), "other_id")
	assert.Error(t, err)
}

func TestServiceRootRepository_Update(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectExec("UPDATE resource SET data = (.*) WHERE id = (.*)").
		WithArgs([]byte(`some: "json"`), "some_id").
		WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE resource SET data = (.*) WHERE id = (.*)").
		WillReturnError(&pq.Error{Code: pq.ErrorCode("internal")})

	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	err = repo.Update(context.Background(), &dto.ResourceDto{
		Id:   "some_id",
		Data: []byte(`some: "json"`),
	})
	assert.NoError(t, err)

	err = repo.Update(context.Background(), &dto.ResourceDto{
		Id:   "other_id",
		Data: []byte(`some: "json"`),
	})
	assert.Error(t, err)
}

func TestServiceRootRepository_DeleteStartsWith(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectExec("DELETE FROM resource WHERE id LIKE (.*)").
		WithArgs("some_id%").
		WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM resource WHERE id LIKE (.*)").
		WillReturnError(&pq.Error{Code: pq.ErrorCode("internal")})

	repo := NewPsqlResourceRepository(&provider.DbProvider{DB: db})

	err = repo.DeleteStartsWith(context.Background(), "some_id")
	assert.NoError(t, err)

	err = repo.DeleteStartsWith(context.Background(), "other_id")
	assert.Error(t, err)
}
