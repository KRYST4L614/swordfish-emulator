package provider

type DbConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type EmbeddedPsqlConfig struct {
	Name     string `yaml:"name"`
	Port     uint32 `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Version  string `yaml:"version"`
}
