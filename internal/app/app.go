package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/pressly/goose/v3"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/handler"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/master"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/server"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/service"
)

type App struct {
	server   *server.Server
	provider *provider.DbProvider
	embedded *provider.EmbeddedPsql
	master   *master.InitialConfigurationMaster
	config   *Config
	repos    *repository.Repository
}

func NewApp(config *Config, notify chan error) (*App, error) {
	var embedded *provider.EmbeddedPsql
	if config.EmbeddedConfig != nil {
		embedded = provider.NewEmbeddedPsql(config.EmbeddedConfig)
	}

	provider, err := provider.NewPsqlProvider(&config.DbConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize db with error: %w", err)
	}

	repos := repository.NewRepository(provider)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	router := mux.NewRouter()
	handlers.SetRouter(router)

	server := server.NewServer(&config.ServerConfig, router, notify)

	return &App{
		server:   server,
		provider: provider,
		embedded: embedded,
		master:   master.NewInitialConfigurationMaster(repos),
		config:   config,
		repos:    repos,
	}, nil
}

func (app *App) Start() error {
	if app.embedded != nil {
		if err := app.embedded.Start(); err != nil {
			return err
		}
	}
	if err := goose.Up(app.provider.DB.DB, "./database/migration"); err != nil {
		return err
	}

	if err := app.master.LoadResources(app.config.DatasetConfig.Path); err != nil {
		if !errors.Is(err, errlib.ErrResourceAlreadyExists) {
			return err
		}

		if errors.Is(err, errlib.ErrResourceAlreadyExists) && app.config.DatasetConfig.Overwrite {
			if err := app.repos.ResourceRepository.DeleteAll(context.Background()); err != nil {
				return err
			}
			if err := app.master.LoadResources(app.config.DatasetConfig.Path); err != nil {
				return err
			}
		}
	}
	app.server.Start()
	return nil
}

func (app *App) Stop() error {
	var errComposition error
	if err := app.server.Stop(); err != nil {
		errComposition = fmt.Errorf("%w. %w", err, errComposition)
	}
	if err := app.provider.Close(); err != nil {
		errComposition = fmt.Errorf("%w. %w", err, errComposition)
	}
	if app.embedded != nil {
		if err := app.embedded.Stop(); err != nil {
			errComposition = fmt.Errorf("%w. %w", err, errComposition)
		}
	}
	return errComposition
}
