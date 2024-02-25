package provider

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
)

// DbProvider
//
// Structure to make sql request to the DB through sqlx
type DbProvider struct {
	*sqlx.DB
}

// NewPsqlProvider creates new DbProvider to access PostgreSQL DB
func NewPsqlProvider(config *DbConfig) (*DbProvider, error) {
	connectionFmt := "postgresql://@%s/%s?user=%s&password=%s&sslmode=disable"
	db, err := sqlx.Open("pgx", fmt.Sprintf(connectionFmt, config.Host, config.Name, config.User, config.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to add database to pool. Error: %w", errlib.ErrHttpInternal)
	}
	if db.Ping() != nil {
		return nil, fmt.Errorf("failed to ping database. Error: %w", errlib.ErrHttpInternal)
	}

	return &DbProvider{
		DB: db,
	}, nil
}
