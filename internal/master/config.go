package master

type DatasetConfig struct {
	Path      string `yaml:"path"`
	Overwrite bool   `yaml:"overwrite"`
}
