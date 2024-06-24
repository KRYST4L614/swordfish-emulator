package master

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
)

// InitialConfigurationMaster loads all default json resources from
// /datasets
type InitialConfigurationMaster struct {
	repo   *repository.Repository
	config *DatasetConfig
}

func NewInitialConfigurationMaster(repo *repository.Repository, config *DatasetConfig) *InitialConfigurationMaster {
	return &InitialConfigurationMaster{
		repo:   repo,
		config: config,
	}
}

// LoadResources main function that loads all files from datasetPath and
// stores them in DB
func (m *InitialConfigurationMaster) LoadResources(datasetPath string) error {
	return filepath.WalkDir(datasetPath, m.loadAllJson)
}

func (m *InitialConfigurationMaster) loadAllJson(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	extension := filepath.Ext(path)
	if d.IsDir() || !(extension == ".json" || extension == ".xml") {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	id, err := idFromPath(m.config.Path, path)
	if err != nil {
		return err
	}

	dto := &dto.ResourceDto{
		Id:   id,
		Data: content,
	}

	err = m.repo.ResourceRepository.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	return nil
}

func idFromPath(parent, path string) (string, error) {
	var id = path
	if filepath.Base(path) == "index.json" || filepath.Base(path) == "index.xml" {
		id = filepath.Dir(path)
	}
	id, err := filepath.Rel(parent, id)
	if err != nil {
		return "", err
	}
	if id == "." {
		id = ""
	}

	id = "/redfish/v1/" + filepath.ToSlash(id)
	if id[len(id)-1] == '/' {
		id = id[:len(id)-1]
	}

	return id, nil
}
