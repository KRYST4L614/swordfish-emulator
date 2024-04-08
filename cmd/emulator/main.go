package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/app"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/config"
)

const (
	configFile string = "configs/emulator/config.yaml"
)

func main() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	conf, err := config.ReadConfigFromYAML[app.Config](configFile)
	if err != nil {
		panic(fmt.Errorf("Read of config from '%s' failed: %w", configFile, err))
	}
	err = config.ValidateConfig(conf)
	if err != nil {
		panic(fmt.Errorf("'%s' parsing failed: %w", configFile, err))
	}

	logrus.Info("Starting...")

	notify := make(chan error, 1)
	defer close(notify)

	app, err := app.NewApp(conf, notify)
	if err != nil {
		logrus.Panic(err)
	}

	err = app.Start()
	defer func() {
		appErr := app.Stop()
		if appErr != nil {
			logrus.Error(appErr)
		}
	}()

	if err != nil {
		logrus.Panic(err)
	}

	interupt := make(chan os.Signal, 1)
	defer close(interupt)

	signal.Notify(interupt, os.Interrupt, syscall.SIGTERM)

	select {
	case serr := <-notify:
		logrus.Errorf("Notified with app error: %s", serr.Error())
	case signl := <-interupt:
		logrus.Info("Cought signal while App running: " + signl.String())
	}

	logrus.Info("Shutting down...")
}
