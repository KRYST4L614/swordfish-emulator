package provider

import (
	"fmt"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
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
		return nil, fmt.Errorf("failed to add database to pool. Error: %w", errlib.ErrInternal)
	}

	return &DbProvider{
		DB: db,
	}, nil
}

type EmbeddedPsql struct {
	*embeddedpostgres.EmbeddedPostgres
}

func NewEmbeddedPsql(config *EmbeddedPsqlConfig) *EmbeddedPsql {
	db := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Database(config.Name).
		Port(config.Port).
		Username(config.UserName).
		Password(config.Password).
		Logger(nil).
		DataPath(config.DataPath))

	return &EmbeddedPsql{
		EmbeddedPostgres: db,
	}
}
