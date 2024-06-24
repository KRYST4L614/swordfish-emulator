package logger

type LoggerConfig struct {
	LogFile  string `yaml:"log-file" validate:"required"`
	WarnFile string `yaml:"warn-file" validate:"required"`
}

// TODO: replace logrus with slog, make configurable
// log level and write in files
