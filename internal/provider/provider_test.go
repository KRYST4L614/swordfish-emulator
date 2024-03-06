package provider

import (
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

type BeerCatalogue struct {
	ID       int64
	Name     string
	Consumed bool
	Rating   float64
}

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

func TestEmbeddedPsqlOnPersistentData(t *testing.T) {
	tempDir := t.TempDir()
	db := NewEmbeddedPsql(&EmbeddedPsqlConfig{
		Name:     "name",
		Port:     5432,
		UserName: "user",
		Password: "pass",
		Version:  "v13.8",
		DataPath: tempDir,
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

	// Create table and insert default data
	err = goose.Up(provider.DB.DB, "./testdata")
	assert.NoError(t, err)

	// End DB process
	err = db.Stop()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	// Start DB again and check if data exist
	err = db.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	beers := make([]BeerCatalogue, 0)
	err = provider.Select(&beers, "SELECT * FROM beer_catalogue")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	assert.Len(t, beers, 1)
}
