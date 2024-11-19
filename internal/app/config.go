package app

import (
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/logger"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/master"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/server"
)

type Config struct {
	ServerConfig         server.ServerConfig          `yaml:"server"`
	DbConfig             provider.DbConfig            `yaml:"db"`
	EmbeddedConfig       *provider.EmbeddedPsqlConfig `yaml:"embedded-psql"`
	DatasetConfig        master.DatasetConfig         `yaml:"dataset"`
	LoggerConfig         *logger.LoggerConfig         `yaml:"logger"`
	LoggerRotationConfig *logger.LoggerRotationConfig `yaml:"logger-rotation"`
}
