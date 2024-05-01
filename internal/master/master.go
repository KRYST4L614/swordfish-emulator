package master

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
)

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

func (m *InitialConfigurationMaster) LoadResources(datasetPath string) error {
	return filepath.WalkDir(datasetPath, m.loadAllJson)
}

func (m *InitialConfigurationMaster) loadAllJson(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() || filepath.Ext(path) != ".json" {
		return nil
	}

	logrus.Tracef("Configuration master loaded resource from %s", path)
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	id, err := idFromPath(m.config.Path, path)
	if err != nil {
		return err
	}

	logrus.Tracef("New resource with id %s", id)
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
	if filepath.Base(path) == "index.json" {
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
