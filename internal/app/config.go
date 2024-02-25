package app

import "gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/provider"

type Config struct {
	DbConfig provider.DbConfig `yaml:"db"`
}
