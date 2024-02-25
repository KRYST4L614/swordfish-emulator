package provider

type DbConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}
