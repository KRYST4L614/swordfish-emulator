package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/repository"
)

func TestNewService(t *testing.T) {
	service := NewService(&repository.Repository{})
	assert.NotNil(t, service)
}
