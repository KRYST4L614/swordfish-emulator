package impl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getParent(t *testing.T) {
	uri := "redfish/v1/Storage/1"
	assert.Equal(t, "redfish/v1/Storage", getParent(uri))
}
