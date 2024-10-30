package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log"
	"log/slog"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/app"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/config"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/logger"
)

const (
	configFile string = "configs/emulator/config.yaml"
)

func main() {

	conf, err := config.ReadConfigFromYAML[app.Config](configFile)
	if err != nil {
		panic(fmt.Errorf("read of config from '%s' failed: %w", configFile, err))
	}
	err = config.ValidateConfig(conf)
	if err != nil {
		panic(fmt.Errorf("'%s' parsing failed: %w", configFile, err))
	}

	slog.SetDefault(logger.GetLogger(conf.LoggerRotationConfig, conf.LoggerConfig))

	slog.Info("Starting...")

	notify := make(chan error, 1)
	defer close(notify)

	app, err := app.NewApp(conf, notify)
	if err != nil {
		log.Panic(err)
	}

	err = app.Start()
	defer func() {
		appErr := app.Stop()
		if appErr != nil {
			slog.Error(appErr.Error())
		}
	}()

	if err != nil {
		log.Panic(err)
	}

	interupt := make(chan os.Signal, 1)
	defer close(interupt)

	signal.Notify(interupt, os.Interrupt, syscall.SIGTERM)

	select {
	case serr := <-notify:
		slog.Error(fmt.Sprintf("Notified with app error: %s", serr.Error()))
	case signl := <-interupt:
		slog.Info("Cought signal while App running: " + signl.String())
	}

	slog.Info("Shutting down...")
}
