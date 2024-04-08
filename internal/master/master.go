package master

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository/dto"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/util"
)

type InitialConfigurationMaster struct {
	repo *repository.Repository
}

func NewInitialConfigurationMaster(repo *repository.Repository) *InitialConfigurationMaster {
	return &InitialConfigurationMaster{
		repo: repo,
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

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	id, err := util.Unmarshal[domain.InlineODataId](content)
	if err != nil {
		return err
	}

	dto := &dto.ResourceDto{
		Id:   id.ODataId,
		Data: content,
	}

	err = m.repo.ResourceRepository.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	return nil
}
