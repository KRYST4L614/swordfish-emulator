package util

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
)

func TestMarshal(t *testing.T) {
	root := &domain.ServiceRoot{
		OdataId: Addr("/v1/redfish"),
	}
	bytes, err := Marshal(root)

	if assert.NoError(t, err) {
		assert.NotEmpty(t, bytes)
	}
}

func TestUnmarshal(t *testing.T) {
	root := &domain.ServiceRoot{
		OdataId: Addr("/v1/redfish"),
	}

	bytes, err := Marshal(root)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, bytes)
	}

	decodedRoot, err := Unmarshal[domain.ServiceRoot](bytes)
	if assert.NoError(t, err) {
		assert.Equal(t, "/v1/redfish", *decodedRoot.OdataId)
	}
}

func TestUnmarshalFromReader(t *testing.T) {
	root := &domain.ServiceRoot{
		OdataId: Addr("/v1/redfish"),
	}

	bytes, err := Marshal(root)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, bytes)
	}

	decodedRoot, err := UnmarshalFromReader[domain.ServiceRoot](strings.NewReader(string(bytes)))
	if assert.NoError(t, err) {
		assert.Equal(t, "/v1/redfish", *decodedRoot.OdataId)
	}
}

func TestWriteJSON(t *testing.T) {
	root := &domain.ServiceRoot{
		OdataId: Addr("/v1/redfish"),
	}
	rr := httptest.NewRecorder()
	WriteJSON(rr, root)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	decodedRoot, err := UnmarshalFromReader[domain.ServiceRoot](rr.Body)
	if assert.NoError(t, err) {
		assert.Equal(t, "/v1/redfish", *decodedRoot.OdataId)
	}
}
