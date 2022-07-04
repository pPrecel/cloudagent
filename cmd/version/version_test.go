package version

import (
	"testing"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		o := &command.Options{}
		c := NewCmd(o)

		assert.Equal(t, 0, len(c.Commands()))
		assert.NoError(t, c.RunE(c, nil))
	})
}
