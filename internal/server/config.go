package server

import "time"

type ServerConfig struct {
	Host         *string        `yaml:"host" validate:"hostname_port"`
	ReadTimeout  *time.Duration `yaml:"read-timeout"`
	WriteTimeout *time.Duration `yaml:"write-timeout"`
}
