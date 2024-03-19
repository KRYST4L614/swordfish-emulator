package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
)

func TestMarshal(t *testing.T) {
	root := &domain.ServiceRoot{
		Base: domain.Base{
			InlineODataId: domain.InlineODataId{
				ODataId: "redfish/v1",
			},
		},
	}
	bytes, err := Marshal(root)

	if assert.NoError(t, err) {
		assert.NotEmpty(t, bytes)
	}
}

func TestUnmarshal(t *testing.T) {
	root := &domain.ServiceRoot{
		Base: domain.Base{
			InlineODataId: domain.InlineODataId{
				ODataId: "redfish/v1",
			},
		},
	}

	bytes, err := Marshal(root)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, bytes)
	}

	decodedRoot, err := Unmarshal[domain.ServiceRoot](bytes)
	if assert.NoError(t, err) {
		assert.Equal(t, "redfish/v1", decodedRoot.ODataId)
	}
}
