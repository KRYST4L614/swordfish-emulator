package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimpleIdGenerator(t *testing.T) {
	gen := NewSimpleIdGenerator()
	assert.NotNil(t, gen)
}

func Test_simpleIdGenerator_Generate(t *testing.T) {
	gen := NewSimpleIdGenerator()

	if assert.NotNil(t, gen) {
		id1, err := gen.Generate(123)
		assert.NoError(t, err)
		id2, err := gen.Generate(124)
		assert.NoError(t, err)

		assert.NotEqual(t, id1, id2)
	}
}
