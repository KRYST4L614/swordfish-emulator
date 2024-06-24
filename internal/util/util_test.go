package util

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/domain"
	"gitlab.com/IgorNikiforov/swordfish-emulator-go/internal/errlib"
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

func TestGetParent(t *testing.T) {
	uri := "redfish/v1/Storage/1"
	assert.Equal(t, "redfish/v1/Storage", GetParent(uri))
}

func TestWriteJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteJSONError(rr, errlib.ErrInternal)
	assert.Equal(t, 500, rr.Code)
	decoededError, err := UnmarshalFromReader[errlib.JSONError](rr.Body)
	if assert.NoError(t, err) {
		assert.Equal(t, 500, decoededError.Error.Code)
		assert.Contains(t, decoededError.Error.Message, errlib.ErrInternal.Error())
	}
}

func TestIdGenerator(t *testing.T) {
	gen := IdGenerator()
	cache := make(map[string]struct{}, 10)
	for count := 0; count < 10; count++ {
		id, err := gen()
		if assert.NoError(t, err) {
			_, ok := cache[id]
			if !assert.False(t, ok) {
				t.FailNow()
			}

			cache[id] = struct{}{}
		}
	}
}
