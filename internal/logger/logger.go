package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Stdout bool   `yaml:"stdout"`
	Format string `yaml:"format" validate:"required,oneof=json text"`
	File   bool   `yaml:"file"`
}

type LoggerRotationConfig struct {
	Filename   string `yaml:"file-name" validate:"required"`
	MaxSize    int    `yaml:"max-size" validate:"required"`
	MaxBackups int    `yaml:"max-backups" validate:"required"`
	MaxAge     int    `yaml:"max-age" validate:"required"`
	LocalTime  bool   `yaml:"local-time" validate:"required"`
	Compress   bool   `yaml:"compress"`
}

func GetLogger(loggerRotationConfig *LoggerRotationConfig, loggerConfig *LoggerConfig) *slog.Logger {
	loggerRotation := &lumberjack.Logger{
		Filename:   loggerRotationConfig.Filename,
		MaxSize:    loggerRotationConfig.MaxSize,
		MaxBackups: loggerRotationConfig.MaxBackups,
		MaxAge:     loggerRotationConfig.MaxAge,
		Compress:   loggerRotationConfig.Compress,
		LocalTime:  loggerRotationConfig.LocalTime,
	}

	var logOutput io.Writer

	switch {
	case loggerConfig.File && loggerConfig.Stdout:
		logOutput = io.MultiWriter(loggerRotation, os.Stdout)

	case loggerConfig.File:
		logOutput = loggerRotation
	case loggerConfig.File:
		logOutput = os.Stdout
	}

	var logger *slog.Logger

	if loggerConfig.Format == "json" {
		logger = slog.New(slog.NewJSONHandler(logOutput, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(logOutput, nil))
	}

	return logger
}
