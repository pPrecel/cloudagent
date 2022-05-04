package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectSchema(t *testing.T) {
	t.Run("reflect schema", func(t *testing.T) {
		b, err := JSONSchema()
		assert.NoError(t, err)
		assert.NotEmpty(t, b)
	})
}
