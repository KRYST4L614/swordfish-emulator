package provider

import (
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddedPsqlConnection(t *testing.T) {
	db := NewEmbeddedPsql(&EmbeddedPsqlConfig{
		Name:     "name",
		Port:     5432,
		UserName: "user",
		Password: "pass",
		Version:  "v13.8",
	})

	err := db.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	t.Cleanup(func() {
		err := db.Stop()

		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
	})

	provider, err := NewPsqlProvider(&DbConfig{
		Host:     "localhost:5432",
		Name:     "name",
		User:     "user",
		Password: "pass",
	})
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	t.Cleanup(func() {
		err := provider.Close()

		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
	})
}

func TestEmbeddedPsqlWithMigrations(t *testing.T) {
	db := NewEmbeddedPsql(&EmbeddedPsqlConfig{
		Name:     "name",
		Port:     5432,
		UserName: "user",
		Password: "pass",
		Version:  "v13.8",
	})

	err := db.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	t.Cleanup(func() {
		err := db.Stop()

		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
	})

	provider, err := NewPsqlProvider(&DbConfig{
		Host:     "localhost:5432",
		Name:     "name",
		User:     "user",
		Password: "pass",
	})
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	t.Cleanup(func() {
		err := provider.Close()

		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
	})

	err = goose.Up(provider.DB.DB, "./testdata")
	assert.NoError(t, err)
}
